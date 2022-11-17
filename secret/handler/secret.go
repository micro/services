package handler

import (
	"context"
	"encoding/base64"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/micro/micro/v3/service/config"
	merrors "github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	"github.com/micro/services/pkg/tenant"
	pb "github.com/micro/services/secret/proto"
)

const (
	defaultNamespace = "micro"
	pathSplitter     = "."
)

var (
	// we now support json only
	mtx sync.RWMutex
)

type Secret struct {
	Key []byte
}

type Value struct {
	Key     string
	Data    string
	Created time.Time
	Updated time.Time
}

func New() *Secret {
	var dec []byte
	var err error

	val, err := config.Get("secret.key")
	if err != nil {
		logger.Fatal("Missing key")
	}

	key := val.String("")

	if len(key) == 0 {
		logger.Warn("No encryption key provided")
	} else {
		dec, err = base64.StdEncoding.DecodeString(key)
		if err != nil {
			logger.Warnf("Error decoding key: %v", err)
		}
	}

	return &Secret{
		Key: dec,
	}
}

func (s *Secret) Get(ctx context.Context, req *pb.GetRequest, rsp *pb.GetResponse) error {
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		tnt = "micro"
	}

	key := path.Join(tnt, req.Key)

	rec, err := store.Read(key)
	if err == store.ErrNotFound {
		return merrors.NotFound("secret.get", "Not found")
	} else if err != nil {
		return merrors.BadRequest("secret.get", err.Error())
	}

	// extract value
	v := new(Value)
	rec[0].Decode(v)

	//decode the val
	dec, err := base64.StdEncoding.DecodeString(v.Data)
	if err != nil {
		return err
	}

	// decrypt it
	decrypted, err := decrypt(string(dec), []byte(s.Key))
	if err != nil {
		return err
	}

	var val []byte

	// check path
	if len(req.Path) > 0 {
		path := strings.Replace(req.Path, "/", ".", -1)
		vals := config.NewJSONValues([]byte(decrypted))
		val = vals.Get(path).Bytes()
	} else {
		// take whole value
		val = []byte(decrypted)
	}

	// set response values
	rsp.Key = req.Key
	rsp.Path = req.Path
	rsp.Value = string(val)
	rsp.Created = v.Created.Format(time.RFC3339Nano)
	rsp.Updated = v.Updated.Format(time.RFC3339Nano)

	return nil
}

func (s *Secret) Set(ctx context.Context, req *pb.SetRequest, rsp *pb.SetResponse) error {
	if len(req.Key) == 0 {
		return merrors.BadRequest("secret.set", "missing key")
	}
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		tnt = "micro"
	}

	// key to store under
	key := path.Join(tnt, req.Key)

	v := new(Value)
	data := []byte{}

	// get existing value
	rec, err := store.Read(key)
	if err == nil && len(rec) > 0 {
		rec[0].Decode(v)

		//decode the val
		dec, err := base64.StdEncoding.DecodeString(v.Data)
		if err != nil {
			return err
		}

		// decrypt it
		decrypted, err := decrypt(string(dec), []byte(s.Key))
		if err != nil {
			return err
		}

		// set the decrypted value
		data = []byte(decrypted)
	} else {
		data = []byte(`{}`)
		v.Key = req.Key
		v.Created = time.Now()
	}

	// update timestamp
	v.Updated = time.Now()

	// there is a path to deal with
	if len(req.Path) > 0 {
		path := strings.Replace(req.Path, "/", ".", -1)
		vals := config.NewJSONValues(data)
		vals.Set(path, req.Value)
		data = vals.Bytes()
	} else {
		data = []byte(req.Value)
	}

	// encrypt the data
	encrypted, err := encrypt(string(data), s.Key)
	if err != nil {
		return merrors.InternalServerError("secret.set", "Failed to encrypt: %v", err)
	}

	// base64 encode the value
	v.Data = base64.StdEncoding.EncodeToString([]byte(encrypted))

	return store.Write(store.NewRecord(key, v))
}

func (s *Secret) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	if len(req.Key) == 0 {
		return merrors.BadRequest("secret.delete", "missing key")
	}

	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		tnt = "micro"
	}

	// delete key, no check
	key := path.Join(tnt, req.Key)

	// no path, delete whole key
	if len(req.Path) == 0 {
		return store.Delete(key)
	}

	// replace / with .
	path := strings.Replace(req.Path, "/", ".", -1)
	// get existing value
	rec, err := store.Read(key)
	if err != nil {
		return err
	}

	if len(rec) == 0 {
		return nil
	}

	// extract value
	v := new(Value)
	rec[0].Decode(v)

	//decode the val
	dec, err := base64.StdEncoding.DecodeString(v.Data)
	if err != nil {
		return err
	}

	// decrypt it
	decrypted, err := decrypt(string(dec), []byte(s.Key))
	if err != nil {
		return err
	}

	// delete the path
	vals := config.NewJSONValues([]byte(decrypted))
	vals.Delete(path)

	// get the data
	data := vals.Bytes()

	// encrypt the data
	encrypted, err := encrypt(string(data), s.Key)
	if err != nil {
		return merrors.InternalServerError("secret.set", "Failed to encrypt: %v", err)
	}

	// base64 encode the value
	v.Data = base64.StdEncoding.EncodeToString([]byte(encrypted))
	// updated
	v.Updated = time.Now()

	// put it back
	return store.Write(store.NewRecord(key, v))
}

func (s *Secret) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		tnt = "micro"
	}

	key := tnt + "/"

	// list keys by tenant prefix
	keys, err := store.List(store.ListPrefix(key))
	if err != nil {
		return merrors.BadRequest("secret.list", err.Error())
	}

	// return the list of keys with tenant stripped
	for _, val := range keys {
		rsp.Keys = append(rsp.Keys, strings.TrimPrefix(val, key))
	}

	return nil
}
