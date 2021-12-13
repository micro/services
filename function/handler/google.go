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
	"github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/runtime/source/git"
	"github.com/micro/micro/v3/service/store"
	function "github.com/micro/services/function/proto"
	"github.com/micro/services/pkg/tenant"
	"gopkg.in/yaml.v2"
)

type GoogleFunction struct {
	project string
	// eg. https://us-central1-m3o-apis.cloudfunctions.net/
	address string
	limit   int
}

var (
	GoogleRuntimes = []string{
		"nodejs16",
		"nodejs14",
		"nodejs12",
		"nodejs10",
		"nodejs8",
		"nodejs6",
		"python39",
		"python38",
		"python37",
		"go116",
		"go113",
		"go111",
		"java11",
		"dotnet3",
		"ruby27",
		"ruby26",
		"php74",
	}
)

func NewFunction() *GoogleFunction {
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
	v, err = config.Get("function.limit")
	if err != nil {
		log.Fatalf("function.limit: %v", err)
	}
	limit := v.Int(0)
	if limit > 0 {
		log.Infof("Function limit is %d", limit)
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
	return &GoogleFunction{project: project, address: address, limit: limit}
}

func (e *GoogleFunction) Deploy(ctx context.Context, req *function.DeployRequest, rsp *function.DeployResponse) error {
	log.Info("Received Function.Deploy request")

	if len(req.Name) == 0 {
		return errors.BadRequest("function.deploy", "Missing name")
	}

	if len(req.Repo) == 0 {
		return errors.BadRequest("function.deploy", "Missing repo")
	}
	if len(req.Runtime) == 0 {
		return errors.BadRequest("function.deploy", "invalid runtime")
	}

	var match bool
	for _, r := range GoogleRuntimes {
		if r == req.Runtime {
			match = true
			break
		}
	}

	if !match {
		return errors.BadRequest("function.deploy", "invalid runtime")
	}

	gitter := git.NewGitter(map[string]string{})

	var err error

	for _, branch := range []string{"master", "main"} {
		err = gitter.Checkout(req.Repo, branch)
		if err == nil {
			break
		}
	}

	if err != nil {
		return errors.InternalServerError("function.deploy", err.Error())
	}

	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		tenantId = "micro"
	}

	multitenantPrefix := strings.Replace(tenantId, "/", "-", -1)
	if req.Entrypoint == "" {
		req.Entrypoint = req.Name
	}

	project := req.Project
	if project == "" {
		project = "default"
	}

	key := fmt.Sprintf("function/%s/%s/%s", tenantId, project, req.Name)

	records, err := store.Read(key)
	if err != nil && err != store.ErrNotFound {
		return err
	}

	if len(records) > 0 {
		return errors.BadRequest("function.deploy", "already exists")
	}

	// check for function limit
	if e.limit > 0 {
		// read all the records for the user
		records, err := store.Read("function/"+tenantId+"/", store.ReadPrefix())
		if err != nil {
			return err
		}

		if v := len(records); v >= e.limit {
			return errors.BadRequest("function.deploy", "deployment limit reached")
		}
	}

	// process the env vars to the required format
	var envVars []string

	for k, v := range req.EnvVars {
		envVars = append(envVars, k+"="+v)
	}

	go func() {
		// https://jsoverson.medium.com/how-to-deploy-node-js-functions-to-google-cloud-8bba05e9c10a
		cmd := exec.Command("gcloud", "functions", "deploy",
			multitenantPrefix+"-"+req.Name, "--region", "europe-west1",
			"--allow-unauthenticated", "--entry-point", req.Entrypoint,
			"--trigger-http", "--project", e.project, "--runtime", req.Runtime)

		// if env vars exist then set them
		if len(envVars) > 0 {
			cmd.Args = append(cmd.Args, "--set-env-vars", strings.Join(envVars, ","))
		}

		cmd.Dir = filepath.Join(gitter.RepoDir(), req.Subfolder)
		outp, err := cmd.CombinedOutput()
		if err != nil {
			log.Error(fmt.Errorf(string(outp)))
		}
	}()

	id := fmt.Sprintf("%v-%v-%v", tenantId, project, req.Name)
	rec := store.NewRecord(key, map[string]interface{}{
		"id":         id,
		"project":    project,
		"name":       req.Name,
		"tenantId":   tenantId,
		"repo":       req.Repo,
		"subfolder":  req.Subfolder,
		"entrypoint": req.Entrypoint,
		"runtime":    req.Runtime,
		"env_vars":   envVars,
	})

	// write the record
	return store.Write(rec)
}

func (e *GoogleFunction) Update(ctx context.Context, req *function.UpdateRequest, rsp *function.UpdateResponse) error {
	log.Info("Received Function.Update request")

	if len(req.Name) == 0 {
		return errors.BadRequest("function.update", "Missing name")
	}

	if len(req.Repo) == 0 {
		return errors.BadRequest("function.update", "Missing repo")
	}
	if req.Runtime == "" {
		return fmt.Errorf("missing runtime field, please specify nodejs14, go116 etc")
	}

	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		tenantId = "micro"
	}

	multitenantPrefix := strings.Replace(tenantId, "/", "-", -1)
	if req.Entrypoint == "" {
		req.Entrypoint = req.Name
	}

	project := req.Project
	if project == "" {
		project = "default"
	}

	key := fmt.Sprintf("function/%s/%s/%s", tenantId, project, req.Name)

	records, err := store.Read(key)
	if err != nil {
		return err
	}

	if len(records) == 0 {
		return errors.BadRequest("function.deploy", "function does not exist")
	}

	gitter := git.NewGitter(map[string]string{})

	for _, branch := range []string{"master", "main"} {
		err = gitter.Checkout(req.Repo, branch)
		if err == nil {
			break
		}
	}

	if err != nil {
		return errors.InternalServerError("function.update", err.Error())
	}

	// process the env vars to the required format
	var envVars []string

	for k, v := range req.EnvVars {
		envVars = append(envVars, k+"="+v)
	}

	go func() {
		// https://jsoverson.medium.com/how-to-deploy-node-js-functions-to-google-cloud-8bba05e9c10a
		cmd := exec.Command("gcloud", "functions", "deploy",
			multitenantPrefix+"-"+req.Name, "--region", "europe-west1",
			"--allow-unauthenticated", "--entry-point", req.Entrypoint,
			"--trigger-http", "--project", e.project, "--runtime", req.Runtime)

		// if env vars exist then set them
		if len(envVars) > 0 {
			cmd.Args = append(cmd.Args, "--set-env-vars", strings.Join(envVars, ","))
		}

		cmd.Dir = filepath.Join(gitter.RepoDir(), req.Subfolder)
		outp, err := cmd.CombinedOutput()
		if err != nil {
			log.Error(fmt.Errorf(string(outp)))
		}
	}()

	id := fmt.Sprintf("%v-%v-%v", tenantId, project, req.Name)
	rec := store.NewRecord(key, map[string]interface{}{
		"id":         id,
		"project":    project,
		"name":       req.Name,
		"tenantId":   tenantId,
		"repo":       req.Repo,
		"subfolder":  req.Subfolder,
		"entrypoint": req.Entrypoint,
		"runtime":    req.Runtime,
		"env_vars":   envVars,
	})

	// write the record
	return store.Write(rec)
}

func (e *GoogleFunction) Call(ctx context.Context, req *function.CallRequest, rsp *function.CallResponse) error {
	log.Info("Received Function.Call request")

	if len(req.Name) == 0 {
		return errors.BadRequest("function.call", "Missing function name")
	}

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

func (e *GoogleFunction) Delete(ctx context.Context, req *function.DeleteRequest, rsp *function.DeleteResponse) error {
	log.Info("Received Function.Delete request")

	if len(req.Name) == 0 {
		return errors.BadRequest("function.delete", "Missing function name")
	}

	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		tenantId = "micro"
	}
	multitenantPrefix := strings.Replace(tenantId, "/", "-", -1)

	project := req.Project
	if project == "" {
		project = "default"
	}

	cmd := exec.Command("gcloud", "functions", "delete", "--project", e.project, "--region", "europe-west1", multitenantPrefix+"-"+req.Name)
	outp, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(fmt.Errorf(string(outp)))
		return err
	}

	key := fmt.Sprintf("function/%v/%v/%v", tenantId, project, req.Name)

	return store.Delete(key)
}

func (e *GoogleFunction) List(ctx context.Context, req *function.ListRequest, rsp *function.ListResponse) error {
	log.Info("Received Function.List request")

	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		tenantId = "micro"
	}

	key := "function/" + tenantId + "/"

	project := req.Project
	if len(project) > 0 {
		key = key + "/" + project + "/"
	}

	records, err := store.Read(key, store.ReadPrefix())
	if err != nil {
		return err
	}

	multitenantPrefix := strings.Replace(tenantId, "/", "-", -1)
	cmd := exec.Command("gcloud", "functions", "list", "--project", e.project, "--filter", "name~"+multitenantPrefix+"*")
	outp, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(fmt.Errorf(string(outp)))
	}

	lines := strings.Split(string(outp), "\n")
	statuses := map[string]string{}

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		statuses[fields[0]] = fields[1]
	}

	rsp.Functions = []*function.Func{}

	for _, record := range records {
		f := new(function.Func)
		if err := record.Decode(f); err != nil {
			return err
		}
		f.Status = statuses[multitenantPrefix+"-"+f.Name]
		rsp.Functions = append(rsp.Functions, f)
	}

	return nil
}

func (e *GoogleFunction) Describe(ctx context.Context, req *function.DescribeRequest, rsp *function.DescribeResponse) error {
	if len(req.Name) == 0 {
		return errors.BadRequest("function.describe", "Missing function name")
	}

	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		tenantId = "micro"
	}

	project := req.Project
	if project == "" {
		project = "default"
	}

	multitenantPrefix := strings.Replace(tenantId, "/", "-", -1)
	key := fmt.Sprintf("function/%v/%v/%v", tenantId, project, req.Name)

	records, err := store.Read(key)
	if err != nil {
		return err
	}

	cmd := exec.Command("gcloud", "functions", "describe", "--region", "europe-west1", "--project", e.project, multitenantPrefix+"-"+req.Name)
	outp, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(fmt.Errorf(string(outp)))
		return fmt.Errorf("function does not exist")
	}

	log.Info(string(outp))
	m := map[string]interface{}{}
	err = yaml.Unmarshal(outp, m)
	if err != nil {
		return err
	}

	if len(records) > 0 {
		f := &function.Func{}
		if err := records[0].Decode(f); err != nil {
			return err
		}
		rsp.Function = f
	} else {
		rsp.Function = &function.Func{
			Name:    req.Name,
			Project: req.Project,
		}
	}

	// set describe info
	rsp.Function.Status = m["status"].(string)
	rsp.Timeout = m["timeout"].(string)
	rsp.UpdatedAt = m["updateTime"].(string)

	return nil
}
