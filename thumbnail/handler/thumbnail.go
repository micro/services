package handler

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/logger"
	iproto "github.com/micro/services/image/proto"
	thumbnail "github.com/micro/services/thumbnail/proto"
)

const screenshotPath = "/usr/src/app"

type Thumbnail struct {
	imageService iproto.ImageService
}

func NewThumbnail(imageService iproto.ImageService) *Thumbnail {
	return &Thumbnail{}
}

func (e *Thumbnail) Screenshot(ctx context.Context, req *thumbnail.ScreenshotRequest, rsp *thumbnail.ScreenshotResponse) error {
	id := uuid.New().String() + ".png"
	outp, err := exec.Command("/usr/bin/chromium-browser", "--headless", "--no-sandbox", "--screenshot="+id, "--hide-scrollbars", "https://www.chromestatus.com/").CombinedOutput()
	if err != nil {
		logger.Error(string(outp) + err.Error())
		return err
	}
	file, err := ioutil.ReadFile(filepath.Join(screenshotPath, id))
	if err != nil {
		return err
	}
	base := base64.RawStdEncoding.EncodeToString(file)
	resp, err := e.imageService.Upload(ctx, &iproto.UploadRequest{
		Base64:  base,
		ImageID: id,
	})
	if err != nil {
		return err
	}
	rsp.Url = resp.Url
	return nil
}
