package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/micro/micro/v3/proto/api"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	github "github.com/micro/services/github/proto"
	"github.com/micro/services/pkg/tenant"
)

type conf struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	PEM          string `json:"pem"`
}

func NewHandler(srv *service.Service) *Github {
	c, err := config.Get("micro.github")
	if err != nil {
		log.Fatalf("Error loading config %s", err)
	}
	var cfg conf
	if err := c.Scan(&cfg); err != nil {
		log.Fatalf("Error loading config %s", err)
	}
	if len(cfg.ClientSecret) == 0 || len(cfg.ClientID) == 0 || len(cfg.PEM) == 0 {
		log.Fatalf("Missing github config %+v", cfg)
	}
	return &Github{
		conf: cfg,
	}
}

type Github struct {
	conf conf
}

/*
1. Ask user to install app; send them to https://github.com/apps/m3o-apps-dev/installations/new?state=foobar
2. On success GH callback to m3o.com/<SOMETHING/>?code=0f957374a642fe783c0c&installation_id=24302523&setup_action=install&state=foobar (GH -> FRONTEND)
3. call https://api.m3o.com/github/authorize with the code and installation ID so we can store the installation ID and exchange code for a GH user token (FRONTEND -> BACKEND)
3. exchange the code for an access token; https://github.com/login/oauth/access_token -d '{"client_id": "<client ID>", "client_secret":"<secret>>", "code":"0f957374a642fe783c0c"}' -H "Content-Type: application/json"
response: access_token=<TOKEN>&scope=&token_type=bearer (BACKEND -> GH)
5. how can you be sure that an m3o user has access to an installation?? Probably pass the user access token and validate the user against the webhook passed in user
6. On success load up orgs repo list and display to user to pick (probably proxy query through back end so we don't store the token on the frontend)
7. UI will load list of projects via m3o.com rather than github directly
8. once chosen, backend can pull the contents of a repo at will by authenticating as the installation https://docs.github.com/en/developers/apps/building-github-apps/authenticating-with-github-apps#authenticating-as-an-installation

how can you be sure that an m3o user has access to an installation??

How does second app work, i.e. already installed the GH app, how do i select and run a second GH repo? Store the user access token or something? What happens if that gets blown away? Do we need to redo something??
*/

const (
	urlOauthAccessToken = "https://github.com/login/oauth/access_token"

	prefixInstallByTenant = "installByTenant"
)

type installation struct {
	TenantID       string
	InstallationID string
}

func (e *Github) Authorize(ctx context.Context, req *github.AuthorizeRequest, rsp *github.AuthorizeResponse) error {
	method := "github.Authorize"
	tenantID, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.Unauthorized(method, "Unauthorized")
	}
	if len(req.Code) == 0 {
		return errors.BadRequest(method, "Missing code param")
	}
	if len(req.InstallationId) == 0 {
		return errors.BadRequest(method, "Missing installation_id param")
	}
	b, _ := json.Marshal(map[string]string{
		"client_id":     e.conf.ClientID,
		"client_secret": e.conf.ClientSecret,
		"code":          req.Code,
	})
	oauthReq, err := http.NewRequest("POST", urlOauthAccessToken, bytes.NewReader(b))
	if err != nil {
		logger.Errorf("Error authorising code with Github %s", err)
		return errors.InternalServerError(method, "Error authorising code with Github")
	}
	oauthReq.Header.Set("Content-Type", "application/json")
	oauthReq.Header.Set("Accept", "application/json")
	oauthRsp, err := http.DefaultClient.Do(oauthReq)
	if err != nil {
		logger.Errorf("Error authorising code with Github %s", err)
		return errors.InternalServerError(method, "Error authorising code with Github")
	}
	defer oauthRsp.Body.Close()
	b, err := ioutil.ReadAll(oauthRsp.Body)
	if err != nil || oauthRsp.StatusCode > 399 {
		logger.Errorf("Error authorising code with Github %s %s", oauthRsp.Status, err)
		return errors.InternalServerError(method, "Error authorising code with Github")
	}
	var oauthRspMap map[string]string
	if err := json.Unmarshal(b, &oauthRspMap); err != nil {
		logger.Errorf("Error authorising code with Github %s %s", oauthRsp.Status, err)
		return errors.InternalServerError(method, "Error authorising code with Github")
	}
	acccessTok := oauthRspMap["access_token"]
	rsp.Token = acccessTok
	if err := store.Write(store.NewRecord(installationKey(tenantID, req.InstallationId), installation{
		InstallationID: req.InstallationId,
		TenantID:       tenantID})); err != nil {
		logger.Errorf("Error storing installation ID %s", err)
		return errors.InternalServerError(method, "Error authorising code with Github")
	}
	return nil
}

func installationKey(tenantID, installationID string) string {
	return fmt.Sprintf("%s/%s/%s", prefixInstallByTenant, tenantID, installationID)
}

func (e *Github) Webhook(ctx context.Context, req *api.Request, rsp *api.Response) error {
	//md, ok := metadata.FromContext(ctx)
	//if !ok {
	//	log.Errorf("Missing metadata from request")
	//	return errors.BadRequest("github.Webhook", "Missing headers")
	//}
	logger.Infof("Received webhook %v", req.Body)
	return nil
}
