package handler

import (
	"context"
	"io/ioutil"
	"os/exec"

	"github.com/micro/micro/v3/service/config"
	log "github.com/micro/micro/v3/service/logger"

	"github.com/micro/micro/v3/service/runtime/source/git"

	function "github.com/micro/services/function/proto"
)

type Function struct {
}

func NewFunction() *Function {
	v, err := config.Get("function.service_account_keyfile")
	if err != nil {
		log.Fatalf("function.service_account_json: %v", err)
	}
	keyfile := v.String("")

	v, err = config.Get("function.service_account")
	if err != nil {
		log.Fatalf("function.service_account: %v", err)
	}
	accName := v.String("")

	// only root
	err = ioutil.WriteFile("/acc.json", []byte(keyfile), 0700)
	if err != nil {
		log.Fatalf("function.service_account: %v", err)
	}

	// https://cloud.google.com/sdk/docs/authorizing#authorizing_with_a_service_account
	exec.Command("gcloud", "auth", "activate-service-account", accName, "--key-file", "/acc.json")
	return &Function{}
}

func (e *Function) Deploy(ctx context.Context, req *function.DeployRequest, rsp *function.DeployResponse) error {
	log.Info("Received Function.Deploy request")
	gitter := git.NewGitter(map[string]string{})
	err := gitter.Checkout(req.Repo, "master")
	if err != nil {
		return err
	}
	return nil
}

func (e *Function) Call(ctx context.Context, req *function.CallRequest, rsp *function.CallResponse) error {
	log.Info("Received Function.Call request")

	return nil
}
