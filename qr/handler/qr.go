package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	pauth "github.com/micro/services/pkg/auth"
	adminpb "github.com/micro/services/pkg/service/proto"
	"github.com/micro/services/pkg/tenant"
	qr "github.com/micro/services/qr/proto"
	"github.com/skip2/go-qrcode"
)

const (
	prefixByTenant  = "qrByTenant"
	defaultCodeSize = 256
)

type QrCode struct {
	Filename string `json:"filename"`
	Created  int64  `json:"created"`
	Text     string `json:"text"`
}

type Qr struct {
	cdnPrefix string
}

func New() *Qr {
	v, err := config.Get("micro.qr.cdnprefix")
	if err != nil {
		log.Fatalf("Failed to get CDN prefix %s", err)
	}
	pref := v.String("")
	if len(pref) == 0 {
		log.Fatalf("Failed to get CDN prefix")
	}
	return &Qr{cdnPrefix: pref}
}

func namespacePrefix(tenantID string) string {
	return "micro/qr/" + tenantID
}

func (q *Qr) Generate(ctx context.Context, request *qr.GenerateRequest, response *qr.GenerateResponse) error {
	if len(request.Text) == 0 {
		return errors.BadRequest("qr.generate", "Missing parameter text")
	}
	ten, ok := tenant.FromContext(ctx)
	if !ok {
		log.Errorf("Error retrieving tenant")
		return errors.Unauthorized("qr.generate", "Unauthorized")
	}
	size := defaultCodeSize
	if request.Size > 0 {
		size = int(request.Size)
	}
	qrc, err := qrcode.Encode(request.Text, qrcode.Medium, size)
	if err != nil {
		log.Errorf("Error generating QR code %s", err)
		return errors.InternalServerError("qr.generate", "Error while generating QR code")
	}

	nsPrefix := namespacePrefix(ten)
	fileName := fmt.Sprintf("%s.png", uuid.New().String())
	if err := store.DefaultBlobStore.Write(
		fileName, bytes.NewBuffer(qrc),
		store.BlobContentType("image/png"),
		store.BlobPublic(true),
		store.BlobNamespace(nsPrefix)); err != nil {
		log.Errorf("Error saving QR code to blob store %s", err)
		return errors.InternalServerError("qr.generate", "Error while generating QR code")
	}
	// store record of it
	rec := QrCode{
		Filename: fileName,
		Created:  time.Now().Unix(),
		Text:     request.Text,
	}
	b, _ := json.Marshal(&rec)
	if err := store.Write(&store.Record{
		Key:   fmt.Sprintf("%s/%s/%s", prefixByTenant, nsPrefix, fileName),
		Value: b,
	}); err != nil {
		log.Errorf("Error saving QR code record %s", err)
		return errors.InternalServerError("qr.generate", "Error while generating QR code")
	}
	response.Qr = fmt.Sprintf("%s/%s/%s", q.cdnPrefix, nsPrefix, rec.Filename)
	return nil
}

func (q *Qr) DeleteData(ctx context.Context, request *adminpb.DeleteDataRequest, response *adminpb.DeleteDataResponse) error {
	method := "admin.DeleteData"
	_, err := pauth.VerifyMicroAdmin(ctx, method)
	if err != nil {
		return err
	}

	if len(request.TenantId) < 10 { // deliberate length check so we don't delete all the things
		return errors.BadRequest(method, "Missing tenant ID")
	}
	ns := namespacePrefix(request.TenantId)
	keys, err := store.DefaultBlobStore.List(store.BlobListNamespace(ns))
	if err != nil {
		return err
	}

	for _, key := range keys {
		err = store.DefaultBlobStore.Delete(key, store.BlobNamespace(ns))
		if err != nil {
			return err
		}
	}
	log.Infof("Deleted %d objects from S3 for %s", len(keys), request.TenantId)

	keys, err = store.List(store.ListPrefix(fmt.Sprintf("%s/%s/", prefixByTenant, ns)))
	if err != nil {
		return err
	}
	for _, key := range keys {
		if err := store.Delete(key); err != nil {
			return err
		}
	}

	log.Infof("Deleted %d objects from store for %s", len(keys), request.TenantId)

	return nil
}
