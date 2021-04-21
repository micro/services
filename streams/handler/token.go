package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/streams/proto"
)

func (s *Streams) Token(ctx context.Context, req *pb.TokenRequest, rsp *pb.TokenResponse) error {
	_, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized("UNAUTHORIZED", "Unauthorized")
	}
	if len(req.Topic) > 0 {
		if err := validateTopicInput(req.Topic); err != nil {
			return err
		}
	}

	// construct the token and write it to the database
	t := Token{
		Token:     uuid.New().String(),
		ExpiresAt: s.Time().Add(TokenTTL),
		Topic:     req.Topic,
	}
	dbConn, err := s.GetDBConn(ctx)
	if err != nil {
		logger.Errorf("Error creating token in store: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error writing token to database")
	}
	if err := dbConn.Create(&t).Error; err != nil {
		logger.Errorf("Error creating token in store: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error writing token to database")
	}

	// return the token in the response
	rsp.Token = t.Token
	return nil
}
