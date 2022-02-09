package domain

import (
	"context"
	"fmt"
	"strings"

	"github.com/micro/services/pkg/tenant"
)

func getStoreKeyPrefix(ctx context.Context) string {
	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		tenantId = "micro"
	}

	return getStoreKeyPrefixForTenent(tenantId)
}

func getTenantKey(ctx context.Context) string {
	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		tenantId = "micro"
	}

	return strings.Replace(strings.Replace(tenantId, "/", "_", -1), "-", "_", -1)
}

func getStoreKeyPrefixForTenent(tenantId string) string {
	tid := strings.Replace(strings.Replace(tenantId, "/", "_", -1), "-", "_", -1)

	return fmt.Sprintf("user/%s/", tid)
}

func generateAccountTenantKey(tenantId, userId string) string {
	return fmt.Sprintf("%saccount/id/%s", getStoreKeyPrefixForTenent(tenantId), userId)
}
func generateAccountTenantEmailKey(tenantId, email string) string {
	return fmt.Sprintf("%sacccount/email/%s", getStoreKeyPrefixForTenent(tenantId), email)
}
func generateAccountTenantUsernameKey(tenantId, username string) string {
	return fmt.Sprintf("%saccount/username/%s", getStoreKeyPrefixForTenent(tenantId), username)
}

func generateAccountStoreKey(ctx context.Context, userId string) string {
	return fmt.Sprintf("%saccount/id/%s", getStoreKeyPrefix(ctx), userId)
}

func generateAccountEmailStoreKey(ctx context.Context, email string) string {
	return fmt.Sprintf("%sacccount/email/%s", getStoreKeyPrefix(ctx), email)
}

func generateAccountUsernameStoreKey(ctx context.Context, username string) string {
	return fmt.Sprintf("%saccount/username/%s", getStoreKeyPrefix(ctx), username)
}

func generatePasswordStoreKey(ctx context.Context, userId string) string {
	return fmt.Sprintf("%spassword/%s", getStoreKeyPrefix(ctx), userId)
}

func generatePasswordResetCodeStoreKey(ctx context.Context, userId, code string) string {
	return fmt.Sprintf("%spassword-reset-codes/%s-%s", getStoreKeyPrefix(ctx), userId, code)
}

func generateSessionStoreKey(ctx context.Context, sessionId string) string {
	return fmt.Sprintf("%ssession/%s", getStoreKeyPrefix(ctx), sessionId)
}

func generateVerificationTokenStoreKey(token string) string {
	return fmt.Sprintf("user/verification-token/%s", token)
}
