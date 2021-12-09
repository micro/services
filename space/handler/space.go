package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/micro/micro/v3/proto/api"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
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
	mdACL        = "X-Amz-Acl"
	mdACLPublic  = "public-read"
	mdCreated    = "Micro-Created"
	mdVisibility = "Micro-Visibility"

	visibilityPrivate = "private"
	visibilityPublic  = "public"

	prefixByUser = "byUser"
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
	hoo, err := s.client.HeadObject(&sthree.HeadObjectInput{
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

	createTime := aws.String(time.Now().Format(time.RFC3339Nano))
	if exists {
		createTime = hoo.Metadata[mdCreated]
	}

	if len(visibility) == 0 {
		visibility = visibilityPrivate
	}
	putInput := &sthree.PutObjectInput{
		Body:   bytes.NewReader(object),
		Key:    aws.String(objectName),
		Bucket: aws.String(s.conf.SpaceName),
		Metadata: map[string]*string{
			mdVisibility: aws.String(visibility),
			mdCreated:    createTime,
		},
	}
	// TODO flesh out options - might want to do content-type for better serving of object
	if visibility == visibilityPublic {
		putInput.ACL = aws.String(mdACLPublic)
	}

	if _, err := s.client.PutObject(putInput); err != nil {
		log.Errorf("Error creating object %s", err)
		return "", errors.InternalServerError(method, "Error creating object")
	}

	// store the metadata for easy retrieval for listing
	if err := store.Write(store.NewRecord(
		fmt.Sprintf("%s/%s", prefixByUser, objectName),
		meta{Visibility: visibility, CreateTime: *createTime, ModifiedTime: time.Now().Format(time.RFC3339Nano)})); err != nil {
		log.Errorf("Error writing object to store %s", err)
		return "", errors.InternalServerError(method, "Error creating object")
	}
	retUrl := ""
	if visibility == "public" {
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
	if err != nil {
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

	vis := visibilityPrivate
	if md, ok := goo.Metadata[mdVisibility]; ok && len(*md) > 0 {
		vis = *md
	}
	var created string
	if md, ok := goo.Metadata[mdCreated]; ok && len(*md) > 0 {
		t, err := time.Parse(time.RFC3339Nano, *md)
		if err != nil {
			// try as unix ts
			createdI, err := strconv.ParseInt(*md, 10, 64)
			if err != nil {
				log.Errorf("Error %s", err)
			} else {
				t = time.Unix(createdI, 0)
			}
		}
		created = t.Format(time.RFC3339Nano)
	}

	url := ""
	if vis == "public" {
		url = fmt.Sprintf("%s/%s", s.conf.BaseURL, objectName)
	}
	response.Object = &pb.HeadObject{
		Name:       request.Name,
		Modified:   goo.LastModified.Format(time.RFC3339Nano),
		Created:    created,
		Visibility: vis,
		Url:        url,
	}

	return nil
}

func (s *Space) Read(ctx context.Context, req *api.Request, rsp *api.Response) error {
	method := "space.Read"
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
	rsp.Header = map[string]*api.Pair{
		"Location": {
			Key:    "Location",
			Values: []string{urlStr},
		},
	}
	rsp.StatusCode = 302

	return nil
}
