package handler

import (
	"context"

	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/users/proto"
	"gorm.io/gorm"
)

// Validate a token, each time a token is validated it extends its lifetime for another week
func (u *Users) Validate(ctx context.Context, req *pb.ValidateRequest, rsp *pb.ValidateResponse) error {
	_, ok := auth.AccountFromContext(ctx)
	if !ok {
		errors.Unauthorized("UNAUTHORIZED", "Unauthorized")
	}
	// validate the request
	if len(req.Token) == 0 {
		return ErrMissingToken
	}

	db, err := u.getDBConn(ctx)
	if err != nil {
		logger.Errorf("Error connecting to DB: %v", err)
		return errors.InternalServerError("DB_ERROR", "Error connecting to DB")
	}
	return db.Transaction(func(tx *gorm.DB) error {
		// lookup the token
		var token Token
		if err := tx.Where(&Token{Key: req.Token}).Preload("User").First(&token).Error; err == gorm.ErrRecordNotFound {
			return ErrInvalidToken
		} else if err != nil {
			logger.Errorf("Error reading from the database: %v", err)
			return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
		}

		// ensure the token is valid
		if u.Time().After(token.ExpiresAt) {
			return ErrTokenExpired
		}

		// extend the token for another lifetime
		token.ExpiresAt = u.Time().Add(tokenTTL)
		if err := tx.Save(&token).Error; err != nil {
			logger.Errorf("Error writing to the database: %v", err)
			return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
		}

		// serialize the response
		rsp.User = token.User.Serialize()
		return nil
	})
}
