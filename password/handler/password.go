package handler

import (
	"context"
	"crypto/rand"
	"strings"

	pb "github.com/micro/services/password/proto"
)

const (
	minLength = 8

	special   = "!@#$%&*"
	numbers   = "0123456789"
	lowercase = "abcdefghijklmnopqrstuvwxyz"
	uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	allChars  = special + numbers + lowercase + uppercase
)

func random(chars string, i int) string {
	bytes := make([]byte, i)

	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = chars[b%byte(len(chars))]
	}

	return string(bytes)
}

type Password struct{}

func (p *Password) Generate(ctx context.Context, req *pb.GenerateRequest, rsp *pb.GenerateResponse) error {
	if req.Length <= 0 {
		req.Length = int32(minLength)
	}

	charSpace := ""
	if req.Numbers {
		charSpace += numbers
	}
	if req.Lowercase {
		charSpace += lowercase
	}
	if req.Uppercase {
		charSpace += uppercase
	}
	if req.Special {
		charSpace += special
	}
	if len(charSpace) == 0 {
		charSpace = allChars
	}

	for {
		// generate and return password
		rsp.Password = random(charSpace, int(req.Length))
		// validate we have minimums needed
		reqsSatisfied := true
		if req.Numbers {
			reqsSatisfied = reqsSatisfied && strings.ContainsAny(rsp.Password, numbers)
		}
		if req.Lowercase {
			reqsSatisfied = reqsSatisfied && strings.ContainsAny(rsp.Password, lowercase)
		}
		if req.Uppercase {
			reqsSatisfied = reqsSatisfied && strings.ContainsAny(rsp.Password, uppercase)
		}
		if req.Special {
			reqsSatisfied = reqsSatisfied && strings.ContainsAny(rsp.Password, special)
		}
		if reqsSatisfied {
			break
		}
		// failed to satisfy all reqs, try again
	}

	return nil
}
