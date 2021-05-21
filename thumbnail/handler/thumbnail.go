package handler

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/logger"
	iproto "github.com/micro/services/image/proto"
	thumbnail "github.com/micro/services/thumbnail/proto"
)

const screenshotPath = "/usr/src/app"

type Thumbnail struct {
	imageService iproto.ImageService
}

func NewThumbnail(imageService iproto.ImageService) *Thumbnail {
	return &Thumbnail{
		imageService: imageService,
	}
}

func (e *Thumbnail) Screenshot(ctx context.Context, req *thumbnail.ScreenshotRequest, rsp *thumbnail.ScreenshotResponse) error {
	imageName := uuid.New().String() + ".png"
	imagePath := filepath.Join(screenshotPath, imageName)
	width := "800"
	height := "600"
	if req.Width != 0 {
		width = fmt.Sprintf("%v", req.Width)
	}
	if req.Height != 0 {
		height = fmt.Sprintf("%v", req.Height)
	}

	outp, err := exec.Command("/usr/bin/chromium-browser", "--headless", "--window-size="+width+","+height, "--no-sandbox", "--screenshot="+imagePath, "--hide-scrollbars", req.Url).CombinedOutput()
	logger.Info(string(outp))
	if err != nil {
		logger.Error(string(outp) + err.Error())
		return err
	}
	file, err := ioutil.ReadFile(imagePath)
	if err != nil {
		return err
	}
	base := base64.StdEncoding.EncodeToString(file)
	resp, err := e.imageService.Upload(ctx, &iproto.UploadRequest{
		Base64:  "data:image/png;base64, " + base,
		ImageID: imageName,
	}, client.WithDialTimeout(20*time.Second), client.WithRequestTimeout(20*time.Second))
	if err != nil {
		return err
	}
	rsp.ImageURL = resp.Url
	return nil
}
