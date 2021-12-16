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
	"time"

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
	// max functions deployed
	limit int
	// function identity
	identity string
	// custom domain
	domain string
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

	// hardcoded list of supported regions
	GoogleRegions = []string{"europe-west1", "us-east1", "us-west1"}
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
	v, err = config.Get("function.service_identity")
	if err != nil {
		log.Fatalf("function.service_identity: %v", err)
	}
	identity := v.String("")
	v, err = config.Get("function.domain")
	if err != nil {
		log.Fatalf("function.domain: %v", err)
	}
	domain := v.String("")

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

	return &GoogleFunction{
		project:  project,
		address:  address,
		limit:    limit,
		identity: identity,
		domain:   domain,
	}
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

	var supportedRegion bool

	if len(req.Region) == 0 {
		// set to default region
		req.Region = GoogleRegions[0]
		supportedRegion = true
	}

	// check if its in the supported regions
	for _, reg := range GoogleRegions {
		if req.Region == reg {
			supportedRegion = true
			break
		}
	}

	// unsupported region requested
	if !supportedRegion {
		return errors.BadRequest("function.deploy", "Unsupported region")
	}

	if len(req.Branch) == 0 {
		req.Branch = "master"
	}

	gitter := git.NewGitter(map[string]string{})

	var err error

	err = gitter.Checkout(req.Repo, req.Branch)
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

	id := fmt.Sprintf("%v-%v", multitenantPrefix, req.Name)
	fn := &function.Func{
		Id: id,
		Name: req.Name,
		Project: project,
		Repo: req.Repo,
		Subfolder: req.Subfolder,
		Entrypoint: req.Entrypoint,
		Runtime: req.Runtime,
		EnvVars: req.EnvVars,
		Region: req.Region,
		Branch: req.Branch,
		Created: time.Now().Format(time.RFC3339Nano),
		Status: "Deploying",
	}

	rec := store.NewRecord(key, fn)
	store.Write(rec)

	rsp.Function = fn

	go func(fn *function.Func) {
		// https://jsoverson.medium.com/how-to-deploy-node-js-functions-to-google-cloud-8bba05e9c10a
		cmd := exec.Command("gcloud", "functions", "deploy",
			multitenantPrefix+"-"+fn.Name, "--region", fn.Region, "--service-account", e.identity,
			"--allow-unauthenticated", "--entry-point", fn.Entrypoint,
			"--trigger-http", "--project", e.project, "--runtime", fn.Runtime)

		// if env vars exist then set them
		if len(envVars) > 0 {
			cmd.Args = append(cmd.Args, "--set-env-vars", strings.Join(envVars, ","))
		}

		cmd.Dir = filepath.Join(gitter.RepoDir(), fn.Subfolder)
		outp, err := cmd.CombinedOutput()
		if err != nil {
			log.Error(fmt.Errorf(string(outp)))
			fn.Status = "DeploymentError"
			store.Write(store.NewRecord(key, fn))
			return
		}

		cmd = exec.Command("gcloud", "functions", "describe", "--format", "json",
			"--region", fn.Region, "--project", e.project, multitenantPrefix+"-"+fn.Name)

		outp, err = cmd.CombinedOutput()
		if err != nil {
			log.Error(fmt.Errorf(string(outp)))
			return
		}

		log.Info(string(outp))
		m := map[string]interface{}{}
		if err := json.Unmarshal(outp, m); err != nil {
			return
		}

		// write back the url
		trigger := m["httpsTrigger"].(map[string]interface{})
		fn.Url = trigger["url"].(string)
		fn.Updated = time.Now().Format(time.RFC3339Nano)
		store.Write(store.NewRecord(key, fn))
	}(fn)

	return nil
}

func (e *GoogleFunction) Update(ctx context.Context, req *function.UpdateRequest, rsp *function.UpdateResponse) error {
	log.Info("Received Function.Update request")

	if len(req.Name) == 0 {
		return errors.BadRequest("function.update", "Missing name")
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

	key := fmt.Sprintf("function/%s/%s/%s", tenantId, project, req.Name)

	records, err := store.Read(key)
	if err != nil {
		return err
	}

	if len(records) == 0 {
		return errors.BadRequest("function.deploy", "function does not exist")
	}

	fn := new(function.Func)
	if err := records[0].Decode(fn); err != nil {
		return err
	}

	if len(fn.Region) == 0 {
		fn.Region = GoogleRegions[0]
	}

	if len(fn.Branch) == 0 {
		fn.Branch = "master"
	}

	gitter := git.NewGitter(map[string]string{})
	if err := gitter.Checkout(fn.Repo, fn.Branch); err != nil {
		return errors.InternalServerError("function.update", err.Error())
	}

	// process the env vars to the required format
	var envVars []string

	for k, v := range fn.EnvVars {
		envVars = append(envVars, k+"="+v)
	}

	go func() {
		// https://jsoverson.medium.com/how-to-deploy-node-js-functions-to-google-cloud-8bba05e9c10a
		cmd := exec.Command("gcloud", "functions", "deploy",
			multitenantPrefix+"-"+fn.Name, "--region", fn.Region, "--service-account", e.identity,
			"--allow-unauthenticated", "--entry-point", fn.Entrypoint,
			"--trigger-http", "--project", e.project, "--runtime", fn.Runtime)

		// if env vars exist then set them
		if len(envVars) > 0 {
			cmd.Args = append(cmd.Args, "--set-env-vars", strings.Join(envVars, ","))
		}

		cmd.Dir = filepath.Join(gitter.RepoDir(), fn.Subfolder)
		outp, err := cmd.CombinedOutput()
		if err != nil {
			log.Error(fmt.Errorf(string(outp)))
		}

		fn.Updated = time.Now().Format(time.RFC3339Nano)
		store.Write(store.NewRecord(key, fn))
	}()

	// TODO: allow updating of branch and related?
	return nil
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

	project := req.Project
	if project == "" {
		project = "default"
	}

	multitenantPrefix := strings.Replace(tenantId, "/", "-", -1)
	key := fmt.Sprintf("function/%v/%v/%v", tenantId, project, req.Name)

	records, err := store.Read(key)
	if err != nil && err == store.ErrNotFound {
		return nil
	}

	if len(records) == 0 {
		return nil
	}

	fn := new(function.Func)
	if err := records[0].Decode(fn); err != nil {
		return err
	}

	// async delete
	go func() {
		cmd := exec.Command("gcloud", "functions", "delete", "--project", e.project, "--region", fn.Region, multitenantPrefix+"-"+req.Name)
		outp, err := cmd.CombinedOutput()
		if err != nil && !strings.Contains(string(outp), "does not exist") {
			log.Error(fmt.Errorf(string(outp)))
			return
		}

		store.Delete(key)
	}()

	return nil
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
		status := strings.Replace(fields[1], "_", " ", -1)
		status = strings.Title(strings.ToLower(status))
		statuses[fields[0]] = status
	}

	rsp.Functions = []*function.Func{}

	for _, record := range records {
		fn := new(function.Func)
		if err := record.Decode(fn); err != nil {
			return err
		}
		if len(fn.Region) == 0 {
			fn.Region = GoogleRegions[0]
		}

		if len(fn.Branch) == 0 {
			fn.Branch = "master"
		}
		// set the custom domain
		if len(e.domain) > 0 {
			fn.Url = fmt.Sprintf("https://%s.%s", fn.Id, e.domain)
		}

		fn.Status = statuses[multitenantPrefix+"-"+fn.Name]
		rsp.Functions = append(rsp.Functions, fn)
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

	if len(records) == 0 {
		return errors.NotFound("function.describe", "function does not exist")
	}

	fn := new(function.Func)

	if err := records[0].Decode(fn); err != nil {
		return err
	}

	if len(fn.Region) == 0 {
		fn.Region = GoogleRegions[0]
	}

	if len(fn.Branch) == 0 {
		fn.Branch = "master"
	}

	// set the custom domain
	if len(e.domain) > 0 {
		fn.Url = fmt.Sprintf("https://%s.%s", fn.Id, e.domain)
	}

	// set the response function
	rsp.Function = fn

	// get the current status
	cmd := exec.Command("gcloud", "functions", "describe", "--region", fn.Region, "--project", e.project, multitenantPrefix+"-"+req.Name)
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

	// set describe info
	status := m["status"].(string)
	status = strings.Replace(status, "_", " ", -1)
	status = strings.Title(strings.ToLower(status))
	rsp.Function.Status = status

	return nil
}

func (g *GoogleFunction) Proxy(ctx context.Context, req *function.ProxyRequest, rsp *function.ProxyResponse) error {
	if len(req.Id) == 0 {
		return errors.BadRequest("function.proxy", "missing id")
	}

	if !IDFormat.MatchString(req.Id) {
		return errors.BadRequest("function.proxy", "invalid id")
	}

	key := FunctionKey + req.Id

	recs, err := store.Read(key, store.ReadLimit(1))
	if err != nil {
		return err
	}

	if len(recs) == 0 {
		return errors.BadRequest("function.proxy", "function does not exist")
	}

	fn := new(function.Func)
	recs[0].Decode(fn)

	rsp.Url = fn.Url
	return nil
}

func (e *GoogleFunction) Regions(ctx context.Context, req *function.RegionsRequest, rsp *function.RegionsResponse) error {
	rsp.Regions = GoogleRegions
	return nil
}
