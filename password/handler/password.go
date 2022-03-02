package handler

import (
	"context"
	"crypto/rand"

	pb "github.com/micro/services/password/proto"
)


var (
	minLength = 8

	special = "!@#$%&*"
	numbers = "0123456789"
	lowercase = "abcdefghijklmnopqrstuvwxyz"
	uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	allChars = special + numbers + lowercase + uppercase
)

func random(chars string, i int) string {
        bytes := make([]byte, i)

        for {
                rand.Read(bytes)
                for i, b := range bytes {
                        bytes[i] = chars[b%byte(len(chars))]
                }
		break
        }

        return string(bytes)
}


type Password struct{}

func (p *Password) Generate(ctx context.Context, req *pb.GenerateRequest, rsp *pb.GenerateResponse) error {
	if req.Length <= 0 {
		req.Length = int32(minLength)
	}

	// TODO; allow user to specify types of characters

	// generate and return password
	rsp.Password = random(allChars, int(req.Length))

	return nil
}
