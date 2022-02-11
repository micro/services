package handler

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	file "github.com/micro/services/file/proto"
	pauth "github.com/micro/services/pkg/auth"
	adminpb "github.com/micro/services/pkg/service/proto"
	"github.com/micro/services/pkg/tenant"
)

const pathPrefix = "files"
const hostPrefix = "https://cdn.m3ocontent.com"

func NewFile() *File {
	var hp string
	cfg, err := config.Get("micro.image.host_prefix")
	if err != nil {
		hp = cfg.String(hostPrefix)
	}
	if len(strings.TrimSpace(hp)) == 0 {
		hp = hostPrefix
	}
	return &File{
		hostPrefix: hp,
	}
}

type File struct {
	hostPrefix string
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
	if err := store.Write(store.NewRecord(path, req.File)); err != nil {
		return err
	}

	// save the file if made public
	if !req.Public {
		return nil
	}

	path = filepath.Join(pathPrefix, tenantId, req.File.Project, req.File.Path)

	// upload to the blob store
	err := store.DefaultBlobStore.Write(path, bytes.NewReader([]byte(req.File.Content)), store.BlobPublic(true))
	if err != nil {
		return err
	}

	rsp.Url = fmt.Sprintf("%v/%v/%v", e.hostPrefix, "micro", path)

	return nil
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

func (e *File) DeleteData(ctx context.Context, request *adminpb.DeleteDataRequest, response *adminpb.DeleteDataResponse) error {
	method := "admin.DeleteData"
	_, err := pauth.VerifyMicroAdmin(ctx, method)
	if err != nil {
		return err
	}

	if len(request.TenantId) < 10 { // deliberate length check so we don't delete all the things
		return errors.BadRequest(method, "Missing tenant ID")
	}

	path := filepath.Join("file", request.TenantId)

	// read all the files for the project
	records, err := store.List(store.ListPrefix(path))
	if err != nil {
		return err
	}

	for _, file := range records {
		if err := store.Delete(file); err != nil {
			return err
		}
	}
	log.Infof("Deleted %d records for %s", len(records), request.TenantId)

	return nil
}
