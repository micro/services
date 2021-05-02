// Package tenant provides multi-tenancy helpers
package tenant

import (
	"context"
	"fmt"

	"github.com/micro/micro/v3/service/auth"
)

// FromContext returns a tenant from the context
func FromContext(ctx context.Context) (string, bool) {
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return "", false
	}
	return FromAccount(acc), true
}

// FromAccount returns a tenant from
func FromAccount(acc *auth.Account) string {
	owner := acc.Metadata["apikey_owner"]
	if len(owner) == 0 {
		owner = acc.ID
	}
	return fmt.Sprintf("%s/%s", acc.Issuer, owner)
}
