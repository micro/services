package handler

import (
	"bytes"
	"context"
	"encoding/base64"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"

	"github.com/disintegration/imaging"
	"github.com/micro/micro/v3/service/store"
	img "github.com/micro/services/image/proto"
)

type Image struct{}

func (e *Image) Upload(ctx context.Context, req *img.UploadRequest, rsp *img.UploadResponse) error {
	var srcImage image.Image
	if len(req.Base64) > 0 {
		res := []byte{}
		_, err := base64.StdEncoding.Decode([]byte(req.Base64), res)
		if err != nil {
			return err
		}
		srcImage, err = png.Decode(bytes.NewReader(res))
		if err != nil {
			return err
		}
	} else {
		response, err := http.Get(req.Url)
		if err != nil {
			return err
		}
		srcImage, err = jpeg.Decode(response.Body)
		if err != nil {
			return err
		}
		defer response.Body.Close()
	}
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, srcImage, nil)
	if err != nil {
		return err
	}
	if req.OutputURL {
		err = store.DefaultBlobStore.Write(req.ImageID, buf)
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

func (e *Image) Resize(ctx context.Context, req *img.ResizeRequest, rsp *img.ResizeResponse) error {
	var srcImage image.Image
	if len(req.Base64) > 0 {
		res := []byte{}
		_, err := base64.StdEncoding.Decode([]byte(req.Base64), res)
		if err != nil {
			return err
		}
		srcImage, err = png.Decode(bytes.NewReader(res))
		if err != nil {
			return err
		}
	} else {
		response, err := http.Get(req.Url)
		if err != nil {
			return err
		}
		srcImage, err = jpeg.Decode(response.Body)
		if err != nil {
			return err
		}
		defer response.Body.Close()
	}
	resultImage := imaging.Resize(srcImage, int(req.Width), int(req.Height), imaging.Lanczos)
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, resultImage, nil)
	if err != nil {
		return err
	}
	if req.OutputURL {
		err = store.DefaultBlobStore.Write(req.ImageID, buf)
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
