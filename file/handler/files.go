package handler

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/model"
	file "github.com/micro/services/file/proto"
	"github.com/micro/services/pkg/tenant"
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
			Key:     "Path",
			Indexes: []model.Index{i},
		},
	)

	return &File{
		db: db,
	}
}

func (e *File) Delete(ctx context.Context, req *file.DeleteRequest, rsp *file.DeleteResponse) error {
	if len(req.Path) == 0 {
		return errors.BadRequest("file.read", "missing file path")
	}

	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		tenantId = "micro"
	}

	path := filepath.Join(tenantId, req.Project, req.Path)
	project := tenantId + "/" + req.Project

	// delete one file
	if !strings.HasSuffix(req.Path, "/") {
		return e.db.Delete(model.QueryEquals("Path", path))
	}

	var files []*file.Record

	// read all the files for the project
	err := e.db.Read(model.QueryEquals("project", project), &files)
	if err != nil {
		return err
	}

	for _, file := range files {
		// delete a list of files
		if file.Project != project {
			continue
		}
		if !strings.HasPrefix(file.Path, path) {
			continue
		}
		// delete the file
		e.db.Delete(model.QueryEquals("Path", file.Path))
	}

	return nil
}

func (e *File) Read(ctx context.Context, req *file.ReadRequest, rsp *file.ReadResponse) error {
	log.Info("Received File.Read request")

	if len(req.Path) == 0 {
		return errors.BadRequest("file.read", "missing file path")
	}

	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		tenantId = "micro"
	}

	var files []*file.Record

	project := tenantId + "/" + req.Project

	// read all the files for the project
	err := e.db.Read(model.QueryEquals("project", project), &files)
	if err != nil {
		return err
	}

	// filter the file
	for _, file := range files {
		// check project matches tenants
		if file.Project != project {
			continue
		}

		// strip the tenant id
		file.Project = strings.TrimPrefix(file.Project, tenantId+"/")
		file.Path = strings.TrimPrefix(file.Path, filepath.Join(tenantId, req.Project))

		// check the path matches the request
		if req.Path == file.Path {
			rsp.File = file
		}
	}

	return nil
}

func (e *File) Save(ctx context.Context, req *file.SaveRequest, rsp *file.SaveResponse) error {
	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		tenantId = "micro"
	}

	log.Info("Received File.Save request")

	// prefix the tenant
	req.File.Project = filepath.Join(tenantId, req.File.Project)
	req.File.Path = filepath.Join(req.File.Project, req.File.Path)

	if len(req.File.Created) == 0 {
		req.File.Created = time.Now().Format(time.RFC3339Nano)
	}

	// set updated time
	req.File.Updated = time.Now().Format(time.RFC3339Nano)

	// create the file
	err := e.db.Create(req.File)
	if err != nil {
		return err
	}

	return nil
}

func (e *File) BatchSave(ctx context.Context, req *file.BatchSaveRequest, rsp *file.BatchSaveResponse) error {
	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		tenantId = "micro"
	}

	log.Info("Received File.BatchSave request")

	for _, reqFile := range req.Files {
		// prefix the tenant
		reqFile.Project = filepath.Join(tenantId, reqFile.Project)
		reqFile.Path = filepath.Join(reqFile.Project, reqFile.Project)

		// create the file
		err := e.db.Create(reqFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *File) List(ctx context.Context, req *file.ListRequest, rsp *file.ListResponse) error {
	log.Info("Received File.List request")

	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		tenantId = "micro"
	}

	// prefix tenant id
	project := tenantId + "/" + req.Project

	var files []*file.Record

	// read all the files for the project
	if err := e.db.Read(model.QueryEquals("project", project), &files); err != nil {
		return err
	}

	// @todo funnily while this is the archetypical
	// query for the KV store interface, it's not supported by the model
	// so we do client side filtering here
	for _, file := range files {
		if file.Project != project {
			continue
		}

		// strip the prefixes
		file.Project = strings.TrimPrefix(file.Project, tenantId+"/")
		file.Path = strings.TrimPrefix(file.Path, filepath.Join(tenantId, req.Project))

		// strip the file contents
		// no file listing ever contains it
		file.Content = ""

		// if requesting all files or path matches
		if req.Path != "" {
			if strings.HasPrefix(file.Path, req.Path) {
				rsp.Files = append(rsp.Files, file)
			}
		} else {
			rsp.Files = append(rsp.Files, file)
		}
	}

	return nil
}
