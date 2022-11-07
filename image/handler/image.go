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
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/micro/micro/v3/service/config"
	merrors "github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	img "github.com/micro/services/image/proto"
	pauth "github.com/micro/services/pkg/auth"
	adminpb "github.com/micro/services/pkg/service/proto"
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
		return merrors.Unauthorized("image.Upload", "Not authorized")
	}
	var imageBytes []byte
	var err error

	if len(req.File) > 0 {
		imageBytes = req.File
	} else if len(req.Base64) > 0 {
		b, _, err := base64ToImage(req.Base64)
		if err != nil {
			return err
		}
		imageBytes = b
	} else if len(req.Url) > 0 {
		_, err := url.Parse(req.Url)
		if err != nil {
			return err
		}
		response, err := http.Get(req.Url)
		if err != nil {
			return err
		}
		defer response.Body.Close()
		b, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		imageBytes = b
	} else {
		return merrors.BadRequest("image.Upload", "file, base64 or url param is required")
	}

	// validate that this is indeed an image file
	_, _, err = image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		if err == image.ErrFormat {
			return merrors.BadRequest("image.Upload", "Unrecognised image format")
		}
		return merrors.InternalServerError("image.Upload", "Error processing upload")
	}

	err = store.DefaultBlobStore.Write(fmt.Sprintf("%v/%v/%v", pathPrefix, tenantID, req.Name), bytes.NewReader(imageBytes), store.BlobPublic(true))
	if err != nil {
		return err
	}
	rsp.Url = fmt.Sprintf("%v/%v/%v/%v/%v", e.hostPrefix, "micro", "images", tenantID, req.Name)
	return nil
}

func base64ToImage(b64 string) ([]byte, string, error) {
	ext := ""
	parts := strings.Split(b64, ",")
	if len(parts) != 2 {
		return nil, "", merrors.BadRequest("image", "Incorrect format for base64 image, expected <encoding prefix>,<image data>")
	}
	prefix := parts[0]
	b64 = strings.TrimSpace(parts[1])

	switch prefix {
	case "data:image/png;base64":
		ext = "png"
	case "data:image/jpg;base64", "data:image/jpeg;base64":
		ext = "jpg"
	default:
		return nil, ext, errors.New("unrecognized base64 prefix: " + prefix)
	}
	b, err := base64.StdEncoding.DecodeString(b64)
	return b, ext, err
}

func (e *Image) Resize(ctx context.Context, req *img.ResizeRequest, rsp *img.ResizeResponse) error {
	tenantID, ok := tenant.FromContext(ctx)
	if !ok {
		return merrors.Unauthorized("image.Resize", "Not authorized")
	}
	var imageBytes []byte
	var err error
	var ext string

	if len(req.File) > 0 {
		imageBytes = req.File
	} else if len(req.Base64) > 0 {
		imageBytes, ext, err = base64ToImage(req.Base64)
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
		imageBytes, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		switch {
		case strings.HasSuffix(ur.Path, ".png"):
			ext = "png"
		case strings.HasSuffix(ur.Path, ".jpg") || strings.HasSuffix(ur.Path, ".jpeg"):
			ext = "jpg"
		}
	} else {
		return errors.New("base64 or url param is required")
	}

	srcImage, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return err
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
		return merrors.Unauthorized("image.Convert", "Not authorized")
	}
	var srcImage image.Image
	var imageBytes []byte
	var err error
	var ext string

	if len(req.File) > 0 {
		imageBytes = req.File
	} else if len(req.Base64) > 0 {
		imageBytes, _, err = base64ToImage(req.Base64)
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
		imageBytes, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		switch {
		case strings.HasSuffix(ur.Path, ".png"):
			ext = "png"
		case strings.HasSuffix(ur.Path, ".jpg") || strings.HasSuffix(ur.Path, ".jpeg"):
			ext = "jpg"
		}
	} else {
		return merrors.BadRequest("image.Convert", "Must pass either base64, url, or file param")
	}

	srcImage, _, err = image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	switch {
	case strings.HasSuffix(req.Name, ".png"):
		err = png.Encode(buf, srcImage)
	case strings.HasSuffix(req.Name, ".jpg") || strings.HasSuffix(req.Url, ".jpeg"):
		err = jpeg.Encode(buf, srcImage, nil)
	case strings.HasSuffix(req.Url, ".png"):
		err = png.Encode(buf, srcImage)
	case strings.HasSuffix(req.Url, ".jpg") || strings.HasSuffix(req.Url, ".jpeg"):
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
		prefix := "data:image/png;base64, "
		if ext == "jpg" {
			prefix = "data:image/jpg;base64, "
		}
		rsp.Base64 = prefix + base64.StdEncoding.EncodeToString(buf.Bytes())
		return nil
	}
	return nil
}

func (e *Image) Delete(ctx context.Context, request *img.DeleteRequest, response *img.DeleteResponse) error {
	tenantID, ok := tenant.FromContext(ctx)
	if !ok {
		return merrors.Unauthorized("image.Delete", "Not authorized")
	}
	if len(request.Url) == 0 {
		return merrors.BadRequest("image.Delete", "Missing URL parameter")
	}
	// parse the url <hostprefix>/micro/images/<tenantid>/<name>
	match, err := regexp.MatchString(fmt.Sprintf("^%s/micro/images/%s/.*", e.hostPrefix, tenantID), request.Url)
	if err != nil {
		logger.Errorf("Error matching req url %s", err)
		return merrors.InternalServerError("image.Delete", "Error processing delete")
	}
	if !match {
		logger.Infof("No match %s", request.Url, tenantID)
		return merrors.BadRequest("image.Delete", "URL not recognised for user")
	}
	tenantAndName := strings.TrimPrefix(request.Url, fmt.Sprintf("%s/micro/images/", e.hostPrefix))
	blobKey := fmt.Sprintf("%s/%s", pathPrefix, tenantAndName)
	if err := store.DefaultBlobStore.Delete(blobKey); err != nil {
		logger.Errorf("Error deleting key %s", err)
		return merrors.InternalServerError("image.Delete", "Error processing delete")
	}
	return nil
}

func (e *Image) DeleteData(ctx context.Context, request *adminpb.DeleteDataRequest, response *adminpb.DeleteDataResponse) error {
	method := "admin.DeleteData"
	_, err := pauth.VerifyMicroAdmin(ctx, method)
	if err != nil {
		return err
	}

	if len(request.TenantId) < 10 { // deliberate length check so we don't delete all the things
		return merrors.BadRequest(method, "Missing tenant ID")
	}

	path := fmt.Sprintf("%v/%v", pathPrefix, request.TenantId)
	keys, err := store.DefaultBlobStore.List(store.BlobListPrefix(path))
	if err != nil {
		return err
	}

	for _, key := range keys {
		err = store.DefaultBlobStore.Delete(key)
		if err != nil {
			return err
		}
	}

	logger.Infof("Deleted %d keys for %s", len(keys), request.TenantId)

	return nil
}

func (i *Image) Usage(ctx context.Context, request *adminpb.UsageRequest, response *adminpb.UsageResponse) error {
	method := "admin.Usage"
	_, err := pauth.VerifyMicroAdmin(ctx, method)
	if err != nil {
		return err
	}

	if len(request.TenantId) < 10 { // deliberate length check so we don't grab all the things
		return merrors.BadRequest(method, "Missing tenant ID")
	}

	key := fmt.Sprintf("%v/%v/", pathPrefix, request.TenantId)

	// list all images for the user
	recs, err := store.DefaultBlobStore.List(store.BlobListPrefix(key))
	if err != nil {
		return err
	}

	response.Usage = map[string]*adminpb.Usage{
		"Image.Upload": &adminpb.Usage{Usage: int64(len(recs)), Units: "images"},
		// all other methods don't add to space so are not usage capped
	}

	return nil
}
