package handler

import (
	"context"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/users/proto"
	"gorm.io/gorm"
)

// Delete a user
func (u *Users) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	// validate the request
	if len(req.Id) == 0 {
		return ErrMissingID
	}

	// delete the users tokens
	return u.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&Token{}, &Token{UserID: req.Id}).Error; err != nil {
			logger.Errorf("Error writing to the database: %v", err)
			return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
		}

		// delete from the database
		if err := tx.Delete(&User{}, &User{ID: req.Id}).Error; err != nil {
			logger.Errorf("Error writing to the database: %v", err)
			return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
		}

		return nil
	})
}
