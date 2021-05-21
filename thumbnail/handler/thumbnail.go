package handler

import (
	"context"
	"encoding/base64"
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

	outp, err := exec.Command("/usr/bin/chromium-browser", "--headless", "--no-sandbox", "--screenshot="+imagePath, "--hide-scrollbars", "https://www.chromestatus.com/").CombinedOutput()
	logger.Info(string(outp))
	if err != nil {
		logger.Error(string(outp) + err.Error())
		return err
	}
	file, err := ioutil.ReadFile(imagePath)
	if err != nil {
		return err
	}
	base := base64.RawStdEncoding.EncodeToString(file)
	resp, err := e.imageService.Upload(ctx, &iproto.UploadRequest{
		Base64:  base,
		ImageID: imageName,
	}, client.WithRequestTimeout(20*time.Second))
	if err != nil {
		return err
	}
	rsp.Url = resp.Url
	return nil
}
