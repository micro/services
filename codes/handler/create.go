package handler

import (
	"context"
	"math/rand"
	"strconv"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/codes/proto"
)

func (c *Codes) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	// validate the request
	if len(req.Identity) == 0 {
		return ErrMissingIdentity
	}

	// construct the code
	code := Code{Code: generateCode(), Identity: req.Identity}
	if req.ExpiresAt != nil {
		code.ExpiresAt = req.ExpiresAt.AsTime()
	} else {
		code.ExpiresAt = c.Time().Add(DefaultTTL)
	}

	// write to the database
	if err := c.DB.Create(&code).Error; err != nil {
		logger.Errorf("Error creating code in database: %v", err)
		return errors.InternalServerError("DATABASE_ERORR", "Error connecting to database")
	}

	// return the code
	rsp.Code = code.Code
	return nil
}

// generateCode generates a random 8 digit code
func generateCode() string {
	v := rand.Intn(89999999) + 10000000
	return strconv.Itoa(v)
}
