package handler

import (
	"context"
	"strings"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/users/proto"
)

// Read users using email
func (u *Users) ReadByEmail(ctx context.Context, req *pb.ReadByEmailRequest, rsp *pb.ReadByEmailResponse) error {
	// validate the request
	if len(req.Emails) == 0 {
		return ErrMissingEmails
	}
	emails := make([]string, len(req.Emails))
	for i, e := range req.Emails {
		emails[i] = strings.ToLower(e)
	}

	// query the database
	db, err := u.getDBConn(ctx)
	if err != nil {
		logger.Errorf("Error connecting to DB: %v", err)
		return errors.InternalServerError("DB_ERROR", "Error connecting to DB")
	}
	var users []User
	if err := db.Model(&User{}).Where("lower(email) IN (?)", emails).Find(&users).Error; err != nil {
		logger.Errorf("Error reading from the database: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
	}

	// serialize the response
	rsp.Users = make(map[string]*pb.User, len(users))
	for _, u := range users {
		rsp.Users[u.Email] = u.Serialize()
	}
	return nil
}
