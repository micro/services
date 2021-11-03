package handler

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
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
	pid := 0
	defer func() {
		if err := os.Remove(imagePath); err != nil {
			logger.Errorf("Error removing file %s", err)
		}
		if pid != 0 {
			// using -ve PID kills the process group
			if err := syscall.Kill(-pid, syscall.SIGTERM); err != nil {
				logger.Errorf("Error killing process %s", err)
			}
		}
	}()
	width := "800"
	height := "600"
	if req.Width != 0 {
		width = fmt.Sprintf("%v", req.Width)
	}
	if req.Height != 0 {
		height = fmt.Sprintf("%v", req.Height)
	}
	cmd := exec.Command("/usr/bin/chromium-browser", "--headless", "--window-size="+width+","+height, "--no-sandbox", "--screenshot="+imagePath, "--hide-scrollbars", req.Url)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b
	if err := cmd.Start(); err != nil {
		logger.Errorf("Error starting %s", err)
		return err
	}
	pid = cmd.Process.Pid
	err := cmd.Wait()
	logger.Info(b.String())
	if err != nil {
		logger.Error(string(b.String()) + err.Error())
		return err
	}
	file, err := ioutil.ReadFile(imagePath)
	if err != nil {
		return err
	}
	base := base64.StdEncoding.EncodeToString(file)
	resp, err := e.imageService.Upload(ctx, &iproto.UploadRequest{
		Base64: "data:image/png;base64, " + base,
		Name:   imageName,
	}, client.WithDialTimeout(20*time.Second), client.WithRequestTimeout(20*time.Second))
	if err != nil {
		return err
	}
	rsp.ImageURL = resp.Url
	return nil
}
