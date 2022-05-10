package handler

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	"github.com/micro/services/pkg/tenant"
	"github.com/micro/services/wordle/lib"
	pb "github.com/micro/services/wordle/proto"
)

type Wordle struct {
	sync.RWMutex

	// current wordle
	w *lib.Wordle
}

func key(date string) string {
	return "game-" + date
}

// Return a new handler
func New() *Wordle {
	// load the current word
	wordle := new(lib.Wordle)
	date := time.Now().Format("2006-01-02")

	// retrieve from the store
	recs, err := store.Read(key(date))
	if err != nil || len(recs) == 0 {
		// word doesn't exist
		wordle = lib.NewWordle()
	} else {
		// load the game
		if err := wordle.Load(recs[0].Value); err != nil {
			wordle = lib.NewWordle()
		}
	}

	w := &Wordle{
		// current wordle
		w: wordle,
	}

	go w.run()
	return w
}

func (w *Wordle) run() {
	t := time.NewTicker(time.Second * 10)
	date := time.Now().Format("2006-01-02")
	var game []byte

	for _ = range t.C {
		data, err := w.w.Save()
		if err != nil {
			logger.Errorf("Failed to save game: %v", err)
		}
		// nothing changed
		if bytes.Equal(data, game) {
			continue
		}

		// overwrite it
		game = data

		// save the game
		logger.Infof("Saving game")
		if err := store.Write(&store.Record{
			Key:   key(date),
			Value: game,
		}); err != nil {
			logger.Errorf("Failed to save game: %v", err)
		}

		// get the date
		today := time.Now().Format("2006-01-02")

		// if we've rolled over, create a new wordle
		if today != date {
			w.Lock()
			w.w = lib.NewWordle()
			w.Unlock()
			date = today
			logger.Info("New day, new game")
		}
	}
}

func (w *Wordle) Guess(ctx context.Context, req *pb.GuessRequest, rsp *pb.GuessResponse) error {
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		tnt = "default"
	}

	// player one has entered the game
	if len(req.Player) == 0 {
		req.Player = "1"
	}

	if len(req.Word) == 0 {
		return errors.BadRequest("wordle.guess", "invalid guess")
	}

	player := tnt + "-" + req.Player

	answer, guesses, err := w.w.Guess(player, req.Word)
	if err == nil && len(answer) > 0 {
		rsp.Correct = true
		rsp.Status = "You won!"
	}

	if len(answer) > 0 {
		rsp.Answer = answer
	}

	rsp.TriesLeft = w.w.Tries - int32(len(guesses))
	rsp.Guesses = guesses

	if err != nil {
		rsp.Status = err.Error()
	}

	return nil
}

func (w *Wordle) Next(ctx context.Context, req *pb.NextRequest, rsp *pb.NextResponse) error {
	now := time.Now()
	year, month, day := time.Now().AddDate(0, 0, 1).Date()
	t := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
	d := time.Until(t)
	rsp.Seconds = int32(d.Seconds())
	rsp.Duration = fmt.Sprintf("%v", d)
	return nil
}
