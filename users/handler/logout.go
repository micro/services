package handler

import (
	"context"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/users/proto"
	"gorm.io/gorm"
)

// Logout expires all tokens for the user
func (u *Users) Logout(ctx context.Context, req *pb.LogoutRequest, rsp *pb.LogoutResponse) error {
	// validate the request
	if len(req.Id) == 0 {
		return ErrMissingID
	}

	return u.DB.Transaction(func(tx *gorm.DB) error {
		// lookup the user
		var user User
		if err := tx.Where(&User{ID: req.Id}).Preload("Tokens").First(&user).Error; err == gorm.ErrRecordNotFound {
			return ErrNotFound
		} else if err != nil {
			logger.Errorf("Error reading from the database: %v", err)
			return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
		}

		// delete the tokens
		if err := tx.Delete(user.Tokens).Error; err != nil {
			logger.Errorf("Error deleting from the database: %v", err)
			return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
		}

		return nil
	})
}
