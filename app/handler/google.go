package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/context/metadata"
	"github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	"github.com/micro/services/app/domain"
	pb "github.com/micro/services/app/proto"
	github "github.com/micro/services/github/proto"
	pauth "github.com/micro/services/pkg/auth"
	"github.com/micro/services/pkg/git"
	adminpb "github.com/micro/services/pkg/service/proto"
	"github.com/micro/services/pkg/tenant"
	"github.com/teris-io/shortid"
)

var (
	buildLogsRegex = regexp.MustCompile("Logs are available at \\[https://console.cloud.google.com/cloud-build/builds/(.*)\\?project=")
)

type GoogleApp struct {
	// the associated google project
	project string
	// eg. https://us-central1-m3o-apis.cloudfunctions.net/
	address string
	// max number of apps per user
	limit int
	// custom domain for apps
	domain string
	// the service account for the app
	identity string
	// Embed the app handler
	*App

	ghSvc github.GithubService
}

var (
	// hardcoded list of supported regions
	GoogleRegions = []string{"europe-west1", "us-central1", "us-east1", "us-west1", "asia-east1"}

	// hardcoded list of valid repos
	GitRepos = []string{"github.com", "gitlab.org", "bitbucket.org"}

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

func New(svc *service.Service) *GoogleApp {
	v, err := config.Get("app.service_account_json")
	if err != nil {
		log.Fatalf("app.service_account_json: %v", err)
	}
	keyfile := v.Bytes()
	if len(keyfile) == 0 {
		log.Fatalf("empty keyfile")
	}

	v, err = config.Get("app.address")
	if err != nil {
		log.Fatalf("app.address: %v", err)
	}
	address := v.String("")
	if len(address) == 0 {
		log.Fatalf("empty address")
	}

	v, err = config.Get("app.domain")
	if err != nil {
		log.Fatalf("app.domain: %v", err)
	}
	domain := v.String("")

	v, err = config.Get("app.project")
	if err != nil {
		log.Fatalf("app.project: %v", err)
	}
	project := v.String("")
	if len(project) == 0 {
		log.Fatalf("empty project")
	}
	v, err = config.Get("app.limit")
	if err != nil {
		log.Fatalf("app.limit: %v", err)
	}
	limit := v.Int(0)
	if limit == 0 {
		log.Infof("App limit is %d", limit)
	}

	v, err = config.Get("app.service_account")
	if err != nil {
		log.Fatalf("app.service_account: %v", err)
	}
	accName := v.String("")

	v, err = config.Get("app.service_identity")
	if err != nil {
		log.Fatalf("app.service_identity: %v", err)
	}
	identity := v.String("")

	m := map[string]interface{}{}
	err = json.Unmarshal(keyfile, &m)
	if err != nil {
		log.Fatalf("invalid json: %v", err)
	}

	// only root
	err = ioutil.WriteFile("account.json", keyfile, 0700)
	if err != nil {
		log.Fatalf("app.service_account: %v", err)
	}

	// DO THIS STEP
	// https://cloud.google.com/functions/docs/reference/iam/roles#additional-configuration

	// https://cloud.google.com/sdk/docs/authorizing#authorizing_with_a_service_account
	outp, err := exec.Command("gcloud", "auth", "activate-service-account", accName, "--key-file", "account.json").CombinedOutput()
	if err != nil {
		log.Fatal(string(outp), err.Error())
	}

	outp, err = exec.Command("gcloud", "auth", "list").CombinedOutput()
	if err != nil {
		log.Fatal(string(outp), err.Error())
	}
	log.Info(string(outp))

	return &GoogleApp{
		domain:   domain,
		project:  project,
		address:  address,
		limit:    limit,
		identity: identity,
		App:      new(App),
		ghSvc:    github.NewGithubService("github", svc.Client()),
	}
}

func (e *GoogleApp) WriteGcloudIgnore(dir string) error {
	data := []byte(strings.Join(GitIgnore, "\n"))
	return ioutil.WriteFile(filepath.Join(dir, ".gcloudignore"), data, 0644)
}

func (e *GoogleApp) Regions(ctx context.Context, req *pb.RegionsRequest, rsp *pb.RegionsResponse) error {
	rsp.Regions = GoogleRegions
	return nil
}

func (e *GoogleApp) Run(ctx context.Context, req *pb.RunRequest, rsp *pb.RunResponse) error {
	log.Info("Received App.Run request")

	if len(req.Name) == 0 {
		return errors.BadRequest("app.run", "missing name")
	}

	if len(req.Repo) == 0 {
		return errors.BadRequest("app.run", "missing repo")
	}

	if req.Port <= 0 {
		req.Port = 8080
	}

	// validations
	if !NameFormat.MatchString(req.Name) {
		return errors.BadRequest("app.run", "invalidate name format")
	}

	id, ok := tenant.FromContext(ctx)
	if !ok {
		id = "micro"
	}

	var reservedApp bool

	// check if the name is reserved by a different owner
	reservedKey := ReservationKey + req.Name

	recs, err := store.Read(reservedKey, store.ReadLimit(1))
	if err == nil && len(recs) > 0 {
		res := new(Reservation)
		recs[0].Decode(res)

		// if its the owners app and there's still time left
		if res.Owner == id && res.Expires.After(time.Now()) {
			reservedApp = true
		}
	}

	var validRepo bool

	// only support github and gitlab
	for _, repo := range GitRepos {
		rp := repo

		if strings.HasPrefix(req.Repo, "https://") {
			rp = "https://" + repo
		}

		if strings.HasPrefix(req.Repo, rp) {
			validRepo = true
			break
		}
	}

	if !validRepo {
		return errors.BadRequest("app.run", "invalid git repo")
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
		return errors.BadRequest("app.run", "Unsupported region")
	}

	// checkout the source code
	if len(req.Branch) == 0 {
		req.Branch = "master"
	}

	// look for the existing service
	key := OwnerKey + id + "/" + req.Name

	// check the owner isn't already running it
	recs, err = store.Read(key, store.ReadLimit(1))
	if err == nil && len(recs) > 0 {
		return errors.BadRequest("app.run", "%s already running", req.Name)
	}

	// check the global namespace
	// look for the existing service
	key = ServiceKey + req.Name

	// set the id
	appId := req.Name

	// check the app isn't already running it
	recs, err = store.Read(key, store.ReadLimit(1))

	// if there's an existing service then generate a unique id
	if err == nil && len(recs) > 0 {
		// generate an id for the service
		sid, err := shortid.Generate()
		if err != nil || len(sid) == 0 {
			sid = random(8)
		}

		sid = strings.ToLower(sid)
		sid = strings.Replace(sid, "-", "", -1)
		sid = strings.Replace(sid, "_", "", -1)
		appId = req.Name + "-" + sid
	}

	// check for app limit
	if e.limit > 0 && !reservedApp {
		ownerKey := OwnerKey + id + "/"
		recs, err := store.Read(ownerKey, store.ReadPrefix())
		if err == nil && len(recs) >= e.limit {
			return errors.BadRequest("app.run", "deployment limit reached")
		}
	}

	// checkout the code
	gitter := git.NewGitter(e.gitCreds(ctx, id))
	if err := gitter.Checkout(req.Repo, req.Branch); err != nil {
		log.Errorf("Failed to download %s@%s\n", req.Repo, req.Branch)
		return errors.InternalServerError("app.run", "Failed to download source")
	}

	// TODO validate name and use custom domain name

	// process the env vars to the required format
	var envVars []string

	for k, v := range req.EnvVars {
		envVars = append(envVars, k+"="+v)
	}

	service := &pb.Service{
		Name:    req.Name,
		Id:      appId,
		Repo:    req.Repo,
		Branch:  req.Branch,
		Region:  req.Region,
		Port:    req.Port,
		Status:  domain.StatusDeploying,
		EnvVars: req.EnvVars,
		Created: time.Now().Format(time.RFC3339Nano),
	}

	keys := []string{
		// service key
		ServiceKey + service.Id,
		// owner key
		OwnerKey + id + "/" + req.Name,
	}

	// function to update records
	updateRecords := func(service *pb.Service) {
		// write the keys for the service
		for _, key := range keys {
			rec := store.NewRecord(key, service)

			if err := store.Write(rec); err != nil {
				log.Error(err)
				return
			}
		}
	}

	// write the record
	updateRecords(service)

	// make copy
	srv := new(pb.Service)
	*srv = *service

	// set the custom domain
	if len(e.domain) > 0 {
		srv.Url = fmt.Sprintf("https://%s.%s", srv.Id, e.domain)
	}

	// set the service in the response
	rsp.Service = srv

	go func(service *pb.Service) {
		imageName := fmt.Sprintf("%s-docker.pkg.dev/%s/cloud-run-source-deploy/%s", req.Region, e.project, service.Id)
		cmd := exec.Command("gcloud", "builds", "submit", "--project", e.project, "--format", "json",
			"--pack", "image="+imageName, ".",
		)

		// set the command dir
		cmd.Dir = gitter.RepoDir()
		// write the gloudignore file
		e.WriteGcloudIgnore(cmd.Dir)

		// execute the command
		outp, err := cmd.CombinedOutput()

		var buildID string
		reRes := buildLogsRegex.FindStringSubmatch(string(outp))
		if len(reRes) > 1 {
			buildID = reRes[1]
		}

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
				logsKey := BuildLogsKey + id + "/" + req.Name

				if err := store.Write(store.NewRecord(logsKey, strings.Join(filteredLogs, "\n"))); err != nil {
					log.Errorf("Error writing logs to store %s", err)
				}
			}

		}

		// generate a unique service account for the app
		// https://jsoverson.medium.com/how-to-deploy-node-js-functions-to-google-cloud-8bba05e9c10a
		cmd = exec.Command("gcloud", "--project", e.project, "--quiet", "--format", "json", "run",
			"deploy", service.Id, "--region", req.Region,
			"--service-account", e.identity,
			"--cpu", "1", "--memory", "256Mi", "--port", fmt.Sprintf("%d", req.Port),
			"--allow-unauthenticated", "--max-instances", "1", "--image", imageName,
		)

		// if env vars exist then set them
		if len(envVars) > 0 {
			cmd.Args = append(cmd.Args, "--set-env-vars", strings.Join(envVars, ","))
		}
		// set the command dir
		cmd.Dir = gitter.RepoDir()
		// execute the command
		outp, err = cmd.CombinedOutput()

		// by this point the context may have been cancelled
		acc, _ := auth.AccountFromContext(ctx)
		md, _ := metadata.FromContext(ctx)

		ctx = metadata.NewContext(context.Background(), md)
		ctx = auth.ContextWithAccount(ctx, acc)

		if err == nil {
			// populate the app status
			err = e.Status(ctx, &pb.StatusRequest{Name: req.Name}, &pb.StatusResponse{})
			if err != nil {
				log.Error(err)
			}
			return
		}

		errString := string(outp)

		log.Error(fmt.Errorf(errString))

		// set the error status
		service.Status = domain.StatusDeploymentError

		if strings.Contains(errString, "Failed to start and then listen on the port defined by the PORT environment variable") {
			service.Status += ": Failed to start and listen on port " + fmt.Sprintf("%d", req.Port)
		} else if strings.Contains(errString, "Build failed") {
			service.Status += ": Build failed"
		}

		// write the records
		updateRecords(service)
	}(service)

	return nil
}

func (e *GoogleApp) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	log.Info("Received App.Update request")

	if len(req.Name) == 0 {
		return errors.BadRequest("app.update", "missing name")
	}

	// validations
	if !NameFormat.MatchString(req.Name) {
		return errors.BadRequest("app.run", "invalidate name format")
	}

	id, ok := tenant.FromContext(ctx)
	if !ok {
		id = "micro"
	}

	key := OwnerKey + id + "/" + req.Name

	// look for the existing service
	recs, err := store.Read(key, store.ReadLimit(1))
	if err != nil && err == store.ErrNotFound {
		return errors.BadRequest("app.update", "%s does not exist", req.Name)
	}

	if len(recs) == 0 {
		return errors.BadRequest("app.update", "%s does not exist", req.Name)
	}

	srv := new(pb.Service)

	if err := recs[0].Decode(srv); err != nil {
		return err
	}

	// check the status
	switch srv.Status {
	case domain.StatusUpdating, domain.StatusDeploying, domain.StatusDeleting:
		log.Errorf("Won't update: % is %s", req.Name, srv.Status)
		return errors.BadRequest("app.update", "% status: %s", req.Name, srv.Status)
	}

	// checkout the code
	gitter := git.NewGitter(e.gitCreds(ctx, id))
	if err := gitter.Checkout(srv.Repo, srv.Branch); err != nil {
		log.Errorf("Failed to download %s@%s\n", srv.Repo, srv.Branch)
		return errors.InternalServerError("app.run", "Failed to download source")
	}

	// TODO validate name and use custom domain name

	// process the env vars to the required format
	var envVars []string

	for k, v := range srv.EnvVars {
		// don't add anything we're updating
		if _, ok := req.EnvVars[k]; ok {
			continue
		}
		// add the value
		envVars = append(envVars, k+"="+v)
	}

	// add the updated env vars
	for k, v := range req.EnvVars {
		// skip zero vals
		if len(v) == 0 {
			// delete it
			delete(srv.EnvVars, k)
			continue
		}
		// update the env vars
		srv.EnvVars[k] = v
		envVars = append(envVars, k+"="+v)
	}

	keys := []string{
		// service key
		ServiceKey + srv.Id,
		// owner key
		OwnerKey + id + "/" + req.Name,
	}

	// function to update records
	updateRecords := func(service *pb.Service) {
		// write the keys for the service
		for _, key := range keys {
			rec := store.NewRecord(key, service)

			if err := store.Write(rec); err != nil {
				log.Error(err)
				return
			}
		}
	}

	srv.Status = domain.StatusUpdating

	// update immediately with status
	updateRecords(srv)

	go func(service *pb.Service) {
		imageName := fmt.Sprintf("%s-docker.pkg.dev/%s/cloud-run-source-deploy/%s", service.Region, e.project, service.Id)
		cmd := exec.Command("gcloud", "builds", "submit", "--project", e.project, "--format", "json",
			"--pack", "image="+imageName, ".",
		)

		// set the command dir
		cmd.Dir = gitter.RepoDir()
		// write the gloudignore file
		e.WriteGcloudIgnore(cmd.Dir)

		// execute the command
		outp, err := cmd.CombinedOutput()

		var buildID string
		reRes := buildLogsRegex.FindStringSubmatch(string(outp))
		if len(reRes) > 1 {
			buildID = reRes[1]
		}

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
				logsKey := BuildLogsKey + id + "/" + req.Name

				if err := store.Write(store.NewRecord(logsKey, strings.Join(filteredLogs, "\n"))); err != nil {
					log.Errorf("Error writing logs to store %s", err)
				}
			}

		}

		// generate a unique service account for the app
		// https://jsoverson.medium.com/how-to-deploy-node-js-functions-to-google-cloud-8bba05e9c10a
		cmd = exec.Command("gcloud", "--project", e.project, "--quiet", "--format", "json", "run",
			"deploy", service.Id, "--region", service.Region,
			"--service-account", e.identity,
			"--cpu", "1", "--memory", "256Mi", "--port", fmt.Sprintf("%d", service.Port),
			"--allow-unauthenticated", "--max-instances", "1", "--image", imageName,
		)

		// if env vars exist then set them
		if len(envVars) > 0 {
			cmd.Args = append(cmd.Args, "--set-env-vars", strings.Join(envVars, ","))
		}
		// set the command dir
		cmd.Dir = gitter.RepoDir()
		// execute the command
		outp, err = cmd.CombinedOutput()

		// by this point the context may have been cancelled
		acc, _ := auth.AccountFromContext(ctx)
		md, _ := metadata.FromContext(ctx)

		ctx = metadata.NewContext(context.Background(), md)
		ctx = auth.ContextWithAccount(ctx, acc)

		if err == nil {
			// populate the app status
			err = e.Status(ctx, &pb.StatusRequest{Name: service.Name}, &pb.StatusResponse{})
			if err != nil {
				log.Error(err)
			}
			return
		}

		errString := string(outp)

		log.Error(fmt.Errorf(errString))

		// set the error status
		service.Status = domain.StatusDeploymentError

		if strings.Contains(errString, "Failed to start and then listen on the port defined by the PORT environment variable") {
			service.Status += ": Failed to start and listen on port " + fmt.Sprintf("%d", service.Port)
		} else if strings.Contains(errString, "Build failed") {
			service.Status += ": Build failed"
		}

		// update the records
		updateRecords(service)
	}(srv)

	return err
}

func (e *GoogleApp) gitCreds(ctx context.Context, tenantID string) map[string]string {
	creds := map[string]string{}
	tokRsp, _ := e.ghSvc.Token(ctx, &github.TokenRequest{TenantId: tenantID}, client.WithAuthToken())
	if tokRsp != nil {
		creds[git.CredentialsKey] = tokRsp.Token
	}
	return creds
}

func (e *GoogleApp) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	log.Info("Received App.Delete request")

	if len(req.Name) == 0 {
		return errors.BadRequest("app.delete", "missing name")
	}

	id, ok := tenant.FromContext(ctx)
	if !ok {
		id = "micro"
	}
	// read the app for the owner
	key := OwnerKey + id + "/" + req.Name
	recs, err := store.Read(key)
	if err != nil {
		return err
	}

	// not running
	if len(recs) == 0 {
		return nil
	}

	// decode the service
	srv := new(pb.Service)

	if err := recs[0].Decode(srv); err != nil {
		return err
	}

	return e.deleteApp(ctx, id, srv)
}

func (e *GoogleApp) deleteApp(ctx context.Context, tenantID string, srv *pb.Service) error {

	// check the status
	switch srv.Status {
	case domain.StatusUpdating, domain.StatusDeploying, domain.StatusDeleting:
		log.Errorf("Won't delete: %s is %s", srv.Name, srv.Status)
		return errors.BadRequest("app.delete", "%s status: %s", srv.Name, srv.Status)
	}

	// delete from the db
	keys := []string{
		// service key
		ServiceKey + srv.Id,
		// owner key
		OwnerKey + tenantID + "/" + srv.Name,
	}

	// set the delete status
	srv.Status = domain.StatusDeleting

	// update the status
	for _, key := range keys {
		rec := store.NewRecord(key, srv)
		if err := store.Write(rec); err != nil {
			log.Error(err)
			return err
		}
	}

	// execute the delete async
	go func(srv *pb.Service) {
		cmd := exec.Command("gcloud", "--quiet", "--project", e.project, "run", "services", "delete", "--region", srv.Region, srv.Id)
		outp, err := cmd.CombinedOutput()
		if err != nil && !strings.Contains(string(outp), "could not be found") {
			log.Error(fmt.Errorf(string(outp)))
			return
		}

		// delete the records
		for _, key := range keys {
			if err := store.Delete(key); err != nil {
				log.Error(err)
				return
			}
		}

		store.Delete(BuildLogsKey + tenantID + "/" + srv.Id)
	}(srv)

	return nil
}

func (e *GoogleApp) DeleteData(ctx context.Context, request *adminpb.DeleteDataRequest, response *adminpb.DeleteDataResponse) error {
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
		var srv pb.Service
		if err := rec.Decode(&srv); err != nil {
			return err
		}
		e.deleteApp(ctx, request.TenantId, &srv)
	}
	log.Infof("Deleted %d functions for %s", len(recs), request.TenantId)
	return nil
}

func (e *GoogleApp) Usage(ctx context.Context, request *adminpb.UsageRequest, response *adminpb.UsageResponse) error {
	return nil
}

func (e *GoogleApp) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	log.Info("Received App.List request")

	id, ok := tenant.FromContext(ctx)
	if !ok {
		id = "micro"
	}

	key := OwnerKey + id + "/"

	recs, err := store.Read(key, store.ReadPrefix())
	if err != nil {
		return err
	}

	rsp.Services = []*pb.Service{}

	for _, rec := range recs {
		srv := new(pb.Service)
		if err := rec.Decode(srv); err != nil {
			continue
		}

		// set the custom domain
		if len(e.domain) > 0 {
			// set the backend
			srv.Backend = srv.Url
			// vanity url
			srv.Url = fmt.Sprintf("https://%s.%s", srv.Id, e.domain)
		}

		rsp.Services = append(rsp.Services, srv)
	}

	return nil
}

func (e *GoogleApp) Status(ctx context.Context, req *pb.StatusRequest, rsp *pb.StatusResponse) error {
	id, ok := tenant.FromContext(ctx)
	if !ok {
		id = "micro"
	}

	key := OwnerKey + id + "/" + req.Name

	recs, err := store.Read(key)
	if err != nil {
		return err
	}

	if len(recs) == 0 {
		return errors.NotFound("app.status", "app not found")
	}

	srv := new(pb.Service)
	if err := recs[0].Decode(srv); err != nil {
		return err
	}

	// get the current app status
	cmd := exec.Command("gcloud", "--project", e.project, "--format", "json", "run", "services", "describe", "--region", srv.Region, srv.Id)
	outp, err := cmd.CombinedOutput()
	if err != nil && srv.Status == "Deploying" {
		log.Error(fmt.Errorf(string(outp)))
		rsp.Service = srv
		return nil
	} else if strings.HasPrefix(srv.Status, "DeploymentError") {
		log.Error(fmt.Errorf(string(outp)))
		rsp.Service = srv
		return nil
	} else if err != nil {
		log.Error(fmt.Errorf(string(outp)))
		return errors.BadRequest("app.status", "service does not exist")
	}

	var output map[string]interface{}
	if err = json.Unmarshal(outp, &output); err != nil {
		return err
	}

	currentStatus := srv.Status
	currentUrl := srv.Url
	updatedAt := srv.Updated

	// get the service status
	status := output["status"].(map[string]interface{})
	deployed := status["conditions"].([]interface{})[0].(map[string]interface{})
	srv.Updated = deployed["lastTransitionTime"].(string)

	switch deployed["status"].(string) {
	case "True":
		srv.Status = "Running"
		srv.Url = status["url"].(string)
	default:
		srv.Status = deployed["status"].(string)
	}

	// set response
	rsp.Service = srv

	// no change in status and we have a pre-existing url
	if srv.Status == currentStatus && srv.Url == currentUrl && srv.Updated == updatedAt {
		// set the custom domain
		if len(e.domain) > 0 {
			// set the backend
			rsp.Service.Backend = srv.Url
			rsp.Service.Url = fmt.Sprintf("https://%s.%s", srv.Id, e.domain)
		}
		return nil
	}

	// update built in status
	// delete from the db
	keys := []string{
		// global key
		ServiceKey + srv.Id,
		// owner key
		OwnerKey + id + "/" + req.Name,
	}

	// delete the keys for the service
	for _, key := range keys {
		rec := store.NewRecord(key, srv)
		// write the app to the db
		if err := store.Write(rec); err != nil {
			log.Error(err)
			return err
		}
	}

	// set the custom domain
	if len(e.domain) > 0 {
		// set the backend
		rsp.Service.Backend = srv.Url
		rsp.Service.Url = fmt.Sprintf("https://%s.%s", srv.Id, e.domain)
	}

	return nil
}

func (a *App) Resolve(ctx context.Context, req *pb.ResolveRequest, rsp *pb.ResolveResponse) error {
	if len(req.Id) == 0 {
		return errors.BadRequest("app.resolve", "missing id")
	}

	key := ServiceKey + req.Id

	recs, err := store.Read(key, store.ReadLimit(1))
	if err != nil {
		return err
	}

	if len(recs) == 0 {
		return errors.BadRequest("app.resolve", "app does not exist")
	}

	srv := new(pb.Service)
	recs[0].Decode(srv)

	rsp.Url = srv.Url
	return nil
}

var (
	logsFuncMap = map[string]func(e *GoogleApp, ctx context.Context, req *pb.LogsRequest, rsp *pb.LogsResponse) error{
		"build": buildLogs, // TODO add runtime logs
	}
)

func (e *GoogleApp) Logs(ctx context.Context, req *pb.LogsRequest, rsp *pb.LogsResponse) error {
	f, ok := logsFuncMap[req.LogsType]
	if !ok {
		return errors.BadRequest("app.Logs", "Invalid logs_type specified")
	}
	return f(e, ctx, req, rsp)
}

func buildLogs(e *GoogleApp, ctx context.Context, req *pb.LogsRequest, rsp *pb.LogsResponse) error {
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
