package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/users/proto"
	"gorm.io/gorm"
)

// Login using email and password returns the users profile and a token
func (u *Users) Login(ctx context.Context, req *pb.LoginRequest, rsp *pb.LoginResponse) error {
	// validate the request
	if len(req.Email) == 0 {
		return ErrMissingEmail
	}
	if len(req.Password) == 0 {
		return ErrInvalidPassword
	}

	db, err := u.getDBConn(ctx)
	if err != nil {
		logger.Errorf("Error connecting to DB: %v", err)
		return errors.InternalServerError("DB_ERROR", "Error connecting to DB")
	}

	return db.Transaction(func(tx *gorm.DB) error {
		// lookup the user
		var user User
		if err := tx.Where(&User{Email: req.Email}).First(&user).Error; err == gorm.ErrRecordNotFound {
			return ErrNotFound
		} else if err != nil {
			logger.Errorf("Error reading from the database: %v", err)
			return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
		}

		// compare the passwords
		if !passwordsMatch(user.Password, req.Password) {
			return ErrIncorrectPassword
		}

		// generate a token for the user
		token := Token{
			UserID:    user.ID,
			Key:       uuid.New().String(),
			ExpiresAt: u.Time().Add(tokenTTL),
		}
		if err := tx.Create(&token).Error; err != nil {
			logger.Errorf("Error writing to the database: %v", err)
			return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
		}

		// serialize the response
		rsp.Token = token.Key
		rsp.User = user.Serialize()
		return nil
	})
}
