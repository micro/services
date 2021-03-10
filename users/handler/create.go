package handler

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/users/proto"
	"gorm.io/gorm"
)

// Create a user
func (u *Users) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	// validate the request
	if len(req.FirstName) == 0 {
		return ErrMissingFirstName
	}
	if len(req.LastName) == 0 {
		return ErrMissingLastName
	}
	if len(req.Email) == 0 {
		return ErrMissingEmail
	}
	if !isEmailValid(req.Email) {
		return ErrInvalidEmail
	}
	if len(req.Password) < 8 {
		return ErrInvalidPassword
	}

	// hash and salt the password using bcrypt
	phash, err := hashAndSalt(req.Password)
	if err != nil {
		logger.Errorf("Error hasing and salting password: %v", err)
		return errors.InternalServerError("HASHING_ERROR", "Error hashing password")
	}

	return u.DB.Transaction(func(tx *gorm.DB) error {
		// write the user to the database
		user := &User{
			ID:        uuid.New().String(),
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     strings.ToLower(req.Email),
			Password:  phash,
		}
		err = u.DB.Create(user).Error
		if err != nil && strings.Contains(err.Error(), "idx_users_email") {
			return ErrDuplicateEmail
		} else if err != nil {
			logger.Errorf("Error writing to the database: %v", err)
			return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
		}

		// generate a token for the user
		token := Token{
			UserID:    user.ID,
			Key:       uuid.New().String(),
			ExpiresAt: u.Time().Add(time.Hour * 24 * 7),
		}
		if err := tx.Create(&token).Error; err != nil {
			logger.Errorf("Error writing to the database: %v", err)
			return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
		}

		// serialize the response
		rsp.User = user.Serialize()
		rsp.Token = token.Key
		return nil
	})
}
