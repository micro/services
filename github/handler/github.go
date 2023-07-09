package handler

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	github "github.com/micro/services/github/proto"
	"github.com/micro/services/pkg/auth"
	pauth "github.com/micro/services/pkg/auth"
	adminpb "github.com/micro/services/pkg/service/proto"
	"github.com/micro/services/pkg/tenant"
	"micro.dev/v4/proto/api"
	"micro.dev/v4/service"
	"micro.dev/v4/service/config"
	"micro.dev/v4/service/errors"
	"micro.dev/v4/service/logger"
	"micro.dev/v4/service/store"
)

type conf struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	PEM          string `json:"pem"`
	AppID        string `json:"app_id"`
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

How does second app work, i.e. already installed the GH app, how do i select and run a second GH repo? Store the user access token or something? What happens if that gets blown away? Do we need to redo something??
*/

const (
	urlOauthAccessToken   = "https://github.com/login/oauth/access_token"
	urlInstallationTokens = "https://api.github.com/app/installations/%s/access_tokens"
	urlRepoList           = "https://api.github.com/installation/repositories"
	urlBranchesList       = "https://api.github.com/repos/%s/branches"

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
	b, err = ioutil.ReadAll(oauthRsp.Body)
	if err != nil || oauthRsp.StatusCode > 399 {
		logger.Errorf("Error authorising code with Github %s %s", oauthRsp.Status, err)
		return errors.InternalServerError(method, "Error authorising code with Github")
	}
	var oauthRspMap map[string]string
	if err := json.Unmarshal(b, &oauthRspMap); err != nil {
		logger.Errorf("Error authorising code with Github %s %s", oauthRsp.Status, err)
		return errors.InternalServerError(method, "Error authorising code with Github")
	}
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

func (e *Github) ListRepos(ctx context.Context, req *github.ListReposRequest, rsp *github.ListReposResponse) error {
	method := "github.ListRepos"
	tenantID, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.Unauthorized(method, "Unauthorized")
	}
	recs, err := store.Read(installationKey(tenantID, ""), store.ReadPrefix()) // TODO support multiple installations
	if err != nil {
		if err == store.ErrNotFound {
			return errors.NotFound(method, "Github app not installed for tenant")
		}
		logger.Errorf("Error retrieving token %s", err)
		return errors.InternalServerError(method, "Error retrieving repos")
	}
	var install installation
	if err := recs[0].Decode(&install); err != nil { // TODO support multiple installations
		logger.Errorf("Error retrieving token %s", err)
		return errors.InternalServerError(method, "Error retrieving repos")
	}

	tok, err := e.getInstallationToken(ctx, tenantID, install.InstallationID)
	if err != nil {
		logger.Errorf("Error retrieving token %s", err)
		return errors.InternalServerError(method, "Error retrieving repos")
	}
	repoReq, err := http.NewRequest("GET", urlRepoList, nil)
	if err != nil {
		logger.Errorf("Error retrieving repos %s", err)
		return errors.InternalServerError(method, "Error retrieving repos")
	}
	repoReq.Header.Set("Authorization", "Bearer "+tok)
	repoReq.Header.Set("Accept", "application/vnd.github.v3+json")
	repoRsp, err := http.DefaultClient.Do(repoReq)
	if err != nil {
		logger.Errorf("Error retrieving repos %s", err)
		return errors.InternalServerError(method, "Error retrieving repos")
	}
	defer repoRsp.Body.Close()
	b, _ := ioutil.ReadAll(repoRsp.Body)

	if repoRsp.StatusCode > 399 {
		logger.Errorf("Error retrieving repos %s %s", repoRsp.Status, string(b))
		return errors.InternalServerError(method, "Error retrieving repos")
	}
	rspObj := struct {
		Repositories []struct {
			FullName string `json:"full_name"`
		} `json:"repositories"`
	}{}
	if err := json.Unmarshal(b, &rspObj); err != nil {
		logger.Errorf("Error retrieving repos %s", err)
		return errors.InternalServerError(method, "Error retrieving repos")
	}
	rsp.Repos = make([]string, len(rspObj.Repositories))
	for i, v := range rspObj.Repositories {
		rsp.Repos[i] = v.FullName
	}
	return nil
}

func (e *Github) Token(ctx context.Context, req *github.TokenRequest, rsp *github.TokenResponse) error {
	method := "github.Token"
	_, err := auth.VerifyMicroAdmin(ctx, method)
	if err != nil {
		return err
	}
	recs, err := store.Read(installationKey(req.TenantId, ""), store.ReadPrefix()) // TODO support multiple installations
	if err != nil {
		if err == store.ErrNotFound {
			return errors.NotFound(method, "No installation found for tenant")
		}
		logger.Errorf("Error retrieving token %s", err)
		return errors.InternalServerError(method, "Error retrieving token")
	}
	var install installation
	if err := recs[0].Decode(&install); err != nil { // TODO support multiple installations
		logger.Errorf("Error retrieving token %s", err)
		return errors.InternalServerError(method, "Error retrieving token")
	}

	rsp.Token, err = e.getInstallationToken(ctx, req.TenantId, install.InstallationID)
	if err != nil {
		logger.Errorf("Error retrieving token %s", err)
		return errors.InternalServerError(method, "Error retrieving token")
	}
	return nil
}

func (e *Github) getInstallationToken(ctx context.Context, tenantID, installationID string) (string, error) {
	// TODO support multiple installations per user
	// make JWT
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),  // 10 minute max
		IssuedAt:  time.Now().Add(-60 * time.Second).Unix(), // allow for clock drift
		Issuer:    e.conf.AppID,
	})
	p, _ := pem.Decode([]byte(e.conf.PEM))
	pk, err := x509.ParsePKCS1PrivateKey(p.Bytes)
	if err != nil {
		logger.Fatalf("Failed to parse PEM %s", err)
	}

	jwtString, err := tok.SignedString(pk)
	if err != nil {
		logger.Errorf("Failed to generate signed JWT %s", err)
		return "", err
	}
	tokReq, err := http.NewRequest("POST", fmt.Sprintf(urlInstallationTokens, installationID), nil)
	if err != nil {
		logger.Errorf("Failed to generate installation token %s", err)
		return "", err
	}
	tokReq.Header.Set("Authorization", "Bearer "+jwtString)
	tokReq.Header.Set("Accept", "application/vnd.github.v3+json")
	tokRsp, err := http.DefaultClient.Do(tokReq)
	if err != nil {
		logger.Errorf("Failed to generate installation token %s", err)
		return "", err
	}
	defer tokRsp.Body.Close()
	b, _ := ioutil.ReadAll(tokRsp.Body)

	if tokRsp.StatusCode > 399 {
		logger.Errorf("Failed to generate installation token %s %s", tokRsp.Status, string(b))
		return "", err
	}

	var rspMap map[string]interface{}
	if err := json.Unmarshal(b, &rspMap); err != nil {
		logger.Errorf("Failed to generate installation token %s 5s", tokRsp.Status, string(b))
		return "", err
	}
	return rspMap["token"].(string), nil
}

func (e *Github) ListBranches(ctx context.Context, req *github.ListBranchesRequest, rsp *github.ListBranchesResponse) error {
	method := "github.ListBranches"
	tenantID, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.Unauthorized(method, "Unauthorized")
	}
	recs, err := store.Read(installationKey(tenantID, ""), store.ReadPrefix()) // TODO support multiple installations
	if err != nil {
		if err == store.ErrNotFound {
			return errors.NotFound(method, "Github app not installed for tenant")
		}
		logger.Errorf("Error retrieving token %s", err)
		return errors.InternalServerError(method, "Error retrieving repos")
	}
	var install installation
	if err := recs[0].Decode(&install); err != nil { // TODO support multiple installations
		logger.Errorf("Error retrieving token %s", err)
		return errors.InternalServerError(method, "Error retrieving repos")
	}

	tok, err := e.getInstallationToken(ctx, tenantID, install.InstallationID)
	if err != nil {
		logger.Errorf("Error retrieving token %s", err)
		return errors.InternalServerError(method, "Error retrieving repos")
	}
	repoReq, err := http.NewRequest("GET", fmt.Sprintf(urlBranchesList, req.Repo), nil)
	if err != nil {
		logger.Errorf("Error retrieving repos %s", err)
		return errors.InternalServerError(method, "Error retrieving repos")
	}
	repoReq.Header.Set("Authorization", "Bearer "+tok)
	repoReq.Header.Set("Accept", "application/vnd.github.v3+json")
	repoRsp, err := http.DefaultClient.Do(repoReq)
	if err != nil {
		logger.Errorf("Error retrieving repos %s", err)
		return errors.InternalServerError(method, "Error retrieving repos")
	}
	defer repoRsp.Body.Close()
	b, _ := ioutil.ReadAll(repoRsp.Body)

	if repoRsp.StatusCode > 399 {
		logger.Errorf("Error retrieving repos %s %s", repoRsp.Status, string(b))
		return errors.InternalServerError(method, "Error retrieving repos")
	}
	var rspObj []struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(b, &rspObj); err != nil {
		logger.Errorf("Error retrieving repos %s", err)
		return errors.InternalServerError(method, "Error retrieving repos")
	}
	rsp.Branches = make([]string, len(rspObj))
	for i, v := range rspObj {
		rsp.Branches[i] = v.Name
	}
	return nil
}

func (e *Github) Webhook(ctx context.Context, req *api.Request, rsp *api.Response) error {
	// TODO
	logger.Infof("Received webhook %v", req.Body)
	return nil
}

func (e *Github) DeleteData(ctx context.Context, request *adminpb.DeleteDataRequest, response *adminpb.DeleteDataResponse) error {
	method := "admin.DeleteData"
	_, err := pauth.VerifyMicroAdmin(ctx, method)
	if err != nil {
		return err
	}

	if len(request.TenantId) < 10 { // deliberate length check so we don't delete all the things
		return errors.BadRequest(method, "Missing tenant ID")
	}
	keys, err := store.List(store.ListPrefix(installationKey(request.TenantId, "")))
	if err != nil {
		return err
	}

	for _, key := range keys {
		err = store.Delete(key)
		if err != nil {
			return err
		}
	}
	logger.Infof("Deleted %d objects from S3 for %s", len(keys), request.TenantId)

	return nil
}

func (e *Github) Usage(ctx context.Context, request *adminpb.UsageRequest, response *adminpb.UsageResponse) error {
	return nil
}
