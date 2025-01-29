package handler

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"

	"github.com/google/uuid"
	"github.com/micro/micro/v5/service/errors"
	"github.com/o1egl/govatar"

	pb "github.com/micro/services/avatar/proto"
	imagePb "github.com/micro/services/image/proto"
)

type avatar struct {
	imageSvc imagePb.ImageService
}

func NewAvatar(service imagePb.ImageService) *avatar {
	return &avatar{
		imageSvc: service,
	}
}

// Generate is used to generate a avatar
func (e *avatar) Generate(ctx context.Context, req *pb.GenerateRequest, rsp *pb.GenerateResponse) error {
	var gender govatar.Gender

	// gender, default is `male`
	if req.Gender == "male" {
		gender = govatar.MALE
	} else if req.Gender == "female" {
		gender = govatar.FEMALE
	} else {
		gender = govatar.MALE
	}

	// generate avatar
	var avatarImg image.Image
	var err error

	if req.Username == "" {
		avatarImg, err = govatar.Generate(gender)
	} else {
		avatarImg, err = govatar.GenerateForUsername(gender, req.Username)
	}
	if err != nil {
		return errors.InternalServerError("avatar.generate", "generate avatarImg error: %v", err)
	}

	// format avatar image, default is `jpeg`
	format := req.Format
	if format != "png" && format != "jpeg" {
		format = "jpeg"
	}

	buf := bytes.NewBuffer(nil)
	if format == "png" {
		err = png.Encode(buf, avatarImg)
	} else {
		err = jpeg.Encode(buf, avatarImg, nil)
	}
	if err != nil {
		return errors.InternalServerError("avatar.generate", "encode avatar image error: %v", err)
	}

	base64String := fmt.Sprintf("data:image/%s;base64,%s", format, base64.StdEncoding.EncodeToString(buf.Bytes()))

	if !req.Upload {
		rsp.Base64 = base64String
		return nil
	}

	// upload to CDN
	name := req.Username
	if name == "" {
		uid, _ := uuid.NewUUID()
		name = uid.String()
	}

	uploadResp, err := e.imageSvc.Upload(ctx, &imagePb.UploadRequest{
		Base64: base64String,
		Name:   fmt.Sprintf("%s.%s", name, format),
	})

	if err != nil {
		return errors.InternalServerError("avatar.generate", "upload avatar image error: %v", err)
	}

	rsp.Base64 = base64String
	rsp.Url = uploadResp.Url

	return nil
}
