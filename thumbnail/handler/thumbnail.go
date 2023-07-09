package handler

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	iproto "github.com/micro/services/image/proto"
	thumbnail "github.com/micro/services/thumbnail/proto"
	"micro.dev/v4/service/client"
	"micro.dev/v4/service/errors"
	"micro.dev/v4/service/logger"
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
	defer func() {
		os.Remove(imagePath)
	}()
	width := "800"
	height := "600"
	if req.Width != 0 {
		width = fmt.Sprintf("%v", req.Width)
	}
	if req.Height != 0 {
		height = fmt.Sprintf("%v", req.Height)
	}
	cmd := exec.Command("/usr/bin/chromium-browser",
		"--headless", "--window-size="+width+","+height, "--no-sandbox", "--screenshot="+imagePath,
		"--hide-scrollbars", "--disable-setuid-sandbox", "--single-process", "--no-zygote", "--disable-gpu", req.Url)
	outp, err := cmd.CombinedOutput()
	logger.Info(string(outp))
	if err != nil {
		logger.Error(string(outp) + err.Error())
		return errors.InternalServerError("thumbnail.Screenshot", "Error taking screenshot")
	}
	file, err := ioutil.ReadFile(imagePath)
	if err != nil {
		logger.Errorf("Error reading file %s", err)
		return errors.InternalServerError("thumbnail.Screenshot", "Error taking screenshot")
	}
	base := base64.StdEncoding.EncodeToString(file)
	resp, err := e.imageService.Upload(ctx, &iproto.UploadRequest{
		Base64: "data:image/png;base64, " + base,
		Name:   imageName,
	}, client.WithDialTimeout(20*time.Second), client.WithRequestTimeout(20*time.Second))
	if err != nil {
		logger.Errorf("Error uploading screenshot %s", err)
		return errors.InternalServerError("thumbnail.Screenshot", "Error taking screenshot")
	}
	rsp.ImageURL = resp.Url
	return nil
}
