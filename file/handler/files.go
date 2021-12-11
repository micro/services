package handler

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	file "github.com/micro/services/file/proto"
	"github.com/micro/services/pkg/tenant"
)

type File struct{}

func NewFile() *File {
	return &File{}
}

func (e *File) Delete(ctx context.Context, req *file.DeleteRequest, rsp *file.DeleteResponse) error {
	if len(req.Path) == 0 {
		return errors.BadRequest("file.read", "missing file path")
	}

	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		tenantId = "micro"
	}

	path := filepath.Join("file", tenantId, req.Project, req.Path)

	// delete one file
	if !strings.HasSuffix(req.Path, "/") {
		return store.Delete(path)
	}

	// read all the files for the project
	records, err := store.List(store.ListPrefix(path))
	if err != nil {
		return err
	}

	for _, file := range records {
		store.Delete(file)
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

	path := filepath.Join("file", tenantId, req.Project, req.Path)

	records, err := store.Read(path)
	if err != nil {
		return err
	}

	if len(records) == 0 {
		return errors.NotFound("file.read", "file not found")
	}

	// filter the file
	rec := records[0]
	file := new(file.Record)

	if err := rec.Decode(file); err != nil {
		return err
	}

	// strip the tenant id
	file.Project = strings.TrimPrefix(file.Project, tenantId+"/")
	file.Path = strings.TrimPrefix(file.Path, filepath.Join(tenantId, req.Project))

	rsp.File = file

	return nil
}

func (e *File) Save(ctx context.Context, req *file.SaveRequest, rsp *file.SaveResponse) error {
	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		tenantId = "micro"
	}

	if req.File == nil {
		return errors.BadRequest("file.save", "missing file")
	}

	log.Info("Received File.Save request")

	path := filepath.Join("file", tenantId, req.File.Project, req.File.Path)

	if len(req.File.Created) == 0 {
		req.File.Created = time.Now().Format(time.RFC3339Nano)
	}

	// set updated time
	req.File.Updated = time.Now().Format(time.RFC3339Nano)

	// create the file
	return store.Write(store.NewRecord(path, req.File))
}

func (e *File) List(ctx context.Context, req *file.ListRequest, rsp *file.ListResponse) error {
	log.Info("Received File.List request")

	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		tenantId = "micro"
	}

	// prefix tenant id
	path := filepath.Join("file", tenantId, req.Project, req.Path)

	records, err := store.Read(path, store.ReadPrefix())
	if err != nil {
		return err
	}

	for _, rec := range records {
		file := new(file.Record)

		if err := rec.Decode(file); err != nil {
			continue
		}

		// strip the prefixes
		file.Project = strings.TrimPrefix(file.Project, tenantId+"/")
		file.Path = strings.TrimPrefix(file.Path, filepath.Join(tenantId, req.Project))

		// strip the file contents
		// no file listing ever contains it
		file.Content = ""

		rsp.Files = append(rsp.Files, file)
	}

	return nil
}
