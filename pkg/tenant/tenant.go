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
	id := acc.ID
	issuer := acc.Issuer
	owner := acc.Metadata["apikey_owner"]

	if len(id) == 0 {
		id = "micro"
	}

	if len(issuer) == 0 {
		issuer = "micro"
	}

	if len(owner) == 0 {
		owner = id
	}

	return fmt.Sprintf("%s/%s", acc.Issuer, owner)
}

// CreateKey generated a unique key for a resource
func CreateKey(ctx context.Context, key string) string {
	t, ok := FromContext(ctx)
	if !ok {
		return key
	}
	// return a tenant prefixed key e.g micro/asim/foobar
	return fmt.Sprintf("%s/%s", t, key)
}
