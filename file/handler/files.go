package handler

import (
	"context"
	"errors"
	"strings"

	"github.com/micro/micro/v3/service/auth"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/model"
	file "github.com/micro/services/file/proto"
)

type File struct {
	db model.Model
}

func NewFile() *File {
	i := model.ByEquality("project")
	i.Order.Type = model.OrderTypeUnordered

	db := model.New(
		file.Record{},
		&model.Options{
			Key:     "Id",
			Indexes: []model.Index{i},
		},
	)

	return &File{
		db: db,
	}
}

func (e *File) Save(ctx context.Context, req *file.SaveRequest, rsp *file.SaveResponse) error {
	// @todo return proper micro errors
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.New("File.Save requires authentication")
	}

	log.Info("Received File.Save request")
	for _, reqFile := range req.Files {
		f := file.Record{}
		err := e.db.Read(model.QueryEquals("Id", reqFile.Id), &f)
		if err != nil && err != model.ErrorNotFound {
			return err
		}
		// if file exists check ownership
		if f.Id != "" && f.Owner != acc.ID {
			return errors.New("Not authorized")
		}
		err = e.db.Create(reqFile)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *File) List(ctx context.Context, req *file.ListRequest, rsp *file.ListResponse) error {
	log.Info("Received File.List request")
	rsp.Files = []*file.Record{}
	err := e.db.Read(model.QueryEquals("project", req.GetProject()), &rsp.Files)
	if err != nil {
		return err
	}
	// @todo funnily while this is the archetypical
	// query for the KV store interface, it's not supported by the model
	// so we do client side filtering here
	if req.Path != "" {
		filtered := []*file.Record{}
		for _, file := range rsp.Files {
			if strings.HasPrefix(file.Path, req.Path) {
				filtered = append(filtered, file)
			}
		}
		rsp.Files = filtered
	}
	return nil
}
