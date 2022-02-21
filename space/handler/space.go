package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/micro/micro/v3/proto/api"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	pauth "github.com/micro/services/pkg/auth"
	adminpb "github.com/micro/services/pkg/service/proto"
	"github.com/micro/services/pkg/tenant"
	pb "github.com/micro/services/space/proto"
	"github.com/minio/minio-go/v7/pkg/s3utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awscreds "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	sthree "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

const (
	mdACL       = "X-Amz-Acl"
	mdACLPublic = "public-read"

	visibilityPrivate = "private"
	visibilityPublic  = "public"

	prefixByUser = "byUser"

	// max read 5mb for small objects
	maxReadSize = 5 * 1024 * 1024
)

type Space struct {
	conf   conf
	client s3iface.S3API
}

type conf struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Endpoint  string `json:"endpoint"`
	SpaceName string `json:"space_name"`
	SSL       bool   `json:"ssl"`
	Region    string `json:"region"`
	BaseURL   string `json:"base_url"`
}

type meta struct {
	Visibility   string
	CreateTime   string
	ModifiedTime string
}

func NewSpace(srv *service.Service) *Space {
	var c conf
	val, err := config.Get("micro.space")
	if err != nil {
		log.Fatalf("Failed to load config %s", err)
	}
	if err := val.Scan(&c); err != nil {
		log.Fatalf("Failed to load config %s", err)
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint:    &c.Endpoint,
		Region:      &c.Region,
		Credentials: awscreds.NewStaticCredentials(c.AccessKey, c.SecretKey, ""),
	}))
	client := sthree.New(sess)

	// make sure this thing exists
	if _, err := client.CreateBucket(&sthree.CreateBucketInput{
		Bucket: aws.String(c.SpaceName),
	}); err != nil &&
		(!strings.Contains(err.Error(), "already exists") && !strings.Contains(err.Error(), "not empty")) {
		log.Fatalf("Error making bucket %s", err)
	}

	return &Space{
		conf:   c,
		client: client,
	}
}

func (s Space) Create(ctx context.Context, request *pb.CreateRequest, response *pb.CreateResponse) error {
	var err error
	response.Url, err = s.upsert(ctx, request.Object, request.Name, request.Visibility, "space.Create", true)
	return err
}

func (s Space) upsert(ctx context.Context, object []byte, name, visibility, method string, create bool) (string, error) {
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		return "", errors.Unauthorized(method, "Unauthorized")
	}
	if len(name) == 0 {
		return "", errors.BadRequest(method, "Missing name param")
	}
	objectName := fmt.Sprintf("%s/%s", tnt, name)
	if err := s3utils.CheckValidObjectName(objectName); err != nil {
		return "", errors.BadRequest(method, "Invalid name")
	}

	exists := false
	_, err := s.client.HeadObject(&sthree.HeadObjectInput{
		Bucket: aws.String(s.conf.SpaceName),
		Key:    aws.String(objectName),
	})
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if !ok || aerr.Code() != "NotFound" {
			return "", errors.InternalServerError(method, "Error creating object")
		}
	} else {
		exists = true
	}

	if create && exists {
		return "", errors.BadRequest(method, "Object already exists")
	}

	if len(visibility) == 0 {
		visibility = visibilityPrivate
	}

	now := time.Now().Format(time.RFC3339Nano)
	md := meta{
		CreateTime:   now,
		ModifiedTime: now,
		Visibility:   visibility,
	}
	if exists {
		m, err := s.objectMeta(objectName)
		if err != nil {
			log.Errorf("Error reading object meta %s", err)
			return "", errors.BadRequest(method, "Error creating object")
		}
		md.CreateTime = m.CreateTime
	}

	putInput := &sthree.PutObjectInput{
		Body:   bytes.NewReader(object),
		Key:    aws.String(objectName),
		Bucket: aws.String(s.conf.SpaceName),
	}
	if visibility == visibilityPublic {
		putInput.ACL = aws.String(mdACLPublic)
	}

	if _, err := s.client.PutObject(putInput); err != nil {
		log.Errorf("Error creating object %s", err)
		return "", errors.InternalServerError(method, "Error creating object")
	}

	// store the metadata for easy retrieval for listing
	if err := store.Write(store.NewRecord(fmt.Sprintf("%s/%s", prefixByUser, objectName), md)); err != nil {
		log.Errorf("Error writing object to store %s", err)
		return "", errors.InternalServerError(method, "Error creating object")
	}
	retUrl := ""
	if visibility == visibilityPublic {
		retUrl = fmt.Sprintf("%s/%s", s.conf.BaseURL, objectName)
	}

	return retUrl, nil

}

func (s Space) Update(ctx context.Context, request *pb.UpdateRequest, response *pb.UpdateResponse) error {
	var err error
	response.Url, err = s.upsert(ctx, request.Object, request.Name, request.Visibility, "space.Update", false)
	return err
}

func (s Space) Delete(ctx context.Context, request *pb.DeleteRequest, response *pb.DeleteResponse) error {
	method := "space.Delete"
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.Unauthorized(method, "Unauthorized")
	}
	if len(request.Name) == 0 {
		return errors.BadRequest(method, "Missing name param")
	}
	objectName := fmt.Sprintf("%s/%s", tnt, request.Name)
	if _, err := s.client.DeleteObject(&sthree.DeleteObjectInput{
		Bucket: aws.String(s.conf.SpaceName),
		Key:    aws.String(objectName),
	}); err != nil {
		log.Errorf("Error deleting object %s", err)
		return errors.InternalServerError(method, "Error deleting object")
	}
	if err := store.Delete(fmt.Sprintf("%s/%s", prefixByUser, objectName)); err != nil {
		log.Errorf("Error deleting store record %s", err)
		return errors.InternalServerError(method, "Error deleting object")
	}
	return nil
}

func (s Space) List(ctx context.Context, request *pb.ListRequest, response *pb.ListResponse) error {
	method := "space.List"
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.Unauthorized(method, "Unauthorized")
	}
	objectName := fmt.Sprintf("%s/%s", tnt, request.Prefix)
	rsp, err := s.client.ListObjects(&sthree.ListObjectsInput{
		Bucket: aws.String(s.conf.SpaceName),
		Prefix: aws.String(objectName),
	})
	if err != nil {
		log.Errorf("Error listing objects %s", err)
		return errors.InternalServerError(method, "Error listing objects")
	}

	recs, err := store.Read(fmt.Sprintf("%s/%s", prefixByUser, objectName), store.ReadPrefix())
	if err != nil && err != store.ErrNotFound {
		log.Errorf("Error listing objects %s", err)
		return errors.InternalServerError(method, "Error listing objects")
	}
	md := map[string]meta{}
	for _, r := range recs {
		var m meta
		if err := json.Unmarshal(r.Value, &m); err != nil {
			log.Errorf("Error unmarshaling meta %s", err)
			return errors.InternalServerError(method, "Error listing objects")
		}
		md[strings.TrimPrefix(r.Key, prefixByUser+"/")] = m
	}
	response.Objects = []*pb.ListObject{}
	for _, oi := range rsp.Contents {
		m, ok := md[*oi.Key]
		if !ok {
			// hack for now
			m = meta{}
		}
		url := ""
		if m.Visibility == "public" {
			url = fmt.Sprintf("%s/%s", s.conf.BaseURL, *oi.Key)
		}
		response.Objects = append(response.Objects, &pb.ListObject{
			Name:       strings.TrimPrefix(*oi.Key, tnt+"/"),
			Modified:   oi.LastModified.Format(time.RFC3339Nano),
			Url:        url,
			Visibility: m.Visibility,
			Created:    m.CreateTime,
		})
	}
	return nil
}

func (s Space) Head(ctx context.Context, request *pb.HeadRequest, response *pb.HeadResponse) error {
	method := "space.Head"
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.Unauthorized(method, "Unauthorized")
	}
	if len(request.Name) == 0 {
		return errors.BadRequest(method, "Missing name param")
	}
	objectName := fmt.Sprintf("%s/%s", tnt, request.Name)

	goo, err := s.client.HeadObject(&sthree.HeadObjectInput{
		Bucket: aws.String(s.conf.SpaceName),
		Key:    aws.String(objectName),
	})
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok && aerr.Code() == "NotFound" {
			return errors.BadRequest(method, "Object not found")
		}
		log.Errorf("Error s3 %s", err)
		return errors.InternalServerError(method, "Error reading object")
	}

	md, err := s.objectMeta(objectName)
	if err != nil && err != store.ErrNotFound {
		log.Errorf("Error reading object meta %s", err)
		return errors.InternalServerError(method, "Error reading object")
	}
	if md == nil {
		md, err = s.reconstructMeta(ctx, method, objectName, *goo.LastModified)
		if err != nil {
			return err
		}
	}

	url := ""
	if md.Visibility == visibilityPublic {
		url = fmt.Sprintf("%s/%s", s.conf.BaseURL, objectName)
	}
	response.Object = &pb.HeadObject{
		Name:       request.Name,
		Modified:   goo.LastModified.Format(time.RFC3339Nano),
		Created:    md.CreateTime,
		Visibility: md.Visibility,
		Url:        url,
	}

	return nil
}

func (s *Space) objectMeta(objName string) (*meta, error) {
	recs, err := store.Read(fmt.Sprintf("%s/%s", prefixByUser, objName))
	if err != nil {
		return nil, err
	}
	var me meta
	if err := json.Unmarshal(recs[0].Value, &me); err != nil {
		return nil, err
	}
	return &me, nil
}

func (s *Space) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	method := "space.Read"
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.Unauthorized(method, "Unauthorized")
	}

	name := req.Name

	if len(req.Name) == 0 {
		return errors.BadRequest(method, "Missing name param")
	}

	objectName := fmt.Sprintf("%s/%s", tnt, name)

	goo, err := s.client.GetObject(&sthree.GetObjectInput{
		Bucket: aws.String(s.conf.SpaceName),
		Key:    aws.String(objectName),
	})
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok && aerr.Code() == "NoSuchKey" {
			return errors.BadRequest(method, "Object not found")
		}
		log.Errorf("Error s3 %s", err)
		return errors.InternalServerError(method, "Error reading object")
	}

	md, err := s.objectMeta(objectName)
	if err != nil && err != store.ErrNotFound {
		log.Errorf("Error reading meta %s", err)
		return errors.InternalServerError(method, "Error reading object")
	}
	if md == nil {
		md, err = s.reconstructMeta(ctx, method, objectName, *goo.LastModified)
		if err != nil {
			return err
		}
	}

	url := ""
	if md.Visibility == visibilityPublic {
		url = fmt.Sprintf("%s/%s", s.conf.BaseURL, objectName)
	}

	if *goo.ContentLength > maxReadSize {
		return errors.BadRequest(method, "Exceeds max read size: %v bytes", maxReadSize)
	}

	b, err := ioutil.ReadAll(goo.Body)
	if err != nil {
		return errors.InternalServerError(method, "Failed to read data")
	}

	rsp.Object = &pb.Object{
		Name:       req.Name,
		Modified:   goo.LastModified.Format(time.RFC3339Nano),
		Created:    md.CreateTime,
		Visibility: md.Visibility,
		Url:        url,
		Data:       b,
	}

	return nil
}

func (s *Space) reconstructMeta(ctx context.Context, method, objectName string, lastMod time.Time) (*meta, error) {
	aclo, err := s.client.GetObjectAcl(&sthree.GetObjectAclInput{
		Bucket: aws.String(s.conf.SpaceName),
		Key:    aws.String(objectName),
	})
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok && aerr.Code() == "NotFound" {
			return nil, errors.BadRequest(method, "Object not found")
		}
		log.Errorf("Error s3 %s", err)
		return nil, errors.InternalServerError(method, "Error reading object")
	}

	vis := visibilityPrivate

	for _, v := range aclo.Grants {
		if v.Grantee != nil &&
			v.Grantee.URI != nil && *(v.Grantee.URI) == "http://acs.amazonaws.com/groups/global/AllUsers" &&
			v.Permission != nil && *(v.Permission) == "READ" {
			vis = visibilityPublic
			break
		}
	}

	md := &meta{
		Visibility:   vis,
		CreateTime:   lastMod.Format(time.RFC3339Nano),
		ModifiedTime: lastMod.Format(time.RFC3339Nano),
	}
	// store the metadata for easy retrieval for listing
	if err := store.Write(store.NewRecord(fmt.Sprintf("%s/%s", prefixByUser, objectName), md)); err != nil {
		log.Errorf("Error writing object to store %s", err)
		return nil, errors.InternalServerError(method, "Error reading object")
	}
	return md, nil
}

func (s *Space) Download(ctx context.Context, req *api.Request, rsp *api.Response) error {
	method := "space.Download"
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.Unauthorized(method, "Unauthorized")
	}
	var input map[string]string
	if err := json.Unmarshal([]byte(req.Body), &input); err != nil {
		log.Errorf("Error unmarshalling %s", err)
		return errors.BadRequest(method, "Request in unexpected format")
	}
	name := input["name"]
	if len(name) == 0 {
		return errors.BadRequest(method, "Missing name param")
	}

	objectName := fmt.Sprintf("%s/%s", tnt, name)

	_, err := s.client.HeadObject(&sthree.HeadObjectInput{
		Bucket: aws.String(s.conf.SpaceName),
		Key:    aws.String(objectName),
	})
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok && aerr.Code() == "NotFound" {
			return errors.BadRequest(method, "Object not found")
		}
		log.Errorf("Error s3 %s", err)
		return errors.InternalServerError(method, "Error reading object")
	}

	gooreq, _ := s.client.GetObjectRequest(&sthree.GetObjectInput{
		Bucket: aws.String(s.conf.SpaceName),
		Key:    aws.String(objectName),
	})
	urlStr, err := gooreq.Presign(5 * time.Second)
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok && aerr.Code() == "NoSuchKey" {
			return errors.BadRequest(method, "Object not found")
		}
		log.Errorf("Error presigning url %s", err)
		return errors.InternalServerError(method, "Error reading object")
	}

	// replace hostname or url with our base
	split := strings.SplitN(urlStr, s.conf.Endpoint, 2)
	urlStr = s.conf.BaseURL + split[1]

	rsp.Header = map[string]*api.Pair{
		"Location": {
			Key:    "Location",
			Values: []string{urlStr},
		},
	}
	rsp.StatusCode = 302

	resp := map[string]interface{}{
		"url": urlStr,
	}

	b, _ := json.Marshal(resp)
	rsp.Body = string(b)

	return nil
}

func (s Space) Upload(ctx context.Context, request *pb.UploadRequest, response *pb.UploadResponse) error {
	method := "space.Upload"
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.Unauthorized(method, "Unauthorized")
	}
	if len(request.Name) == 0 {
		return errors.BadRequest(method, "Missing name param")
	}
	objectName := fmt.Sprintf("%s/%s", tnt, request.Name)
	if err := s3utils.CheckValidObjectName(objectName); err != nil {
		return errors.BadRequest(method, "Invalid name")
	}

	_, err := s.client.HeadObject(&sthree.HeadObjectInput{
		Bucket: aws.String(s.conf.SpaceName),
		Key:    aws.String(objectName),
	})
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if !ok || aerr.Code() != "NotFound" {
			return errors.InternalServerError(method, "Error creating upload URL")
		}
	} else {
		return errors.BadRequest(method, "Object already exists")
	}

	createTime := aws.String(time.Now().Format(time.RFC3339Nano))

	if len(request.Visibility) == 0 {
		request.Visibility = visibilityPrivate
	}
	putInput := &sthree.PutObjectInput{
		Key:    aws.String(objectName),
		Bucket: aws.String(s.conf.SpaceName),
	}
	if request.Visibility == visibilityPublic {
		putInput.ACL = aws.String(mdACLPublic)
	}

	req, _ := s.client.PutObjectRequest(putInput)
	url, err := req.Presign(5 * time.Minute)
	if err != nil {
		return errors.InternalServerError(method, "Error creating upload URL")
	}
	response.Url = url

	// store the metadata for easy retrieval for listing
	if err := store.Write(store.NewRecord(
		fmt.Sprintf("%s/%s", prefixByUser, objectName),
		meta{
			Visibility:   request.Visibility,
			CreateTime:   *createTime,
			ModifiedTime: time.Now().Format(time.RFC3339Nano),
		})); err != nil {
		log.Errorf("Error writing object to store %s", err)
		return errors.InternalServerError(method, "Error creating upload URL")
	}

	return nil
}

func (s Space) DeleteData(ctx context.Context, request *adminpb.DeleteDataRequest, response *adminpb.DeleteDataResponse) error {
	method := "admin.DeleteData"
	_, err := pauth.VerifyMicroAdmin(ctx, method)
	if err != nil {
		return err
	}

	if len(request.TenantId) < 10 { // deliberate length check so we don't delete all the things
		return errors.BadRequest(method, "Missing tenant ID")
	}

	objectName := request.TenantId
	rsp, err := s.client.ListObjects(&sthree.ListObjectsInput{
		Bucket: aws.String(s.conf.SpaceName),
		Prefix: aws.String(objectName),
	})
	if err != nil {
		log.Errorf("Error listing objects %s", err)
		return errors.InternalServerError(method, "Error listing objects")
	}

	oIDs := []*sthree.ObjectIdentifier{}
	for _, v := range rsp.Contents {
		oIDs = append(oIDs, &sthree.ObjectIdentifier{Key: v.Key})
	}

	if _, err := s.client.DeleteObjects(&sthree.DeleteObjectsInput{
		Bucket: aws.String(s.conf.SpaceName),
		Delete: &sthree.Delete{
			Objects: oIDs,
		},
	}); err != nil {
		return err
	}

	log.Infof("Deleted %d objects from s3 for %s", len(oIDs), request.TenantId)

	keys, err := store.List(store.ListPrefix(fmt.Sprintf("%s/%s/", prefixByUser, request.TenantId)))
	if err != nil {
		log.Errorf("Error listing objects %s", err)
		return errors.InternalServerError(method, "Error listing objects")
	}
	for _, k := range keys {
		if err := store.Delete(k); err != nil {
			return err
		}
	}
	log.Infof("Deleted %d objects from store for %s", len(keys), request.TenantId)

	return nil

}
