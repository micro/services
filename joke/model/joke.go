package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/micro/micro/v3/service/logger"
	"github.com/pkg/errors"
)

var (
	// jokes is a memory store to save all jokes
	jokes []Joke
)

// GetAllJokes return all jokes
func GetAllJokes() []Joke {
	return jokes
}

// JokeSource is the source of joke, contains source name and source api
// The Api response content must be an array, and every object has contains `title` or `body` field
// eg:
// [
//    {
//       "title": "joke's title"
//       "body": "joke's body"
//    }
// ]
type JokeSource struct {
	Source string
	Api    string
}

type Joke struct {
	Id       string
	Title    string
	Body     string
	Category string
	Source   string
}

// GetKey is used to get jokes uniq key
func (i Joke) GetKey() string {
	return fmt.Sprintf("%s-%s", i.Source, i.Id)
}

type joke struct {
	source JokeSource
}

func NewJoke(s JokeSource) *joke {
	return &joke{
		source: s,
	}
}

// get is used to get jokes from api
func (j *joke) get() (jokes []Joke, err error) {

	client := http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Get(j.source.Api)
	if err != nil {
		return nil, errors.Wrap(err, "request jokes api error: "+err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrap(err, "request jokes api status is not ok: "+resp.Status)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.Errorf("close response error, %v", err)
			return
		}
	}()

	b, _ := ioutil.ReadAll(resp.Body)

	results := make([]map[string]interface{}, 0)
	if err := json.Unmarshal(b, &results); err != nil {
		return nil, errors.Wrap(err, "json unmarshal jokes api error: "+err.Error())
	}

	for _, r := range results {
		info := Joke{}

		info.Id = func() string {
			if id, ok := r["id"].(string); ok {
				return id
			} else if id, ok := r["id"].(float64); ok {
				return fmt.Sprintf("%d", int32(id))
			} else {
				return "unknown"
			}
		}()

		info.Source = j.source.Source
		info.Title, _ = r["title"].(string)
		info.Body, _ = r["body"].(string)
		info.Category, _ = r["category"].(string)

		if info.Title == "" && info.Body == "" {
			continue
		}

		jokes = append(jokes, info)
	}

	return jokes, nil
}

// Load is used to save jokes in memory
func (j *joke) Load() (err error) {
	if j.source.Source == "" || j.source.Api == "" {
		return errors.New("source or api can not be empty")
	}

	js, err := j.get()
	if err != nil {
		return err
	}

	jokes = append(jokes, js...)

	return nil
}
