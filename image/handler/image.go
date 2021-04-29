package handler

import (
	"bytes"
	"context"
	"encoding/base64"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/micro/micro/v3/service/store"
	img "github.com/micro/services/image/proto"
)

const pathPrefix = "images/"

type Image struct{}

func (e *Image) Upload(ctx context.Context, req *img.UploadRequest, rsp *img.UploadResponse) error {
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
	if req.OutputURL {
		err = store.DefaultBlobStore.Write(pathPrefix+req.ImageID, buf)
		if err != nil {
			return err
		}
		return nil
	} else {
		dst := []byte{}
		base64.StdEncoding.Encode(dst, buf.Bytes())
		rsp.Base64 = string(dst)
		return nil
	}
	return nil
}

func base64ToImage(b64 string) (image.Image, error) {
	var srcImage image.Image
	res := []byte{}
	_, err := base64.StdEncoding.Decode([]byte(strings.Split(b64, ",")[1]), res)
	if err != nil {
		return srcImage, err
	}
	switch {
	case strings.HasPrefix("data:image/png", b64):
		srcImage, err = png.Decode(bytes.NewReader(res))
	case strings.HasPrefix("data:image/jpg", b64) || strings.HasPrefix("data:image/jpeg", b64):
		srcImage, err = jpeg.Decode(bytes.NewReader(res))
	}
	return srcImage, nil
}

func (e *Image) Resize(ctx context.Context, req *img.ResizeRequest, rsp *img.ResizeResponse) error {
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
	resultImage := imaging.Resize(srcImage, int(req.Width), int(req.Height), imaging.Lanczos)
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, resultImage, nil)
	if err != nil {
		return err
	}
	if req.OutputURL {
		err = store.DefaultBlobStore.Write(pathPrefix+req.ImageID, buf)
		if err != nil {
			return err
		}
		return nil
	} else {
		dst := []byte{}
		base64.StdEncoding.Encode(dst, buf.Bytes())
		rsp.Base64 = string(dst)
		return nil
	}
	return nil
}

func (e *Image) Convert(ctx context.Context, req *img.ConvertRequest, rsp *img.ConvertResponse) error {
	return nil
}
