package handler

import (
	"context"
	"strings"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/users/proto"
	"gorm.io/gorm"
)

// Update a user
func (u *Users) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	// validate the request
	if len(req.Id) == 0 {
		return ErrMissingID
	}
	if req.FirstName != nil && len(req.FirstName.Value) == 0 {
		return ErrMissingFirstName
	}
	if req.LastName != nil && len(req.LastName.Value) == 0 {
		return ErrMissingLastName
	}
	if req.Email != nil && len(req.Email.Value) == 0 {
		return ErrMissingEmail
	}
	if req.Email != nil && !isEmailValid(req.Email.Value) {
		return ErrInvalidEmail
	}
	if req.Password != nil && len(req.Password.Value) < 8 {
		return ErrInvalidEmail
	}

	// lookup the user
	var user User
	if err := u.DB.Where(&User{ID: req.Id}).First(&user).Error; err == gorm.ErrRecordNotFound {
		return ErrNotFound
	} else if err != nil {
		logger.Errorf("Error reading from the database: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
	}

	// assign the updated values
	if req.FirstName != nil {
		user.FirstName = req.FirstName.Value
	}
	if req.LastName != nil {
		user.LastName = req.LastName.Value
	}
	if req.Email != nil {
		user.Email = strings.ToLower(req.Email.Value)
	}
	if req.Password != nil {
		p, err := hashAndSalt(req.Password.Value)
		if err != nil {
			logger.Errorf("Error hasing and salting password: %v", err)
			return errors.InternalServerError("HASHING_ERROR", "Error hashing password")
		}
		user.Password = p
	}

	// write the user to the database
	err := u.DB.Save(user).Error
	if err != nil && strings.Contains(err.Error(), "idx_users_email") {
		return ErrDuplicateEmail
	} else if err != nil {
		logger.Errorf("Error writing to the database: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
	}

	// serialize the user
	rsp.User = user.Serialize()
	return nil
}
