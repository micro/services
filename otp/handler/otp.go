package handler

import (
	"context"
	"time"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/otp/proto"
	"github.com/micro/services/pkg/cache"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type Otp struct{}

func (e *Otp) Generate(ctx context.Context, req *pb.GenerateRequest, rsp *pb.GenerateResponse) error {
	if len(req.Id) == 0 {
		return errors.BadRequest("otp.generate", "missing id")
	}

	// check if a key exists for the user
	var secret string

	if err := cache.Context(ctx).Get(req.Id, &secret); err != nil {
		// generate a key
		key, err := totp.Generate(totp.GenerateOpts{
			Issuer:      "Micro",
			AccountName: req.Id,
			Period:      60,
			Algorithm:   otp.AlgorithmSHA1,
		})
		if err != nil {
			logger.Error("Failed to generate secret: %v", err)
			return errors.InternalServerError("otp.generate", "failed to generate code")
		}

		secret = key.Secret()

		if err := cache.Context(ctx).Set(req.Id, secret, time.Now().Add(time.Minute*5)); err != nil {
			logger.Error("Failed to store secret: %v", err)
			return errors.InternalServerError("otp.generate", "failed to generate code")
		}
	}

	// generate a new code
	code, err := totp.GenerateCodeCustom(secret, time.Now(), totp.ValidateOpts{
		Period:    60,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})

	if err != nil {
		return errors.InternalServerError("otp.generate", "failed to generate code: %v", err)
	}

	// return the code
	rsp.Code = code

	return nil
}

func (e *Otp) Validate(ctx context.Context, req *pb.ValidateRequest, rsp *pb.ValidateResponse) error {
	if len(req.Id) == 0 {
		return errors.BadRequest("otp.generate", "missing id")
	}
	if len(req.Code) == 0 {
		return errors.BadRequest("otp.generate", "missing code")
	}

	var secret string

	if err := cache.Context(ctx).Get(req.Id, &secret); err != nil {
		logger.Error("Failed to get secret from store: %v", err)
		return errors.InternalServerError("otp.generate", "failed to validate code")
	}

	ok, err := totp.ValidateCustom(req.Code, secret, time.Now(), totp.ValidateOpts{
		Period:    60,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil {
		return errors.InternalServerError("otp.generate", "failed to validate code")
	}

	// set the response
	rsp.Success = ok

	return nil
}
