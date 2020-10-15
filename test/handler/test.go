package handler

import (
	"context"

	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	test "github.com/micro/services/test/proto"
)

type Test struct{}

// Call is a single request handler called via client.Call or the generated client code
func (t *Test) Call(ctx context.Context, req *test.Request, rsp *test.Response) error {
	log.Info("Received Test.Call request")
	rsp.Status = "OK"
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (t *Test) Stream(ctx context.Context, req *test.Request, stream test.Test_StreamStream) error {
	log.Infof("Received Test.Stream request with id: %v", req.Id)

	for i := 0; i < 10; i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&test.Response{
			Status: "OK",
		}); err != nil {
			return err
		}
	}

	return nil
}

// Config is an endpoint that tests the config
func (t *Test) Config(ctx context.Context, req *test.Request, rsp *test.Response) error {
	type keyConfig struct {
		Subkey  string `json:"subkey"`
		Subkey1 int    `json:"subkey1"`
		Subkey2 string `json:"subkey2"`
	}

	type conf struct {
		Key keyConfig `json:"key"`
	}

	err := config.Set("key", map[string]interface{}{
		"Subkey3": "Merge",
	})
	if err != nil {
		return errors.InternalServerError("test.config", err.Error())
	}

	val, _ := config.Get("key.subkey3")
	if val.String("") != "Merge" {
		return errors.InternalServerError("test.config", "ERROR: key.subkey3 should be 'Merge' but it is:"+val.String(""))
	}

	val, err = config.Get("key.subkey")
	fmt.Println("Value of key.subkey: ", val.String(""), err)

	val, _ = config.Get("key", config.Secret(true))
	c := conf{}
	err = val.Scan(&c.Key)
	fmt.Println("Value of key.subkey1: ", c.Key.Subkey1, err)
	fmt.Println("Value of key.subkey2: ", c.Key.Subkey2)

	val, _ = config.Get("key.subkey3")
	fmt.Println("Value of key.subkey3: ", val.String(""))

	// Test defaults
	val, _ = config.Get("key.subkey_does_not_exist")
	fmt.Println("Default", val.String("Hello"))

	// set status
	rsp.Status = "OK"

	return nil
}

func (t *Test) Store(ctx context.Context, req *test.Request, rsp *test.Response) error {
	key := "Test"
	val := []byte(`foobar`)

	if err := store.Write(&store.Record{Key: key, Value: []byte(val)}); err != nil {
		log.Errorf("Error writing %s", err)
		return errors.InternalServerError("test.store", fmt.Errorf("Error writing record %s with expiry %s", key, err))
	}

	recs, err := store.List()
	if err != nil {
		return fmt.Errorf("Error listing from store %s", err)
	}

	log.Infof("Recs %+v", recs)
	if len(recs) != 1 {
		return fmt.Errorf("Error listing records, expected 3, received %d", len(recs))
	}

	rsp.Status = "OK"
	return nil
}

func (t *Test) Events(ctx context.Context, req *test.Request, rsp *test.Response) error {}

func (t *Test) Broker(ctx context.Context, req *test.Request, rsp *test.Response) error {}

func (t *Test) BlobStore(ctx context.Context, req *test.Request, rsp *test.Response) error {}

func (t *Test) Logger(ctx context.Context, req *test.Request, rsp *test.Response) error {
	log.Infof("Testing logger: %v", req.Id)
	rsp.Status = "OK"
	return nil
}
