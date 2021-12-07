package handler

import (
	"context"
	"math/rand"

	"github.com/micro/services/joke/model"
	pb "github.com/micro/services/joke/proto"
)

type Joke struct{}

// Random is used to get random jokes
func (e *Joke) Random(_ context.Context, req *pb.RandomRequest, rsp *pb.RandomResponse) error {
	jokes := model.GetAllJokes()
	count := req.Count
	
	if count <= 0 {
		count = 1
	} else if count > 10 {
		count = 10
	} 
	
	if count > len(jokes) {
		count = len(jokes)
	}

	for i := int32(0); i < count; i++ {
		random := jokes[rand.Intn(len(jokes))]

		info := &pb.JokeInfo{
			Id:       random.Id,
			Title:    random.Title,
			Body:     random.Body,
			Category: random.Category,
			Source:   random.Source,
		}

		rsp.Jokes = append(rsp.Jokes, info)
	}

	return nil
}
