package handler

import (
	"context"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/users/proto"
)

// List all users
func (u *Users) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	// query the database
	var users []User
	if err := u.DB.Model(&User{}).Find(&users).Error; err != nil {
		logger.Errorf("Error reading from the database: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
	}

	// serialize the response
	rsp.Users = make([]*pb.User, len(users))
	for i, u := range users {
		rsp.Users[i] = u.Serialize()
	}
	return nil
}
