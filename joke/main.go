package main

import (
	"errors"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/micro/services/joke/handler"
	"github.com/micro/services/joke/model"
	pb "github.com/micro/services/joke/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

var (
	sources = []model.JokeSource{
		{
			Source: "https://www.reddit.com/r/jokes",
			Api:    "https://raw.githubusercontent.com/taivop/joke-dataset/master/reddit_jokes.json",
		},
		{
			Source: "http://wocka.com/",
			Api:    "https://raw.githubusercontent.com/taivop/joke-dataset/master/wocka.json",
		},
		{
			Source: "http://stupidstuff.org/jokes/",
			Api:    "https://raw.githubusercontent.com/taivop/joke-dataset/master/stupidstuff.json",
		},
	}
)

// loadJokes is used to load jokes in store
func loadJokes() error {
	wg := sync.WaitGroup{}
	errs := make([]string, 0)
	for _, v := range sources {
		wg.Add(1)

		go func(s model.JokeSource) {
			defer wg.Done()

			err := model.NewJoke(s).Load()

			if err != nil {
				logger.Errorf("load jokes error: %s", err)
				errs = append(errs, err.Error())
				return
			}

		}(v)

	}

	wg.Wait()

	if len(errs) != 0 {
		return errors.New("load jokes error: " + strings.Join(errs, ";"))
	}

	logger.Infof("loaded jokes: %d\n", len(model.GetAllJokes()))

	return nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	// Create service
	srv := service.New(
		service.Name("joke"),
		service.Version("latest"),
		service.BeforeStart(loadJokes),
	)

	// Register handler
	pb.RegisterJokeHandler(srv.Server(), new(handler.Joke))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
