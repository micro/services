package handler

import (
	"context"

	db "github.com/micro/services/db/proto"
	gorm2 "github.com/micro/services/pkg/gorm"
)

type Db struct {
	gorm2.Helper
}

// Call is a single request handler called via client.Call or the generated client code
func (e *Db) Create(ctx context.Context, req *db.CreateRequest, rsp *db.CreateResponse) error {

	return nil
}
