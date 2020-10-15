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

	return nil
}
