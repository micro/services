package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/runtime/source/git"
	"github.com/micro/micro/v3/service/store"
	function "github.com/micro/services/function/proto"
	pauth "github.com/micro/services/pkg/auth"
	adminpb "github.com/micro/services/pkg/service/proto"
	"github.com/micro/services/pkg/tenant"
	"github.com/teris-io/shortid"
)

var (
	buildLogsRegex = regexp.MustCompile("For Cloud Build Logs, visit: https://console.cloud.google.com/cloud-build/builds;region=.*/(.*)\\?project=")
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

	*Function
}

var (
	GoogleRuntimes = []string{
		"nodejs16",
		"python39",
		"go116",
		"java11",
		"dotnet3",
		"ruby27",
		"php74",
	}

	// hardcoded list of supported regions
	GoogleRegions = []string{"europe-west1", "us-central1", "us-east1", "us-west1", "asia-east1"}

	SourceFile = map[string]string{
		"nodejs16": "index.js",
		"nodejs14": "index.js",
		"nodejs12": "index.js",
		"nodejs10": "index.js",
		"nodejs8":  "index.js",
		"nodejs6":  "index.js",
		"python39": "main.py",
		"python38": "main.py",
		"python37": "main.py",
		"go116":    "main.go",
		"go113":    "main.go",
		"go111":    "main.go",
		"java11":   "main.java",
		"dotnet3":  "main.cs",
		"ruby27":   "app.rb",
		"ruby26":   "app.rb",
		"php74":    "index.php",
	}

	GitIgnore = []string{
		".git",
		"dist",
		"node_modules",
		"vendor",
		"*.jar",
	}
)

var (
	alphanum = "abcdefghijklmnopqrstuvwxyz0123456789"
)

func random(i int) string {
	bytes := make([]byte, i)
	for {
		rand.Read(bytes)
		for i, b := range bytes {
			bytes[i] = alphanum[b%byte(len(alphanum))]
		}
		return string(bytes)
	}
	return fmt.Sprintf("%d", time.Now().Unix())
}

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
		Function: new(Function),
	}
}

func (e *GoogleFunction) WriteGcloudIgnore(dir string) error {
	data := []byte(strings.Join(GitIgnore, "\n"))
	return ioutil.WriteFile(filepath.Join(dir, ".gcloudignore"), data, 0644)
}

func (e *GoogleFunction) Deploy(ctx context.Context, req *function.DeployRequest, rsp *function.DeployResponse) error {
	log.Info("Received Function.Deploy request")

	if len(req.Name) == 0 {
		return errors.BadRequest("function.deploy", "Missing name")
	}

	if len(req.Repo) == 0 && len(req.Source) == 0 {
		return errors.BadRequest("function.deploy", "Missing source or repo")
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

	// check format of the name
	if !NameFormat.MatchString(req.Name) {
		return errors.BadRequest("function.deploy", "Invalid name")
	}

	var sourceDir string

	// if we have source then write the file
	if len(req.Source) > 0 {
		dir, err := ioutil.TempDir(".", req.Name)
		if err != nil {
			return errors.InternalServerError("function.deploy", "Failed to create source dir")
		}

		if err := os.WriteFile(filepath.Join(dir, SourceFile[req.Runtime]), []byte(req.Source), 0644); err != nil {
			return errors.InternalServerError("function.deploy", "Failed to save source")
		}
		// set the source dir
		sourceDir = dir
	} else {
		// checkout from github
		gitter := git.NewGitter(map[string]string{})

		var err error

		err = gitter.Checkout(req.Repo, req.Branch)
		if err != nil {
			return errors.InternalServerError("function.deploy", err.Error())
		}
		sourceDir = filepath.Join(gitter.RepoDir(), req.Subfolder)
	}

	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		tenantId = "micro"
	}

	if req.Entrypoint == "" {
		req.Entrypoint = req.Name
	}

	// read the function by owner
	key := fmt.Sprintf(OwnerKey+"%s/%s", tenantId, req.Name)
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
		records, err := store.Read(OwnerKey+tenantId+"/", store.ReadPrefix())
		if err != nil {
			return err
		}

		if v := len(records); v >= e.limit {
			return errors.BadRequest("function.deploy", "deployment limit reached")
		}
	}

	// set the id
	id := req.Name

	// check the owner isn't already running it
	recs, err := store.Read(FunctionKey+req.Name, store.ReadLimit(1))

	// if there's an existing function then generate a unique id
	if err == nil && len(recs) > 0 {
		// generate an id for the service
		sid, err := shortid.Generate()
		if err != nil || len(sid) == 0 {
			sid = random(8)
		}

		sid = strings.ToLower(sid)
		sid = strings.Replace(sid, "-", "", -1)
		sid = strings.Replace(sid, "_", "", -1)
		id = req.Name + "-" + sid
	}

	// process the env vars to the required format
	var envVars []string

	for k, v := range req.EnvVars {
		envVars = append(envVars, k+"="+v)
	}

	fn := &function.Func{
		Id:         id,
		Name:       req.Name,
		Source:     req.Source,
		Repo:       req.Repo,
		Subfolder:  req.Subfolder,
		Entrypoint: req.Entrypoint,
		Runtime:    req.Runtime,
		EnvVars:    req.EnvVars,
		Region:     req.Region,
		Branch:     req.Branch,
		Created:    time.Now().Format(time.RFC3339Nano),
		Status:     "Deploying",
		Url:        fmt.Sprintf("https://%s-%s.cloudfunctions.net/%s", req.Region, e.project, id),
	}

	// write the owner key
	rec := store.NewRecord(key, fn)
	if err := store.Write(rec); err != nil {
		return err
	}

	// write the global key
	rec = store.NewRecord(FunctionKey+fn.Id, fn)
	if err := store.Write(rec); err != nil {
		return err
	}

	// set the custom domain
	if len(e.domain) > 0 {
		fn.Url = fmt.Sprintf("https://%s.%s", fn.Id, e.domain)
	}

	// set the response
	rsp.Function = fn

	go func(fn *function.Func) {
		// delete the source code once done
		defer func() {
			os.RemoveAll(sourceDir)
		}()

		// https://jsoverson.medium.com/how-to-deploy-node-js-functions-to-google-cloud-8bba05e9c10a
		cmd := exec.Command("gcloud", "functions", "deploy", fn.Id, "--quiet",
			"--region", fn.Region, "--service-account", e.identity,
			"--allow-unauthenticated", "--entry-point", fn.Entrypoint,
			"--trigger-http", "--project", e.project, "--runtime", fn.Runtime)

		// if env vars exist then set them
		if len(envVars) > 0 {
			cmd.Args = append(cmd.Args, "--set-env-vars", strings.Join(envVars, ","))
		}

		cmd.Dir = sourceDir

		// write gloud ignore file
		e.WriteGcloudIgnore(cmd.Dir)

		outp, err := cmd.CombinedOutput()

		var buildID string
		reRes := buildLogsRegex.FindStringSubmatch(string(outp))
		if len(reRes) > 1 {
			buildID = reRes[1]
		}
		log.Infof("Build ID %s %s", buildID, string(outp))
		// Logs are a nice to have so don't error out
		logCmd := exec.Command("gcloud", "logging", "read", "--project", e.project, "--format", "json", fmt.Sprintf(`resource.type=build AND resource.labels.build_id=%s`, buildID))
		logOutp, logErr := logCmd.CombinedOutput()
		if logErr != nil {
			log.Errorf("Failed to retrieve logs for %s %s", buildID, logErr, string(logOutp))
		} else {
			// logs are returned in reverse chronological order as json
			var logs []map[string]interface{}
			if err := json.Unmarshal(logOutp, &logs); err != nil {
				log.Errorf("Error unmarshalling logs %s", err)
			} else {
				filteredLogs := []string{}
				for _, l := range logs {
					if tp, ok := l["textPayload"]; ok {
						filteredLogs = append(filteredLogs, tp.(string))
					}
				}
				sort.Sort(sort.Reverse(sort.StringSlice(filteredLogs)))
				// store it
				logsKey := BuildLogsKey + tenantId + "/" + req.Name

				if err := store.Write(store.NewRecord(logsKey, strings.Join(filteredLogs, "\n"))); err != nil {
					log.Errorf("Error writing logs to store %s", err)
				}
			}

		}

		if err != nil {
			log.Error(fmt.Errorf(string(outp)))
			fn.Status = "DeploymentError"
			store.Write(store.NewRecord(key, fn))
			return
		}

		var status string

	LOOP:
		// wait for the deployment and status update
		for i := 0; i < 120; i++ {
			cmd = exec.Command("gcloud", "functions", "describe", "--format", "json",
				"--region", fn.Region, "--project", e.project, fn.Id)

			outp, err = cmd.CombinedOutput()
			if err != nil {
				log.Error(fmt.Errorf(string(outp)))
				return
			}

			log.Info(string(outp))

			var m map[string]interface{}
			if err := json.Unmarshal(outp, &m); err != nil {
				log.Error(err)
				return
			}

			// write back the url
			trigger := m["httpsTrigger"].(map[string]interface{})

			if v := trigger["url"].(string); len(v) > 0 {
				fn.Url = v
			} else {
				fn.Url = fmt.Sprintf("https://%s-%s.cloudfunctions.net/%s", fn.Region, e.project, fn.Id)
			}

			v := m["status"].(string)

			switch v {
			case "ACTIVE":
				status = "Deployed"
				break LOOP
			case "DEPLOY_IN_PROGRESS":
				status = "Deploying"
			case "OFFLINE":
				status = "DeploymentError"
				break LOOP
			}

			// we need to try get the url again
			time.Sleep(time.Second)
		}

		fn.Updated = time.Now().Format(time.RFC3339Nano)
		fn.Status = status

		// write the owners key
		store.Write(store.NewRecord(key, fn))

		// write the global key
		rec = store.NewRecord(FunctionKey+fn.Id, fn)
		store.Write(rec)
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

	key := fmt.Sprintf(OwnerKey+"%s/%s", tenantId, req.Name)

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

	var sourceDir string

	// if we have source then write the file
	if len(req.Source) > 0 {
		dir, err := ioutil.TempDir(".", fn.Name)
		if err != nil {
			return errors.InternalServerError("function.update", "Failed to create source dir")
		}
		if err := os.WriteFile(filepath.Join(dir, SourceFile[fn.Runtime]), []byte(req.Source), 0644); err != nil {
			return errors.InternalServerError("function.update", "Failed to save source")
		}
		// set the source dir
		sourceDir = dir
		// set the source code
		fn.Source = req.Source
	} else {
		// checkout from github
		gitter := git.NewGitter(map[string]string{})
		if err := gitter.Checkout(fn.Repo, fn.Branch); err != nil {
			return errors.InternalServerError("function.update", err.Error())
		}
		// set the source dir
		sourceDir = filepath.Join(gitter.RepoDir(), fn.Subfolder)
	}

	// process the env vars to the required format
	var envVars []string

	for k, v := range fn.EnvVars {
		envVars = append(envVars, k+"="+v)
	}

	var status string

	go func() {
		// delete the source code once done
		defer func() {
			os.RemoveAll(sourceDir)
		}()

		// https://jsoverson.medium.com/how-to-deploy-node-js-functions-to-google-cloud-8bba05e9c10a
		cmd := exec.Command("gcloud", "functions", "deploy", fn.Id, "--quiet",
			"--region", fn.Region, "--service-account", e.identity,
			"--allow-unauthenticated", "--entry-point", fn.Entrypoint,
			"--trigger-http", "--project", e.project, "--runtime", fn.Runtime)

		// if env vars exist then set them
		if len(envVars) > 0 {
			cmd.Args = append(cmd.Args, "--set-env-vars", strings.Join(envVars, ","))
		}

		cmd.Dir = sourceDir

		// write gloud ignore file
		e.WriteGcloudIgnore(cmd.Dir)

		outp, err := cmd.CombinedOutput()

		var buildID string
		reRes := buildLogsRegex.FindStringSubmatch(string(outp))
		if len(reRes) > 1 {
			buildID = reRes[1]
		}
		log.Infof("Build ID %s %s", buildID, string(outp))

		// Logs are a nice to have so don't error out
		logCmd := exec.Command("gcloud", "logging", "read", "--project", e.project, "--format", "json", fmt.Sprintf(`resource.type=build AND resource.labels.build_id=%s`, buildID))
		logOutp, logErr := logCmd.CombinedOutput()
		if logErr != nil {
			log.Errorf("Failed to retrieve logs for %s %s", buildID, logErr, string(logOutp))
		} else {
			// logs are returned in reverse chronological order as json
			var logs []map[string]interface{}
			if err := json.Unmarshal(logOutp, &logs); err != nil {
				log.Errorf("Error unmarshalling logs %s", err)
			} else {
				filteredLogs := []string{}
				for _, l := range logs {
					if tp, ok := l["textPayload"]; ok {
						filteredLogs = append(filteredLogs, tp.(string))
					}
				}
				sort.Sort(sort.Reverse(sort.StringSlice(filteredLogs)))
				// store it
				logsKey := BuildLogsKey + tenantId + "/" + req.Name

				if err := store.Write(store.NewRecord(logsKey, strings.Join(filteredLogs, "\n"))); err != nil {
					log.Errorf("Error writing logs to store %s", err)
				}
			}

		}

		if err != nil {
			log.Error(fmt.Errorf(string(outp)))
			fn.Status = "DeploymentError"
			store.Write(store.NewRecord(key, fn))
			return
		}

	LOOP:
		// wait for the deployment and status update
		for i := 0; i < 120; i++ {
			cmd := exec.Command("gcloud", "functions", "describe", "--quiet", "--format", "json",
				"--region", fn.Region, "--project", e.project, fn.Id)

			outp, err := cmd.CombinedOutput()
			if err != nil {
				log.Error(fmt.Errorf(string(outp)))
				return
			}

			log.Info(string(outp))

			var m map[string]interface{}
			if err := json.Unmarshal(outp, &m); err != nil {
				log.Error(err)
				return
			}

			// write back the url
			trigger := m["httpsTrigger"].(map[string]interface{})
			if v := trigger["url"].(string); len(v) > 0 {
				fn.Url = v
			} else {
				fn.Url = fmt.Sprintf("https://%s-%s.cloudfunctions.net/%s", fn.Region, e.project, fn.Id)
			}

			v := m["status"].(string)

			switch v {
			case "ACTIVE":
				status = "Deployed"
				break LOOP
			case "DEPLOY_IN_PROGRESS":
				status = "Deploying"
			case "OFFLINE":
				status = "DeploymentError"
				break LOOP
			}

			// we need to try get the url again
			time.Sleep(time.Second)
		}

		// update the status
		fn.Status = status
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

	// get the function id based on the tenant
	recs, err := store.Read(OwnerKey + tenantId + "/" + req.Name)
	if err != nil {
		return err
	}
	if len(recs) == 0 {
		return errors.NotFound("function.call", "not found")
	}

	fn := new(function.Func)
	recs[0].Decode(fn)

	if len(fn.Id) == 0 {
		return errors.NotFound("function.call", "not found")
	}

	url := fn.Url

	if len(url) == 0 {
		url = fmt.Sprintf("https://%s-%s.cloudfunctions.net/%s", fn.Region, e.project, fn.Id)
	}

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

	key := fmt.Sprintf(OwnerKey+"%v/%v", tenantId, req.Name)

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
	go e.deleteFunction(fn, key)

	return nil
}

func (e *GoogleFunction) deleteFunction(fn *function.Func, key string) {
	cmd := exec.Command("gcloud", "functions", "delete", "--quiet", "--project", e.project, "--region", fn.Region, fn.Id)
	outp, err := cmd.CombinedOutput()
	if err != nil && !strings.Contains(string(outp), "does not exist") {
		log.Error(fmt.Errorf(string(outp)))
		return
	}

	// delete the owner key
	store.Delete(key)

	// delete the global key
	store.Delete(FunctionKey + fn.Id)
}

func (e *GoogleFunction) List(ctx context.Context, req *function.ListRequest, rsp *function.ListResponse) error {
	log.Info("Received Function.List request")

	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		tenantId = "micro"
	}

	key := OwnerKey + tenantId + "/"

	records, err := store.Read(key, store.ReadPrefix())
	if err != nil {
		return err
	}

	for _, record := range records {
		fn := new(function.Func)
		if err := record.Decode(fn); err != nil {
			continue
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

	key := fmt.Sprintf(OwnerKey+"%v/%v", tenantId, req.Name)

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
	cmd := exec.Command("gcloud", "functions", "describe", "--format", "json", "--region", fn.Region, "--project", e.project, fn.Id)
	outp, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(fmt.Errorf(string(outp)))
		return fmt.Errorf("function does not exist")
	}

	log.Info(string(outp))
	m := map[string]interface{}{}

	if err := json.Unmarshal(outp, &m); err != nil {
		return err
	}

	// set describe info
	status := m["status"].(string)
	status = strings.Replace(status, "_", " ", -1)
	status = strings.Title(strings.ToLower(status))
	fn.Status = status

	// set the url
	if len(fn.Url) == 0 && status == "Deployed" {
		v := m["httpsTrigger"].(map[string]interface{})
		fn.Url = v["url"].(string)
	}

	// write it back
	go store.Write(store.NewRecord(key, fn))

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

	recs, err := store.Read(FunctionKey+req.Id, store.ReadLimit(1))
	if err != nil {
		return err
	}

	if len(recs) == 0 {
		return errors.BadRequest("function.proxy", "function does not exist")
	}

	fn := new(function.Func)
	recs[0].Decode(fn)

	url := fn.Url

	// backup plan is to construct https://region-project.cloudfunctions.net/function-name
	if len(url) == 0 {
		url = fmt.Sprintf("https://%s-%s.cloudfunctions.net/%s", fn.Region, g.project, fn.Id)
	}

	rsp.Url = url

	return nil
}

func (e *GoogleFunction) Regions(ctx context.Context, req *function.RegionsRequest, rsp *function.RegionsResponse) error {
	rsp.Regions = GoogleRegions
	return nil
}

func (e *GoogleFunction) Runtimes(ctx context.Context, req *function.RuntimesRequest, rsp *function.RuntimesResponse) error {
	rsp.Runtimes = GoogleRuntimes
	return nil
}

func (e *GoogleFunction) DeleteData(ctx context.Context, request *adminpb.DeleteDataRequest, response *adminpb.DeleteDataResponse) error {
	method := "admin.DeleteData"
	_, err := pauth.VerifyMicroAdmin(ctx, method)
	if err != nil {
		return err
	}

	if len(request.TenantId) < 10 { // deliberate length check so we don't delete all the things
		return errors.BadRequest(method, "Missing tenant ID")
	}

	prefix := OwnerKey + request.TenantId + "/"

	recs, err := store.Read(prefix, store.ReadPrefix())
	if err != nil {
		return err
	}

	for _, rec := range recs {
		var fn function.Func
		if err := rec.Decode(&fn); err != nil {
			return err
		}
		e.deleteFunction(&fn, rec.Key)
	}
	log.Infof("Deleted %d functions for %s", len(recs), request.TenantId)
	return nil
}

var (
	logsFuncMap = map[string]func(e *GoogleFunction, ctx context.Context, req *function.LogsRequest, rsp *function.LogsResponse) error{
		"build": buildLogs, // TODO add runtime logs
	}
)

func (e *GoogleFunction) Logs(ctx context.Context, req *function.LogsRequest, rsp *function.LogsResponse) error {
	f, ok := logsFuncMap[req.LogsType]
	if !ok {
		return errors.BadRequest("app.Logs", "Invalid logs_type specified")
	}
	return f(e, ctx, req, rsp)
}

func buildLogs(e *GoogleFunction, ctx context.Context, req *function.LogsRequest, rsp *function.LogsResponse) error {
	id, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.Unauthorized("app.Logs", "Unauthorized")
	}
	logsKey := BuildLogsKey + id + "/" + req.Name
	recs, err := store.Read(logsKey)
	if err != nil {
		return errors.NotFound("app.Logs", "Build logs not found")
	}
	var ret string
	if err := json.Unmarshal(recs[0].Value, &ret); err != nil {
		log.Errorf("Error unmarshalling logs %s", err)
		return errors.NotFound("app.Logs", "Build logs not found")
	}
	rsp.Logs = ret
	return nil
}
