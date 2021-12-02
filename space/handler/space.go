package handler

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/pkg/tenant"
	pb "github.com/micro/services/space/proto"
	"github.com/minio/minio-go/v7/pkg/s3utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awscreds "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	sthree "github.com/aws/aws-sdk-go/service/s3"
)

const (
	mdACL        = "X-Amz-Acl"
	mdACLPublic  = "public-read"
	mdCreated    = "X-Amz-Meta-Micro-Created" // need the x-amz-meta prefix because minio does non-obvious things with prefixes
	mdVisibility = "X-Amz-Meta-Micro-Visibility"

	visibilityPrivate = "private"
	visibilityPublic  = "public"
)

type Space struct {
	conf   conf
	client *sthree.S3
}

type conf struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Endpoint  string `json:"endpoint"`
	SpaceName string `json:"space_name"`
	SSL       bool   `json:"ssl"`
	Region    string `json:"region"`
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
	if create { // check that this doesn't alraedy exist
		_, err := s.client.HeadObject(&sthree.HeadObjectInput{
			Bucket: aws.String(s.conf.SpaceName),
			Key:    aws.String(objectName),
		})
		if err == nil {
			return "", errors.BadRequest(method, "Object already exists")
		}

		aerr, ok := err.(awserr.Error)
		if !ok || aerr.Code() != "NotFound" {
			return "", errors.InternalServerError(method, "Error creating object")
		}
	}

	putInput := &sthree.PutObjectInput{
		Body:   bytes.NewReader(object),
		Key:    aws.String(objectName),
		Bucket: aws.String(s.conf.SpaceName),
		Metadata: map[string]*string{
			mdCreated:    aws.String(fmt.Sprintf("%d", time.Now().Unix())),
			mdVisibility: aws.String(visibility),
		},
	}
	// TODO flesh out options - might want to do content-type for better serving of object
	if visibility == visibilityPublic {
		putInput.ACL = aws.String(mdACLPublic)
	}

	_, err := s.client.PutObject(putInput)

	if err != nil {
		log.Errorf("Error creating object %s", err)
		return "", errors.InternalServerError(method, "Error creating object")
	}

	// TODO fix the url
	return "", nil // i.Location, nil

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
	return nil
}

func (s Space) List(ctx context.Context, request *pb.ListRequest, response *pb.ListResponse) error {
	method := "space.List"
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.Unauthorized(method, "Unauthorized")
	}
	if len(request.Prefix) == 0 {
		return errors.BadRequest(method, "Missing prefix param")
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
	response.Objects = []*pb.ListObject{}
	for _, oi := range rsp.Contents {
		response.Objects = append(response.Objects, &pb.ListObject{
			Name:     strings.TrimPrefix(*oi.Key, tnt+"/"),
			Modified: oi.LastModified.Unix(),
		})
	}
	return nil
}

func (s Space) Read(ctx context.Context, request *pb.ReadRequest, response *pb.ReadResponse) error {
	method := "space.Read"
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.Unauthorized(method, "Unauthorized")
	}
	if len(request.Name) == 0 {
		return errors.BadRequest(method, "Missing name param")
	}
	objectName := fmt.Sprintf("%s/%s", tnt, request.Name)

	// TODO replace with HeadObject?
	goo, err := s.client.GetObject(&sthree.GetObjectInput{
		Bucket: aws.String(s.conf.SpaceName),
		Key:    aws.String(objectName),
	})
	if err != nil {
		// TODO check for not exists
		log.Errorf("Error s3 %s", err)
		aerr, ok := err.(awserr.Error)
		if ok && aerr.Code() == "NoSuchKey" {
			return errors.BadRequest(method, "Object not found")
		}
		return errors.InternalServerError(method, "Error reading object")
	}

	log.Infof("OIII %+v", goo)

	vis := visibilityPrivate
	if md, ok := goo.Metadata[mdVisibility]; ok && len(*md) > 0 {
		vis = *md
	}
	var created int64
	if md, ok := goo.Metadata[mdCreated]; ok && len(*md) > 0 {
		created, err = strconv.ParseInt(*md, 10, 64)
		if err != nil {
			log.Errorf("Error %s", err)
		}
	}

	response.Object = &pb.ReadObject{
		Name:       request.Name,
		Modified:   goo.LastModified.Unix(),
		Created:    created,
		Visibility: vis,
	}

	return nil
}
