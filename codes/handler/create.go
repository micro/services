package handler

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/codes/proto"
)

func (c *Codes) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	_, ok := auth.AccountFromContext(ctx)
	if !ok {
		errors.Unauthorized("UNAUTHORIZED", "Unauthorized")
	}
	// validate the request
	if len(req.Identity) == 0 {
		return ErrMissingIdentity
	}

	// construct the code
	code := Code{Code: generateCode(), Identity: req.Identity}
	if req.ExpiresAt == 0 {
		code.ExpiresAt = time.Unix(req.ExpiresAt, 0)
	} else {
		code.ExpiresAt = c.Time().Add(DefaultTTL)
	}

	db, err := c.GetDBConn(ctx)
	if err != nil {
		logger.Errorf("Error connecting to DB: %v", err)
		return errors.InternalServerError("DB_ERROR", "Error connecting to DB")
	}
	// write to the database
	if err := db.Create(&code).Error; err != nil {
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
