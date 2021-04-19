package handler

import (
	"context"
	"regexp"
	"strings"

	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/users/proto"
	"gorm.io/gorm"
)

// Update a user
func (u *Users) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	_, ok := auth.AccountFromContext(ctx)
	if !ok {
		errors.Unauthorized("UNAUTHORIZED", "Unauthorized")
	}
	// validate the request
	if len(req.Id) == 0 {
		return ErrMissingID
	}
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
		return ErrInvalidEmail
	}

	// lookup the user
	var user User
	db, err := u.GetDBConn(ctx)
	if err != nil {
		logger.Errorf("Error connecting to DB: %v", err)
		return errors.InternalServerError("DB_ERROR", "Error connecting to DB")
	}
	if err := db.Where(&User{ID: req.Id}).First(&user).Error; err == gorm.ErrRecordNotFound {
		return ErrNotFound
	} else if err != nil {
		logger.Errorf("Error reading from the database: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
	}

	// assign the updated values
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Email != "" {
		user.Email = strings.ToLower(req.Email)
	}
	if req.Password != "" {
		p, err := hashAndSalt(req.Password)
		if err != nil {
			logger.Errorf("Error hasing and salting password: %v", err)
			return errors.InternalServerError("HASHING_ERROR", "Error hashing password")
		}
		user.Password = p
	}

	// write the user to the database
	err = db.Save(user).Error
	if err != nil {
		if match, _ := regexp.MatchString(`idx_[\S]+_users_email`, err.Error()); match {
			return ErrDuplicateEmail
		}
		logger.Errorf("Error writing to the database: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
	}

	// serialize the user
	rsp.User = user.Serialize()
	return nil
}
