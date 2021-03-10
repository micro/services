package handler

import (
	"context"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/users/proto"
)

// Read users using ID
func (u *Users) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	// validate the request
	if len(req.Ids) == 0 {
		return ErrMissingIDs
	}

	// query the database
	var users []User
	if err := u.DB.Model(&User{}).Where("id IN (?)", req.Ids).Find(&users).Error; err != nil {
		logger.Errorf("Error reading from the database: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
	}

	// serialize the response
	rsp.Users = make(map[string]*pb.User, len(users))
	for _, u := range users {
		rsp.Users[u.ID] = u.Serialize()
	}
	return nil
}
