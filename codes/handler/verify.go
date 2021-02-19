package handler

import (
	"context"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/codes/proto"
	"gorm.io/gorm"
)

func (c *Codes) Verify(ctx context.Context, req *pb.VerifyRequest, rsp *pb.VerifyResponse) error {
	// validate the request
	if len(req.Code) == 0 {
		return ErrMissingCode
	}
	if len(req.Identity) == 0 {
		return ErrMissingIdentity
	}

	// lookup the code
	var code Code
	if err := c.DB.Where(&Code{Code: req.Code, Identity: req.Identity}).First(&code).Error; err == gorm.ErrRecordNotFound {
		return ErrInvalidCode
	} else if err != nil {
		logger.Errorf("Error reading code from database: %v", err)
		return errors.InternalServerError("DATABASE_ERORR", "Error connecting to database")
	}

	// check the invite hasn't expired
	if code.ExpiresAt.Before(c.Time()) {
		return ErrExpiredCode
	}

	return nil
}
