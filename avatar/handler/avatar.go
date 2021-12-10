package handler

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"

	"github.com/pkg/errors"

	"github.com/o1egl/govatar"

	pb "github.com/micro/services/avatar/proto"
)

type Avatar struct{}

// Generate is used to generate a avatar
func (e *Avatar) Generate(_ context.Context, req *pb.GenerateRequest, rsp *pb.GenerateResponse) error {
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
	var avatar image.Image
	var err error

	if req.Username == "" {
		avatar, err = govatar.Generate(gender)
	} else {
		avatar, err = govatar.GenerateForUsername(gender, req.Username)
	}
	if err != nil {
		return errors.Wrap(err, "generate avatar error")
	}

	// format avatar image, default is `jepg`
	format := req.Format
	if format != "png" && format != "jpeg" {
		format = "jpeg"
	}

	buf := bytes.NewBuffer(nil)
	if format == "png" {
		err = png.Encode(buf, avatar)
	} else {
		err = jpeg.Encode(buf, avatar, nil)
	}
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("encode %s format error", req.Format))
	}

	// response
	rsp.Format = req.Format
	rsp.Avatar = base64.StdEncoding.EncodeToString(buf.Bytes())

	return nil
}
