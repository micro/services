package handler

import (
	"regexp"
	"time"

	"github.com/micro/micro/v3/service/errors"
	pb "github.com/micro/services/users/proto"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrMissingFirstName  = errors.BadRequest("MISSING_FIRST_NAME", "Missing first name")
	ErrMissingLastName   = errors.BadRequest("MISSING_LAST_NAME", "Missing last name")
	ErrMissingEmail      = errors.BadRequest("MISSING_EMAIL", "Missing email")
	ErrDuplicateEmail    = errors.BadRequest("DUPLICATE_EMAIL", "A user with this email address already exists")
	ErrInvalidEmail      = errors.BadRequest("INVALID_EMAIL", "The email provided is invalid")
	ErrInvalidPassword   = errors.BadRequest("INVALID_PASSWORD", "Password must be at least 8 characters long")
	ErrMissingEmails     = errors.BadRequest("MISSING_EMAILS", "One or more emails are required")
	ErrMissingIDs        = errors.BadRequest("MISSING_IDS", "One or more ids are required")
	ErrMissingID         = errors.BadRequest("MISSING_ID", "Missing ID")
	ErrMissingToken      = errors.BadRequest("MISSING_TOKEN", "Missing token")
	ErrIncorrectPassword = errors.BadRequest("INCORRECT_PASSWORD", "Incorrect password")
	ErrTokenExpired      = errors.BadRequest("TOKEN_EXPIRED", "Token has expired")
	ErrInvalidToken      = errors.BadRequest("INVALID_TOKEN", "Token is invalid")
	ErrNotFound          = errors.NotFound("NOT_FOUND", "User not found")

	emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	tokenTTL   = time.Hour * 7 * 24
)

type User struct {
	ID        string
	FirstName string
	LastName  string
	Email     string `gorm:"uniqueIndex"`
	Password  string
	CreatedAt time.Time
	Tokens    []Token
}

func (u *User) Serialize() *pb.User {
	return &pb.User{
		Id:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
	}
}

type Token struct {
	Key       string `gorm:"primaryKey"`
	CreatedAt time.Time
	ExpiresAt time.Time
	UserID    string
	User      User
}

type Users struct {
	DB   *gorm.DB
	Time func() time.Time
}

// isEmailValid checks if the email provided passes the required structure and length.
func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

func hashAndSalt(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func passwordsMatch(hashed string, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	return err == nil
}
