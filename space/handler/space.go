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
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/s3utils"
)

const (
	mdACL       = "X-Amz-Acl"
	mdACLPublic = "public-read"
	mdCreated   = "x-amz-meta-micro-created" // need the x-amz-meta prefix because minio does non-obvious things with prefixes

	visibilityPrivate = "private"
	visibilityPublic  = "public"
)

type Space struct {
	client *minio.Client
	conf   conf
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

	cl, err := minio.New(c.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV2(c.AccessKey, c.SecretKey, ""),
		Secure: c.SSL,
		Region: c.Region,
	})
	if err != nil {
		log.Fatalf("Failed to load minio client %s", err)
	}

	bucks, err := cl.ListBuckets(context.Background())
	if err != nil {
		log.Errorf("Error listing buckets %s", err)
	}
	for _, b := range bucks {
		log.Infof("%v", b)
	}
	// make sure this thing exists
	if err := cl.MakeBucket(context.Background(), c.SpaceName, minio.MakeBucketOptions{}); err != nil &&
		(!strings.Contains(err.Error(), "already exists") && !strings.Contains(err.Error(), "not empty")) {
		log.Fatalf("Error making bucket %s", err)
	}

	return &Space{
		client: cl,
		conf:   c,
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
	md := map[string]string{}
	if create { // check that this doesn't alraedy exist
		o, err := s.client.GetObject(ctx, s.conf.SpaceName, objectName, minio.GetObjectOptions{})
		if err == nil {
			_, err := o.Stat()
			if err == nil {
				return "", errors.BadRequest(method, "Object already exists")
			}
		}
		md[mdCreated] = fmt.Sprintf("%d", time.Now().Unix())
	}

	buf := bytes.NewBuffer(object)
	// TODO flesh out options - might want to do content-type for better serving of object
	if visibility == visibilityPublic {
		md[mdACL] = mdACLPublic
	}
	i, err := s.client.PutObject(ctx, s.conf.SpaceName, objectName, buf, int64(buf.Len()), minio.PutObjectOptions{
		UserMetadata: md,
	})

	if err != nil {
		log.Errorf("Error creating object %s", err)
		return "", errors.InternalServerError(method, "Error creating object")
	}

	// TODO fix the url
	return i.Location, nil

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
	if err := s.client.RemoveObject(ctx, s.conf.SpaceName, objectName, minio.RemoveObjectOptions{}); err != nil {
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
	ch := s.client.ListObjects(ctx, s.conf.SpaceName, minio.ListObjectsOptions{
		WithMetadata: true,
		Prefix:       objectName,
		Recursive:    true,
	})
	response.Objects = []*pb.ListObject{}
	for oi := range ch {
		response.Objects = append(response.Objects, &pb.ListObject{
			Name:     strings.TrimPrefix(oi.Key, tnt+"/"),
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
	o, err := s.client.GetObject(ctx, s.conf.SpaceName, objectName, minio.GetObjectOptions{})
	if err != nil {
		log.Errorf("Error reading object %s", err)
		return errors.InternalServerError(method, "Error reading object")
	}
	oi, err := o.Stat()
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return errors.BadRequest(method, "Object does not exist")
		}
		log.Errorf("Error statting object %s", err)
		return errors.InternalServerError(method, "Error reading object")
	}

	log.Infof("OI %+v", oi)
	vis := visibilityPrivate
	if md := oi.Metadata[mdACL]; len(md) > 0 {
		vis = md[0]
	}
	var created int64
	if md := oi.Metadata[strings.Title(mdCreated)]; len(md) > 0 {
		created, err = strconv.ParseInt(md[0], 10, 64)
		if err != nil {
			log.Errorf("Error %s", err)
		}
	}
	response.Object = &pb.ReadObject{
		Name:       strings.TrimPrefix(oi.Key, tnt+"/"),
		Modified:   oi.LastModified.Unix(),
		Created:    created,
		Visibility: vis,
	}

	return nil
}
