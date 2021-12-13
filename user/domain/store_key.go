package domain

import "fmt"

func generateAccountStoreKey(userId string) string {
	return fmt.Sprintf("user/account/id/%s", userId)
}

func generateAccountEmailStoreKey(email string) string {
	return fmt.Sprintf("user/acccount/email/%s", email)
}

func generateAccountUsernameStoreKey(username string) string {
	return fmt.Sprintf("user/account/username/%s", username)
}

func generatePasswordStoreKey(userId string) string {
	return fmt.Sprintf("user/password/%s", userId)
}

func generatePasswordResetCodeStoreKey(userId, code string) string {
	return fmt.Sprintf("user/password-reset-codes/%s-%s", userId, code)
}

func generateSessionStoreKey(sessionId string) string {
	return fmt.Sprintf("user/session/%s", sessionId)
}

func generateVerificationsTokenStoreKey(userId, token string) string {
	return fmt.Sprintf("user/verifycation-token/%s-%s", userId, token)
}
