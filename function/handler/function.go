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
	"gopkg.in/yaml.v2"

	"github.com/micro/micro/v3/service/runtime/source/git"

	_struct "github.com/golang/protobuf/ptypes/struct"
	db "github.com/micro/services/db/proto"
	function "github.com/micro/services/function/proto"
	"github.com/micro/services/pkg/tenant"
)

type Function struct {
	project string
	// eg. https://us-central1-m3o-apis.cloudfunctions.net/
	address string
	db      db.DbService
}

type Func struct {
	Name    string `json:"name"`
	Tenant  string `json:"tenant"`
	Project string `json:"project"`
}

func NewFunction(db db.DbService) *Function {
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
	return &Function{project: project, address: address, db: db}
}

func (e *Function) Deploy(ctx context.Context, req *function.DeployRequest, rsp *function.DeployResponse) error {
	log.Info("Received Function.Deploy request")

	if len(req.Name) == 0 {
		return errors.BadRequest("function.deploy", "Missing name")
	}

	if len(req.Repo) == 0 {
		return errors.BadRequest("function.deploy", "Missing repo")
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

	readRsp, err := e.db.Read(ctx, &db.ReadRequest{
		Table: "functions",
		Query: fmt.Sprintf("tenantId == '%v' and project == '%v' and name == '%v'", tenantId, project, req.Name),
	})
	if err != nil {
		return err
	}
	if req.Runtime == "" {
		return fmt.Errorf("missing runtime field, please specify nodejs14, go116 etc")
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

	s := &_struct.Struct{}
	id := fmt.Sprintf("%v-%v-%v", tenantId, project, req.Name)
	jso, _ := json.Marshal(map[string]interface{}{
		"id":         id,
		"project":    project,
		"name":       req.Name,
		"tenantId":   tenantId,
		"repo":       req.Repo,
		"subfolder":  req.Subfolder,
		"entrypoint": req.Entrypoint,
		"runtime":    req.Runtime,
	})
	err = s.UnmarshalJSON(jso)
	if err != nil {
		return err
	}
	if len(readRsp.Records) > 0 {
		_, err = e.db.Update(ctx, &db.UpdateRequest{
			Table:  "functions",
			Record: s,
			Id:     id,
		})
		if err != nil {
			log.Error(err)
		}
		return err
	}
	_, err = e.db.Create(ctx, &db.CreateRequest{
		Table:  "functions",
		Record: s,
	})
	if err != nil {
		log.Error(err)
	}
	return err
}

func (e *Function) Call(ctx context.Context, req *function.CallRequest, rsp *function.CallResponse) error {
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

func (e *Function) Delete(ctx context.Context, req *function.DeleteRequest, rsp *function.DeleteResponse) error {
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

	id := fmt.Sprintf("%v-%v-%v", tenantId, project, req.Name)
	_, err = e.db.Delete(ctx, &db.DeleteRequest{
		Table: "functions",
		Id:    id,
	})
	return err
}

func (e *Function) List(ctx context.Context, req *function.ListRequest, rsp *function.ListResponse) error {
	log.Info("Received Function.List request")

	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		tenantId = "micro"
	}
	project := req.Project

	q := fmt.Sprintf(`tenantId == "%v"`, tenantId)
	if project != "" {
		q += fmt.Sprintf(` and project == "%v"`, project)
	}
	log.Infof("Making query %v", q)
	readRsp, err := e.db.Read(ctx, &db.ReadRequest{
		Table: "functions",
		Query: q,
	})
	if err != nil {
		return err
	}
	log.Info(readRsp.Records)

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
	for _, record := range readRsp.Records {
		m := record.AsMap()
		bs, _ := json.Marshal(m)
		f := &function.Func{}
		err = json.Unmarshal(bs, f)
		if err != nil {
			return err
		}
		f.Status = statuses[multitenantPrefix+"-"+f.Name]
		rsp.Functions = append(rsp.Functions, f)
	}
	return nil
}

func (e *Function) Describe(ctx context.Context, req *function.DescribeRequest, rsp *function.DescribeResponse) error {
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
	id := fmt.Sprintf("%v-%v-%v", tenantId, project, req.Name)

	readRsp, err := e.db.Read(ctx, &db.ReadRequest{
		Table: "functions",
		Id:    id,
	})
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

	if len(readRsp.Records) > 0 {
		m := readRsp.Records[0].AsMap()
		bs, _ := json.Marshal(m)
		f := &function.Func{}
		err = json.Unmarshal(bs, f)
		if err != nil {
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
