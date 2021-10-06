package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/micro/micro/v3/service/config"
	log "github.com/micro/micro/v3/service/logger"

	"github.com/micro/micro/v3/service/runtime/source/git"

	function "github.com/micro/services/function/proto"
	"github.com/micro/services/pkg/tenant"
)

type Function struct {
	project string
	// eg. https://us-central1-m3o-apis.cloudfunctions.net/
	address string
}

func NewFunction() *Function {
	v, err := config.Get("function.service_account_json")
	if err != nil {
		log.Fatalf("function.service_account_json: %v", err)
	}
	keyfile := v.Bytes()
	if len(keyfile) == 0 {
		log.Fatalf("empty keyfile")
	}

	v, err = config.Get("function.address")
	if err != nil {
		log.Fatalf("function.address: %v", err)
	}
	address := v.String("")
	if len(address) == 0 {
		log.Fatalf("empty address")
	}

	v, err = config.Get("function.project")
	if err != nil {
		log.Fatalf("function.project: %v", err)
	}
	project := v.String("")
	if len(project) == 0 {
		log.Fatalf("empty project")
	}

	v, err = config.Get("function.service_account")
	if err != nil {
		log.Fatalf("function.service_account: %v", err)
	}
	accName := v.String("")

	m := map[string]interface{}{}
	err = json.Unmarshal(keyfile, &m)
	if err != nil {
		log.Fatalf("invalid json: %v", err)
	}

	// only root
	err = ioutil.WriteFile("/acc.json", keyfile, 0700)
	if err != nil {
		log.Fatalf("function.service_account: %v", err)
	}

	// DO THIS STEP
	// https://cloud.google.com/functions/docs/reference/iam/roles#additional-configuration

	// https://cloud.google.com/sdk/docs/authorizing#authorizing_with_a_service_account
	outp, err := exec.Command("gcloud", "auth", "activate-service-account", accName, "--key-file", "/acc.json").CombinedOutput()
	if err != nil {
		log.Fatalf(string(outp))
	}

	outp, err = exec.Command("gcloud", "auth", "list").CombinedOutput()
	if err != nil {
		log.Fatalf(string(outp))
	}
	log.Info(string(outp))
	return &Function{project: project, address: address}
}

func (e *Function) Deploy(ctx context.Context, req *function.DeployRequest, rsp *function.DeployResponse) error {
	log.Info("Received Function.Deploy request")
	gitter := git.NewGitter(map[string]string{})
	err := gitter.Checkout(req.Repo, "master")
	if err != nil {
		return err
	}

	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		tenantId = "micro"
	}
	multitenantPrefix := strings.Replace(tenantId, "/", "-", -1)
	if req.Entrypoint == "" {
		req.Entrypoint = req.Name
	}

	// https://jsoverson.medium.com/how-to-deploy-node-js-functions-to-google-cloud-8bba05e9c10a
	cmd := exec.Command("gcloud", "functions", "deploy", multitenantPrefix+"-"+req.Name, "--region", "europe-west1", "--allow-unauthenticated", "--entry-point", req.Entrypoint, "--trigger-http", "--project", e.project, "--runtime", "nodejs14")
	cmd.Dir = filepath.Join(gitter.RepoDir(), req.Subfolder)
	outp, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(string(outp))
	}
	log.Info(string(outp))
	return nil
}

func (e *Function) Call(ctx context.Context, req *function.CallRequest, rsp *function.CallResponse) error {
	log.Info("Received Function.Call request")

	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		tenantId = "micro"
	}
	multitenantPrefix := strings.Replace(tenantId, "/", "-", -1)

	url := e.address + multitenantPrefix + "-" + req.Name
	fmt.Println("URL:>", url)

	js, _ := json.Marshal(req.Request)
	if req.Request == nil || len(req.Request.Fields) == 0 {
		js = []byte("{}")
	}
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(js))
	if err != nil {
		return err
	}
	r.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		log.Errorf("error making request %v", err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("error reading body %v", string(body))
		return err
	}

	err = json.Unmarshal(body, &rsp.Response)
	if err != nil {
		log.Errorf("error unmarshaling %v", string(body))
		return err
	}

	return nil
}
