package handler

import (
	"context"
	"regexp"

	pauth "github.com/micro/services/pkg/auth"
	adminpb "github.com/micro/services/pkg/service/proto"
	"micro.dev/v4/service/errors"
	"micro.dev/v4/service/store"
)

var (
	IDFormat   = regexp.MustCompilePOSIX("^[a-z0-9-]+$")
	NameFormat = regexp.MustCompilePOSIX("^[a-z0-9]+$")

	FunctionKey    = "function/func/"
	OwnerKey       = "function/owner/"
	ReservationKey = "function/reservation/"
	BuildLogsKey   = "function/buildlogs/"
)

type Function struct{}

func (a *Function) Usage(ctx context.Context, request *adminpb.UsageRequest, response *adminpb.UsageResponse) error {
	method := "admin.Usage"
	_, err := pauth.VerifyMicroAdmin(ctx, method)
	if err != nil {
		return err
	}

	if len(request.TenantId) < 10 { // deliberate length check so we don't grab all the things
		return errors.BadRequest(method, "Missing tenant ID")
	}

	id := request.TenantId
	key := OwnerKey + id + "/"

	recs, err := store.List(store.ListPrefix(key))
	if err != nil {
		return err
	}

	funcs := len(recs)

	response.Usage = map[string]*adminpb.Usage{
		"Function.Deploy": &adminpb.Usage{Usage: int64(funcs), Units: "functions"},
	}

	return nil
}
