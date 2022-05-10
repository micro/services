package lib

import (
	_ "embed"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	pb "github.com/micro/services/wordle/proto"
)

//go:embed words.txt
var wordLib string

//go:embed all.txt
var allWordLib string

var wordArr []string
var allWordArr []string

var rounds = 6

func init() {
	if wordLib == "" {
		log.Fatalln("No word list provided")
	}

	allWords := strings.Split(allWordLib, "\n")
	for _, word := range allWords {
		if len(word) < 5 {
			continue
		}
		allWordArr = append(allWordArr, word[:5])
	}

	words := strings.Split(wordLib, "\n")
	i := 0
	defer func() {
		err := recover()
		if err != nil {
			log.Fatalf("line %d, %s", i+1, err)
		}
	}()
	for ; i < len(words); i++ {
		word := words[i]
		if len(word) < 5 {
			continue
		}
		wordArr = append(wordArr, word[:5])
	}

	// seed once
	rand.Seed(time.Now().Unix())
}

type Wordle struct {
	sync.RWMutex

	Word    string                 `json:"word"`
	Tries   int32                  `json:"tries"`
	Guesses map[string][]*pb.Guess `json:"guesses"`
}

func NewWordle() *Wordle {
	w := &Wordle{}
	w.Tries = int32(rounds)
	w.Next()
	return w
}

func (w *Wordle) Guess(player, word string) (string, []*pb.Guess, error) {
	w.Lock()
	defer w.Unlock()

	guesses, ok := w.Guesses[player]
	if !ok {
		guesses = []*pb.Guess{}
	}
	if len(guesses) >= int(w.Tries) {
		return w.Word, guesses, errors.New("Reached guess limit")
	}

	if !isLetters(word) || len(word) != 5 || !in(word, allWordArr) {
		return "", guesses, errors.New("Invalid guess")
	}

	guess := &pb.Guess{Word: word}

	for i := 0; i < 5; i++ {
		char := new(pb.Char)
		char.Letter = string(word[i])
		char.Position = int32(i)

		switch {
		case word[i] == w.Word[i]:
			char.Correct = true
			char.InWord = true
			guess.Highlight += "[" + char.Letter + "]"
		case contains(word[i], w.Word):
			char.InWord = true
			guess.Highlight += "{" + char.Letter + "}"
		default:
			guess.Highlight += char.Letter
		}

		guess.Chars = append(guess.Chars, char)
	}

	guesses = append(guesses, guess)
	w.Guesses[player] = guesses

	if word == w.Word {
		return word, guesses, nil
	}

	return "", guesses, nil
}

func (w *Wordle) Next() {
	w.Lock()
	w.Word = wordArr[rand.Intn(len(wordArr))]
	w.Guesses = make(map[string][]*pb.Guess)
	w.Unlock()
}

func (w *Wordle) Load(b []byte) error {
	return json.Unmarshal(b, &w)
}

func (w *Wordle) Save() ([]byte, error) {
	return json.Marshal(w)
}

// determine if a word is in the word list
func in(str string, arr []string) bool {
	sort.Strings(arr)
	i := sort.SearchStrings(arr, str)
	if i < len(arr) && arr[i] == str {
		return true
	}
	return false
}

// determine if a word has the letter in it
func contains(letter uint8, word string) bool {
	for i := 0; i < len(word); i++ {
		if word[i] == letter {
			return true
		}
	}
	return false
}

// determine if the input is combined by letters
func isLetters(str string) bool {
	match, _ := regexp.MatchString(`^[A-Za-z]+$`, str)
	return match
}
