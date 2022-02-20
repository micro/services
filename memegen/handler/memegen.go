package handler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/memegen/proto"
	"github.com/micro/services/pkg/api"
)

type Memegen struct {
	username string
	password string
}

func New() *Memegen {
	v, err := config.Get("imgflip.username")
	if err != nil {
		logger.Fatalf("imgflip.username config not found: %v", err)
	}
	username := v.String("")
	if len(username) == 0 {
		logger.Fatal("imgflip.username config not found")
	}
	v, err = config.Get("imgflip.password")
	if err != nil {
		logger.Fatalf("imgflip.password config not found: %v", err)
	}
	password := v.String("")
	if len(password) == 0 {
		logger.Fatal("imgflip.password config not found")
	}
	return &Memegen{
		username: username,
		password: password,
	}
}

type TemplateRequest struct {
	Success bool `json:"success"`
	Data    Data `json":data"`
}

type Data struct {
	Memes []*pb.Template `json:"memes"`
}

func (m *Memegen) Templates(ctx context.Context, req *pb.TemplatesRequest, rsp *pb.TemplatesResponse) error {
	templateRsp := new(TemplateRequest)
	if err := api.Get("https://api.imgflip.com/get_memes", templateRsp); err != nil {
		return errors.InternalServerError("memegen.templates", err.Error())
	}
	rsp.Templates = templateRsp.Data.Memes
	return nil
}

func (m *Memegen) Generate(ctx context.Context, req *pb.GenerateRequest, rsp *pb.GenerateResponse) error {
	vals := url.Values{}
	vals.Set("template_id", req.Id)
	vals.Set("text0", req.TopText)
	vals.Set("text1", req.BottomText)
	vals.Set("font", req.Font)
	vals.Set("max_font_size", req.MaxFontSize)
	vals.Set("username", m.username)
	vals.Set("password", m.password)

	genRsp := map[string]interface{}{}

	resp, err := http.PostForm("https://api.imgflip.com/caption_image", vals)
	if err != nil {
		return errors.InternalServerError("memegen.generate", err.Error())
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(b, &genRsp)

	success := genRsp["success"].(bool)
	if !success {
		return errors.BadRequest("memegen.generate", genRsp["error_message"].(string))
	}

	// set response url
	rsp.Url = genRsp["data"].(map[string]interface{})["url"].(string)
	return nil
}
