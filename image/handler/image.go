package handler

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	img "github.com/micro/services/image/proto"
	"github.com/micro/services/pkg/tenant"
)

const pathPrefix = "images"
const hostPrefix = "https://micro-store-bucket-125b9f0.ams3.cdn.digitaloceanspaces.com"

type Image struct {
	hostPrefix string
}

func NewImage() *Image {
	var hp string
	cfg, err := config.Get("micro.image.host_prefix")
	if err != nil {
		hp = cfg.String(hostPrefix)
	}
	if len(hp) == 0 {
		hp = hostPrefix
	}
	return &Image{
		hostPrefix: hp,
	}
}

func (e *Image) Upload(ctx context.Context, req *img.UploadRequest, rsp *img.UploadResponse) error {
	tenantID, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.New("Not authorized")
	}
	var srcImage image.Image
	var err error
	if len(req.Base64) > 0 {
		srcImage, err = base64ToImage(req.Base64)
		if err != nil {
			return err
		}
	} else {
		response, err := http.Get(req.Url)
		if err != nil {
			return err
		}
		switch {
		case strings.HasSuffix(req.Url, ".png"):
			srcImage, err = png.Decode(response.Body)
		case strings.HasSuffix(req.Url, ".jpg") || strings.HasSuffix(req.Url, ".jpeg"):
			srcImage, err = jpeg.Decode(response.Body)
		}
		if err != nil {
			return err
		}
		defer response.Body.Close()
	}
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, srcImage, nil)
	if err != nil {
		return err
	}

	err = store.DefaultBlobStore.Write(fmt.Sprintf("%v/%v/%v", pathPrefix, tenantID, req.ImageID), buf)
	if err != nil {
		return err
	}
	rsp.Url = fmt.Sprintf("%v/%v/%v/%v/%v", e.hostPrefix, "micro", "images", tenantID, req.ImageID)
	return nil
}

func base64ToImage(b64 string) (image.Image, error) {
	var srcImage image.Image

	res := []byte{}
	parts := strings.Split(b64, ",")
	prefix := parts[0]
	b64 = strings.TrimSpace(parts[1])
	_, err := base64.StdEncoding.Decode([]byte(b64), res)
	if err != nil {
		return srcImage, err
	}

	switch prefix {
	case "data:image/png;base64":
		srcImage, err = png.Decode(bytes.NewReader(res))
	case "data:image/jpg;base64", "data:image/jpeg;base64":
		srcImage, err = jpeg.Decode(bytes.NewReader(res))
	default:
		return srcImage, errors.New("unrecognized base64 prefix: " + prefix)
	}
	return srcImage, nil
}

func (e *Image) Resize(ctx context.Context, req *img.ResizeRequest, rsp *img.ResizeResponse) error {
	tenantID, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.New("Not authorized")
	}
	var srcImage image.Image
	var err error
	if len(req.Base64) > 0 {
		srcImage, err = base64ToImage(req.Base64)
		if err != nil {
			return err
		}
	} else {
		response, err := http.Get(req.Url)
		if err != nil {
			return err
		}
		switch {
		case strings.HasSuffix(req.Url, ".png"):
			srcImage, err = png.Decode(response.Body)
		case strings.HasSuffix(req.Url, ".jpg") || strings.HasSuffix(req.Url, ".jpeg"):
			srcImage, err = jpeg.Decode(response.Body)
		}
		if err != nil {
			return err
		}
		defer response.Body.Close()
	}
	logger.Infof("yo %v", srcImage)
	resultImage := imaging.Resize(srcImage, int(req.Width), int(req.Height), imaging.Lanczos)
	buf := new(bytes.Buffer)
	switch {
	case strings.HasSuffix(req.ImageID, ".png"):
		err = png.Encode(buf, resultImage)
	case strings.HasSuffix(req.ImageID, ".jpg") || strings.HasSuffix(req.Url, ".jpeg"):
		err = jpeg.Encode(buf, resultImage, nil)
	}

	if err != nil {
		return err
	}
	if req.OutputURL {
		err = store.DefaultBlobStore.Write(fmt.Sprintf("%v/%v/%v", pathPrefix, tenantID, req.ImageID), buf)
		if err != nil {
			return err
		}
		rsp.Url = fmt.Sprintf("%v/%v/%v/%v/%v", e.hostPrefix, "micro", "images", tenantID, req.ImageID)
	} else {
		dst := []byte{}
		base64.StdEncoding.Encode(dst, buf.Bytes())
		rsp.Base64 = string(dst)
		return nil
	}
	return nil
}

func (e *Image) Convert(ctx context.Context, req *img.ConvertRequest, rsp *img.ConvertResponse) error {
	tenantID, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.New("Not authorized")
	}
	var srcImage image.Image
	var err error
	if len(req.Base64) > 0 {
		srcImage, err = base64ToImage(req.Base64)
		if err != nil {
			return err
		}
	} else {
		response, err := http.Get(req.Url)
		if err != nil {
			return err
		}
		switch {
		case strings.HasSuffix(req.Url, ".png"):
			srcImage, err = png.Decode(response.Body)
		case strings.HasSuffix(req.Url, ".jpg") || strings.HasSuffix(req.Url, ".jpeg"):
			srcImage, err = jpeg.Decode(response.Body)
		}
		if err != nil {
			return err
		}
		defer response.Body.Close()
	}

	buf := new(bytes.Buffer)
	switch {
	case strings.HasSuffix(req.ImageID, ".png"):
		err = png.Encode(buf, srcImage)
	case strings.HasSuffix(req.ImageID, ".jpg") || strings.HasSuffix(req.Url, ".jpeg"):
		err = jpeg.Encode(buf, srcImage, nil)
	}

	if err != nil {
		return err
	}
	if req.OutputURL {
		err = store.DefaultBlobStore.Write(fmt.Sprintf("%v/%v/%v", pathPrefix, tenantID, req.ImageID), buf)
		if err != nil {
			return err
		}
		rsp.Url = fmt.Sprintf("%v/%v/%v/%v/%v", e.hostPrefix, "micro", "images", tenantID, req.ImageID)
	} else {
		dst := []byte{}
		base64.StdEncoding.Encode(dst, buf.Bytes())
		rsp.Base64 = string(dst)
		return nil
	}
	return nil
}
