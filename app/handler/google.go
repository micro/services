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
	"github.com/micro/micro/v3/service/config"
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
	err = ioutil.WriteFile("/acc.json", keyfile, 0700)
	if err != nil {
		log.Fatalf("app.service_account: %v", err)
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

	return &GoogleApp{project: project, address: address, db: db, App: new(App)}
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
		if strings.HasPrefix(req.Repo, "https://"+repo) {
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

	gitter := git.NewGitter(map[string]string{})
	err := gitter.Checkout(req.Repo, "master")
	if err != nil {
		return err
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

	id := req.Name + "-" + strings.Replace(sid, "-", "", -1)

	service := &pb.Service{
		Name:    req.Name,
		Id:      id,
		Repo:    req.Repo,
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

	go func() {
		// https://jsoverson.medium.com/how-to-deploy-node-js-functions-to-google-cloud-8bba05e9c10a
		cmd := exec.Command("gcloud", "--format", "json", "run", "deploy", service.Id, "--region", req.Region,
			"--cpu", "100m", "--memory", "128Mi", "--port", fmt.Sprintf("%d", req.Port), "--use-http2",
			"--allow-unauthenticated", "--max-instances", "1",
		)

		// if env vars exist then set them
		if len(envVars) > 0 {
			cmd.Args = append(cmd.Args, "--set-env-vars", strings.Join(envVars, ","))
		}

		// set the command dir
		cmd.Dir = gitter.RepoDir()

		// execute the command
		outp, err := cmd.CombinedOutput()
		if err != nil {
			log.Error(fmt.Errorf(string(outp)))
		}

		// check the app name reservation for custom domain mapping
		// apply the custom domain mapping if it exists and is valid
		var output map[string]interface{}

		// get the status output and deployment url
		if err := json.Unmarshal(outp, &output); err == nil {
			status := output["status"].(map[string]interface{})
			url := status["address"].(map[string]interface{})["url"].(string)
			deployed := status["conditions"].([]interface{})[0].(map[string]interface{})

			service.Url = url
			service.DeployedAt = deployed["lastTransitionTime"].(string)

			if deployed["status"] == "True" {
				service.Status = "Running"
			}
		} else {
			// TODO: return error
			service.Status = "Error"
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
	}()

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

	if err = json.Unmarshal(b, srv); err != nil {
		return err
	}

	// delete the app
	cmd := exec.Command("gcloud", "run", "services", "delete", "--region", srv.Region, srv.Id)
	outp, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(fmt.Errorf(string(outp)))
		return errors.InternalServerError("app.delete", "Failed to delete app")
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
	cmd := exec.Command("gcloud", "--format", "json", "run", "services", "describe", "--region", srv.Region, srv.Id)
	outp, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(fmt.Errorf(string(outp)))
		return fmt.Errorf("function does not exist")
	}

	output := map[string]interface{}{}
	if err = json.Unmarshal(outp, output); err != nil {
		return err
	}

	// get the service status
	status := output["status"].(map[string]interface{})
	deployed := status["conditions"].([]interface{})[0].(map[string]interface{})

	srv.DeployedAt = deployed["lastTransitionTime"].(string)

	if deployed["status"] == "True" {
		srv.Status = "Running"
	}

	// set response
	rsp.Service = srv

	return nil
}
