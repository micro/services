package handler

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/teris-io/shortid"
	"github.com/google/uuid"
	"github.com/mattheath/kala/bigflake"
	"github.com/mattheath/kala/snowflake"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/errors"
	pb "github.com/micro/services/id/proto"
)

type Id struct {
	Snowflake *snowflake.Snowflake
	Bigflake  *bigflake.Bigflake
}

func New() *Id {
	id := rand.Intn(100)

	sf, err := snowflake.New(uint32(id))
	if err != nil {
		panic(err.Error())
	}
	bg, err := bigflake.New(uint64(id))
	if err != nil {
		panic(err.Error())
	}

	return &Id{
		Snowflake: sf,
		Bigflake:  bg,
	}
}

func (id *Id) Generate(ctx context.Context, req *pb.GenerateRequest, rsp *pb.GenerateResponse) error {
	if len(req.Type) == 0 {
		req.Type = "uuid"
	}

	switch req.Type {
	case "uuid":
		rsp.Type = "uuid"
		rsp.Id = uuid.New().String()
	case "snowflake":
		id, err := id.Snowflake.Mint()
		if err != nil {
			logger.Errorf("Failed to generate snowflake id: %v", err)
			return errors.InternalServerError("id.generate", "failed to mint snowflake id")
		}
		rsp.Type = "snowflake"
		rsp.Id = fmt.Sprintf("%v", id)
	case "bigflake":
		id, err := id.Bigflake.Mint()
		if err != nil {
			logger.Errorf("Failed to generate bigflake id: %v", err)
			return errors.InternalServerError("id.generate", "failed to mint bigflake id")
		}
		rsp.Type = "bigflake"
		rsp.Id = fmt.Sprintf("%v", id)
	case "shortid":
		id, err := shortid.Generate()
		if err != nil {
			logger.Errorf("Failed to generate shortid id: %v", err)
			return errors.InternalServerError("id.generate", "failed to generate short id")
		}
		rsp.Type = "shortid"
		rsp.Id = id
	default:
		return errors.BadRequest("id.generate", "unsupported id type")
	}

	return nil
}
