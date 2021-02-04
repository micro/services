package handler

import (
	"context"
	"strings"

	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/model"
	files "github.com/micro/services/files/proto"
)

type Files struct {
	db model.Model
}

func NewFiles() *Files {
	i := model.ByEquality("project")
	i.Order.Type = model.OrderTypeUnordered

	db := model.NewModel(
		model.WithIndexes(i),
	)

	db.Register(new(files.File))

	return &Files{
		db: db,
	}
}

func (e *Files) Save(ctx context.Context, req *files.SaveRequest, rsp *files.SaveResponse) error {
	log.Info("Received Files.Call request")
	for _, file := range req.Files {
		err := e.db.Create(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Files) List(ctx context.Context, req *files.ListRequest, rsp *files.ListResponse) error {
	log.Info("Received Files.Call request")
	rsp.Files = []*files.File{}
	err := e.db.Read(model.QueryEquals("project", req.GetProject()), &rsp.Files)
	if err != nil {
		return err
	}
	// @todo funnily while this is the archetypical
	// query for the KV store interface, it's not supported by the model
	// so we do client side filtering here
	if req.Path != "" {
		filtered := []*files.File{}
		for _, file := range rsp.Files {
			if strings.HasPrefix(file.Path, req.Path) {
				filtered = append(filtered, file)
			}
		}
		rsp.Files = filtered
	}
	return nil
}
