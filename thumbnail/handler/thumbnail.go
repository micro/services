package handler

import (
	"context"
	"os/exec"

	"github.com/micro/micro/v3/service/logger"
	thumbnail "github.com/micro/services/thumbnail/proto"
)

type Thumbnail struct{}

func NewThumbnail() *Thumbnail {
	outp, err := exec.Command("/usr/bin/chromium-browser").CombinedOutput()
	logger.Info(string(outp) + err.Error())
	return &Thumbnail{}
}

func (e *Thumbnail) Screenshot(ctx context.Context, req *thumbnail.ScreenshotRequest, rsp *thumbnail.ScreenshotResponse) error {
	return nil
}
