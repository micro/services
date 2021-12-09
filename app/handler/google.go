package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os/exec"
	"strings"
	"time"

	_struct "github.com/golang/protobuf/ptypes/struct"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/context/metadata"
	"github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/runtime/source/git"
	pb "github.com/micro/services/app/proto"
	db "github.com/micro/services/db/proto"
	"github.com/teris-io/shortid"
)

type GoogleApp struct {
	project string
	// eg. https://us-central1-m3o-apis.cloudfunctions.net/
	address string
	db      db.DbService
	limit int
	regions []string

	// Embed the app handler
	*App
}

var (
	// hardcoded list of supported regions
	GoogleRegions = []string{"europe-west1", "us-east1", "us-west1"}

	// hardcoded list of valid repos
	GitRepos = []string{"github.com", "gitlab.org", "bitbucket.org"}
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

func New(db db.DbService) *GoogleApp {
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

	return &GoogleApp{project: project, address: address, db: db, limit: limit, App: new(App)}
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

	gitter := git.NewGitter(map[string]string{})
	err := gitter.Checkout(req.Repo, req.Branch)
	if err != nil {
		log.Errorf("Failed to download %s@%s\n", req.Repo, req.Branch)
		return errors.InternalServerError("app.run", "Failed to download source")
	}

	// check for the existing app
	readRsp, err := e.db.Read(ctx, &db.ReadRequest{
		Table: "apps",
		Id:    req.Name,
	})
	if err != nil {
		return err
	}

	// app is already running
	if len(readRsp.Records) > 0 {
		return errors.BadRequest("app.run", "%s already exists", req.Name)
	}

	// check for app limit
	if e.limit > 0 {
		// check for the existing app
		countRsp, err := e.db.Count(ctx, &db.CountRequest{
			Table: "apps",
		})
		if err != nil {
			return err
		}

		if int(countRsp.Count) >= e.limit {
			return errors.BadRequest("app.run", "deployment limit reached")
		}
	}

	// TODO validate name and use custom domain name

	// process the env vars to the required format
	var envVars []string

	for k, v := range req.EnvVars {
		envVars = append(envVars, k+"="+v)
	}

	// generate an id for the service
	sid, err := shortid.Generate()
	if err != nil || len(sid) == 0 {
		sid = random(8)
	}

	sid = strings.ToLower(sid)
	id := req.Name + "-" + strings.Replace(sid, "-", "", -1)

	service := &pb.Service{
		Name:    req.Name,
		Id:      id,
		Repo:    req.Repo,
		Branch:  req.Branch,
		Region:  req.Region,
		Port:    req.Port,
		Status:  "Deploying",
		EnvVars: req.EnvVars,
	}

	s := &_struct.Struct{}
	b, _ := json.Marshal(service)
	err = s.UnmarshalJSON(b)
	if err != nil {
		return err
	}

	// write the app to the db
	_, err = e.db.Create(ctx, &db.CreateRequest{
		Table:  "apps",
		Record: s,
		Id:     req.Name,
	})
	if err != nil {
		log.Error(err)
	}

	go func(service *pb.Service) {
		// https://jsoverson.medium.com/how-to-deploy-node-js-functions-to-google-cloud-8bba05e9c10a
		cmd := exec.Command("gcloud", "--project", e.project, "--format", "json", "run", "deploy", service.Id, "--region", req.Region,
			"--cpu", "1", "--memory", "256Mi", "--port", fmt.Sprintf("%d", req.Port),
			"--allow-unauthenticated", "--max-instances", "1", "--source", ".",
		)

		// if env vars exist then set them
		if len(envVars) > 0 {
			cmd.Args = append(cmd.Args, "--set-env-vars", strings.Join(envVars, ","))
		}

		// set the command dir
		cmd.Dir = gitter.RepoDir()

		// execute the command
		outp, err := cmd.CombinedOutput()

		// by this point the context may have been cancelled
		acc, _ := auth.AccountFromContext(ctx)
		md, _ := metadata.FromContext(ctx)

		ctx = metadata.NewContext(context.Background(), md)
		ctx = auth.ContextWithAccount(ctx, acc)

		if err == nil {
			// populate the app status
			e.Status(ctx, &pb.StatusRequest{Name: req.Name}, &pb.StatusResponse{})
			return
		}


		errString := string(outp)

		log.Error(fmt.Errorf(errString))

		// set the error status
		service.Status = "DeploymentError"

		if strings.Contains(errString, "Failed to start and then listen on the port defined by the PORT environment variable") {
			service.Status += ": Failed to start and listen on port " + fmt.Sprintf("%d", req.Port)
		} else if strings.Contains(errString, "Build failed") {
			service.Status += ": Build failed"
		}

		// crazy garbage structs
		s := &_struct.Struct{}
		b, _ := json.Marshal(service)

		if err = s.UnmarshalJSON(b); err != nil {
			log.Error(err)
			return
		}

		// write the app to the db
		_, err = e.db.Update(ctx, &db.UpdateRequest{
			Table:  "apps",
			Record: s,
			Id:     req.Name,
		})
		if err != nil {
			log.Error(err)
		}
	}(service)

	// set the service in the response
	rsp.Service = service

	return err
}

func (e *GoogleApp) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	log.Info("Received App.Delete request")

	// read the app from the db
	readRsp, err := e.db.Read(ctx, &db.ReadRequest{
		Table: "apps",
		Id:    req.Name,
	})
	if err != nil {
		return err
	}

	// not running
	if len(readRsp.Records) == 0 {
		return nil
	}

	record := readRsp.Records[0]
	b, _ := json.Marshal(record.AsMap())
	srv := new(pb.Service)

	if err := json.Unmarshal(b, srv); err != nil {
		return err
	}

	// execute the delete async
	// delete the app
	cmd := exec.Command("gcloud", "--quiet", "--project", e.project, "run", "services", "delete", "--region", srv.Region, srv.Id)
	outp, err := cmd.CombinedOutput()
	if err != nil && !strings.Contains(string(outp), "could not be found") {
		log.Error(fmt.Errorf(string(outp)))
		return err
	}

	// delete from the db
	_, err = e.db.Delete(ctx, &db.DeleteRequest{
		Table: "apps",
		Id:    req.Name,
	})
	return err
}

func (e *GoogleApp) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	log.Info("Received App.List request")

	readRsp, err := e.db.Read(ctx, &db.ReadRequest{
		Table: "apps",
	})
	if err != nil {
		return err
	}

	rsp.Services = []*pb.Service{}

	for _, record := range readRsp.Records {
		b, _ := json.Marshal(record.AsMap())
		srv := new(pb.Service)

		if err = json.Unmarshal(b, srv); err != nil {
			return err
		}

		rsp.Services = append(rsp.Services, srv)
	}

	return nil
}

func (e *GoogleApp) Status(ctx context.Context, req *pb.StatusRequest, rsp *pb.StatusResponse) error {
	readRsp, err := e.db.Read(ctx, &db.ReadRequest{
		Table: "apps",
		Id:    req.Name,
	})
	if err != nil {
		return err
	}

	if len(readRsp.Records) == 0 {
		return errors.NotFound("app.status", "app not found")
	}

	srv := new(pb.Service)

	b, _ := json.Marshal(readRsp.Records[0].AsMap())
	if err = json.Unmarshal(b, srv); err != nil {
		return err
	}

	// get the current app status
	cmd := exec.Command("gcloud", "--project", e.project, "--format", "json", "run", "services", "describe", "--region", srv.Region, srv.Id)
	outp, err := cmd.CombinedOutput()
	if err != nil && srv.Status == "Deploying" {
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
	deployedAt := srv.DeployedAt

	// get the service status
	status := output["status"].(map[string]interface{})
	deployed := status["conditions"].([]interface{})[0].(map[string]interface{})
	srv.DeployedAt = deployed["lastTransitionTime"].(string)

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
	if srv.Status == currentStatus && srv.Url == currentUrl && srv.DeployedAt == deployedAt {
		return nil
	}

	// update built in status

	s := &_struct.Struct{}
	b, _ = json.Marshal(srv)

	if err = s.UnmarshalJSON(b); err != nil {
		log.Error(err)
		return err
	}

	// write the app to the db
	_, err = e.db.Update(ctx, &db.UpdateRequest{
		Table:  "apps",
		Record: s,
		Id:     req.Name,
	})

	return err
}