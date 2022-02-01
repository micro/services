package auth

import (
	"context"

	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
)

func VerifyMicroAdmin(ctx context.Context, method string) (*auth.Account, error) {
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return nil, errors.Unauthorized(method, "Unauthorized")
	}
	if err := doVerifyMicroAdmin(acc, method); err != nil {
		return nil, err
	}
	return acc, nil
}

func doVerifyMicroAdmin(acc *auth.Account, method string) error {
	errForbid := errors.Forbidden(method, "Forbidden")
	if acc.Issuer != "micro" {
		return errForbid
	}

	for _, s := range acc.Scopes {
		if (s == "admin" && acc.Type == "user") || (s == "service" && acc.Type == "service") {
			return nil
		}
	}
	return errForbid

}
