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
	"net/url"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/micro/micro/v3/service/config"
	merrors "github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/store"
	img "github.com/micro/services/image/proto"
	"github.com/micro/services/pkg/tenant"
)

const pathPrefix = "images"
const hostPrefix = "https://cdn.m3ocontent.com"

type Image struct {
	hostPrefix string
}

func NewImage() *Image {
	var hp string
	cfg, err := config.Get("micro.image.host_prefix")
	if err != nil {
		hp = cfg.String(hostPrefix)
	}
	if len(strings.TrimSpace(hp)) == 0 {
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
	var ext string

	if len(req.Base64) > 0 {
		srcImage, ext, err = base64ToImage(req.Base64)
		if err != nil {
			return err
		}
	} else if len(req.Url) > 0 {
		ur, err := url.Parse(req.Url)
		if err != nil {
			return err
		}
		response, err := http.Get(req.Url)
		if err != nil {
			return err
		}
		defer response.Body.Close()
		switch {
		case strings.HasSuffix(ur.Path, ".png"):
			srcImage, err = png.Decode(response.Body)
		case strings.HasSuffix(ur.Path, ".jpg") || strings.HasSuffix(ur.Path, ".jpeg"):
			srcImage, err = jpeg.Decode(response.Body)
		}
		if err != nil {
			return err
		}

	} else {
		return errors.New("base64 or url param is required")
	}

	buf := new(bytes.Buffer)

	switch {
	case strings.HasSuffix(req.Name, ".png") || ext == "png":
		err = png.Encode(buf, srcImage)
	case strings.HasSuffix(req.Name, ".jpg") || strings.HasSuffix(req.Url, ".jpeg") || ext == "jpg":
		err = jpeg.Encode(buf, srcImage, nil)
	default:
		return errors.New("could not determine extension")
	}

	if err != nil {
		return err
	}

	err = store.DefaultBlobStore.Write(fmt.Sprintf("%v/%v/%v", pathPrefix, tenantID, req.Name), buf, store.BlobPublic(true))
	if err != nil {
		return err
	}
	rsp.Url = fmt.Sprintf("%v/%v/%v/%v/%v", e.hostPrefix, "micro", "images", tenantID, req.Name)
	return nil
}

func base64ToImage(b64 string) (image.Image, string, error) {
	var srcImage image.Image
	ext := ""

	parts := strings.Split(b64, ",")
	if len(parts) != 2 {
		return srcImage, "", merrors.BadRequest("image", "Incorrect format for base64 image, expected <encoding prefix>,<image data>")
	}
	prefix := parts[0]
	b64 = strings.TrimSpace(parts[1])
	res, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return srcImage, ext, err
	}

	switch prefix {
	case "data:image/png;base64":
		srcImage, err = png.Decode(bytes.NewReader(res))
		ext = "png"
	case "data:image/jpg;base64", "data:image/jpeg;base64":
		srcImage, err = jpeg.Decode(bytes.NewReader(res))
		ext = "jpg"
	default:
		return srcImage, ext, errors.New("unrecognized base64 prefix: " + prefix)
	}
	return srcImage, ext, err
}

func (e *Image) Resize(ctx context.Context, req *img.ResizeRequest, rsp *img.ResizeResponse) error {
	tenantID, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.New("Not authorized")
	}
	var srcImage image.Image
	var err error
	var ext string

	if len(req.Base64) > 0 {
		srcImage, ext, err = base64ToImage(req.Base64)
		if err != nil {
			return err
		}
	} else if len(req.Url) > 0 {
		ur, err := url.Parse(req.Url)
		if err != nil {
			return err
		}
		response, err := http.Get(req.Url)
		if err != nil {
			return err
		}
		defer response.Body.Close()
		switch {
		case strings.HasSuffix(ur.Path, ".png"):
			srcImage, err = png.Decode(response.Body)
		case strings.HasSuffix(ur.Path, ".jpg") || strings.HasSuffix(ur.Path, ".jpeg"):
			srcImage, err = jpeg.Decode(response.Body)
		}
		if err != nil {
			return err
		}

	} else {
		return errors.New("base64 or url param is required")
	}

	resultImage := imaging.Resize(srcImage, int(req.Width), int(req.Height), imaging.Lanczos)
	if req.CropOptions != nil {
		anchor := imaging.Center
		switch req.CropOptions.Anchor {
		case "top left":
			anchor = imaging.TopLeft
		case "top":
			anchor = imaging.Top
		case "top right":
			anchor = imaging.TopRight
		case "left":
			anchor = imaging.Left
		case "bottom left":
			anchor = imaging.BottomLeft
		case "bottom":
			anchor = imaging.Bottom
		case "bottom right":
			anchor = imaging.BottomRight
		}
		resultImage = imaging.CropAnchor(resultImage, int(req.CropOptions.Width), int(req.CropOptions.Height),
			anchor)
	}
	buf := new(bytes.Buffer)

	switch {
	case strings.HasSuffix(req.Name, ".png") || ext == "png":
		err = png.Encode(buf, resultImage)
	case strings.HasSuffix(req.Name, ".jpg") || strings.HasSuffix(req.Url, ".jpeg") || ext == "jpg":
		err = jpeg.Encode(buf, resultImage, nil)
	default:
		return errors.New("could not determine extension")
	}

	if err != nil {
		return err
	}
	if req.OutputURL {
		err = store.DefaultBlobStore.Write(fmt.Sprintf("%v/%v/%v", pathPrefix, tenantID, req.Name), buf, store.BlobPublic(true))
		if err != nil {
			return err
		}
		rsp.Url = fmt.Sprintf("%v/%v/%v/%v/%v", e.hostPrefix, "micro", "images", tenantID, req.Name)
	} else {
		prefix := "data:image/png;base64, "
		if ext == "jpg" {
			prefix = "data:image/jpg;base64, "
		}
		rsp.Base64 = prefix + base64.StdEncoding.EncodeToString(buf.Bytes())
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
		srcImage, _, err = base64ToImage(req.Base64)
		if err != nil {
			return err
		}
	} else {
		ur, err := url.Parse(req.Url)
		if err != nil {
			return err
		}

		response, err := http.Get(req.Url)
		if err != nil {
			return err
		}
		defer response.Body.Close()
		switch {
		case strings.HasSuffix(ur.Path, ".png"):
			srcImage, err = png.Decode(response.Body)
		case strings.HasSuffix(ur.Path, ".jpg") || strings.HasSuffix(ur.Path, ".jpeg"):
			srcImage, err = jpeg.Decode(response.Body)
		}
		if err != nil {
			return err
		}

	}

	buf := new(bytes.Buffer)
	switch {
	case strings.HasSuffix(req.Name, ".png"):
		err = png.Encode(buf, srcImage)
	case strings.HasSuffix(req.Name, ".jpg") || strings.HasSuffix(req.Url, ".jpeg"):
		err = jpeg.Encode(buf, srcImage, nil)
	}

	if err != nil {
		return err
	}
	if req.OutputURL {
		err = store.DefaultBlobStore.Write(fmt.Sprintf("%v/%v/%v", pathPrefix, tenantID, req.Name), buf, store.BlobPublic(true))
		if err != nil {
			return err
		}
		rsp.Url = fmt.Sprintf("%v/%v/%v/%v/%v", e.hostPrefix, "micro", "images", tenantID, req.Name)
	} else {
		src := buf.Bytes()
		length := base64.StdEncoding.EncodedLen(len(src))
		dst := make([]byte, length)
		base64.StdEncoding.Encode(dst, src)

		rsp.Base64 = string(dst)
		return nil
	}
	return nil
}
