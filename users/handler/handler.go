package handler

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
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

// Create a user
func (u *Users) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	// validate the request
	if len(req.FirstName) == 0 {
		return ErrMissingFirstName
	}
	if len(req.LastName) == 0 {
		return ErrMissingLastName
	}
	if len(req.Email) == 0 {
		return ErrMissingEmail
	}
	if !isEmailValid(req.Email) {
		return ErrInvalidEmail
	}
	if len(req.Password) < 8 {
		return ErrInvalidPassword
	}

	// hash and salt the password using bcrypt
	phash, err := hashAndSalt(req.Password)
	if err != nil {
		logger.Errorf("Error hasing and salting password: %v", err)
		return errors.InternalServerError("HASHING_ERROR", "Error hashing password")
	}

	return u.DB.Transaction(func(tx *gorm.DB) error {
		// write the user to the database
		user := &User{
			ID:        uuid.New().String(),
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Password:  phash,
		}
		err = u.DB.Create(user).Error
		if err != nil && strings.Contains(err.Error(), "idx_users_email") {
			return ErrDuplicateEmail
		} else if err != nil {
			logger.Errorf("Error writing to the database: %v", err)
			return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
		}

		// generate a token for the user
		token := Token{
			UserID:    user.ID,
			Key:       uuid.New().String(),
			ExpiresAt: u.Time().Add(time.Hour * 24 * 7),
		}
		if err := tx.Create(&token).Error; err != nil {
			logger.Errorf("Error writing to the database: %v", err)
			return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
		}

		// serialize the response
		rsp.User = user.Serialize()
		rsp.Token = token.Key
		return nil
	})
}

// Read users using ID
func (u *Users) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	// validate the request
	if len(req.Ids) == 0 {
		return ErrMissingIDs
	}

	// query the database
	var users []User
	if err := u.DB.Model(&User{}).Where("id IN (?)", req.Ids).Find(&users).Error; err != nil {
		logger.Errorf("Error reading from the database: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
	}

	// serialize the response
	rsp.Users = make(map[string]*pb.User, len(users))
	for _, u := range users {
		rsp.Users[u.ID] = u.Serialize()
	}
	return nil
}

// Update a user
func (u *Users) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	// validate the request
	if len(req.Id) == 0 {
		return ErrMissingID
	}
	if req.FirstName != nil && len(req.FirstName.Value) == 0 {
		return ErrMissingFirstName
	}
	if req.LastName != nil && len(req.LastName.Value) == 0 {
		return ErrMissingLastName
	}
	if req.Email != nil && len(req.Email.Value) == 0 {
		return ErrMissingEmail
	}
	if req.Email != nil && !isEmailValid(req.Email.Value) {
		return ErrInvalidEmail
	}

	// lookup the user
	var user User
	if err := u.DB.Where(&User{ID: req.Id}).First(&user).Error; err == gorm.ErrRecordNotFound {
		return ErrNotFound
	} else if err != nil {
		logger.Errorf("Error reading from the database: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
	}

	// assign the updated values
	if req.FirstName != nil {
		user.FirstName = req.FirstName.Value
	}
	if req.LastName != nil {
		user.LastName = req.LastName.Value
	}
	if req.Email != nil {
		user.Email = req.Email.Value
	}

	// write the user to the database
	err := u.DB.Save(user).Error
	if err != nil && strings.Contains(err.Error(), "idx_users_email") {
		return ErrDuplicateEmail
	} else if err != nil {
		logger.Errorf("Error writing to the database: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
	}

	// serialize the user
	rsp.User = user.Serialize()
	return nil
}

// Delete a user
func (u *Users) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	// validate the request
	if len(req.Id) == 0 {
		return ErrMissingID
	}

	// delete the users tokens
	return u.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&Token{}, &Token{UserID: req.Id}).Error; err != nil {
			logger.Errorf("Error writing to the database: %v", err)
			return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
		}

		// delete from the database
		if err := tx.Delete(&User{}, &User{ID: req.Id}).Error; err != nil {
			logger.Errorf("Error writing to the database: %v", err)
			return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
		}

		return nil
	})
}

// List all users
func (u *Users) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	// query the database
	var users []User
	if err := u.DB.Model(&User{}).Find(&users).Error; err != nil {
		logger.Errorf("Error reading from the database: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
	}

	// serialize the response
	rsp.Users = make([]*pb.User, len(users))
	for i, u := range users {
		rsp.Users[i] = u.Serialize()
	}
	return nil
}

// Login using email and password returns the users profile and a token
func (u *Users) Login(ctx context.Context, req *pb.LoginRequest, rsp *pb.LoginResponse) error {
	// validate the request
	if len(req.Email) == 0 {
		return ErrMissingEmail
	}
	if len(req.Password) == 0 {
		return ErrInvalidPassword
	}

	return u.DB.Transaction(func(tx *gorm.DB) error {
		// lookup the user
		var user User
		if err := tx.Where(&User{Email: req.Email}).First(&user).Error; err == gorm.ErrRecordNotFound {
			return ErrNotFound
		} else if err != nil {
			logger.Errorf("Error reading from the database: %v", err)
			return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
		}

		// compare the passwords
		if !passwordsMatch(user.Password, req.Password) {
			return ErrIncorrectPassword
		}

		// generate a token for the user
		token := Token{
			UserID:    user.ID,
			Key:       uuid.New().String(),
			ExpiresAt: u.Time().Add(tokenTTL),
		}
		if err := tx.Create(&token).Error; err != nil {
			logger.Errorf("Error writing to the database: %v", err)
			return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
		}

		// serialize the response
		rsp.Token = token.Key
		rsp.User = user.Serialize()
		return nil
	})
}

// Logout expires all tokens for the user
func (u *Users) Logout(ctx context.Context, req *pb.LogoutRequest, rsp *pb.LogoutResponse) error {
	// validate the request
	if len(req.Id) == 0 {
		return ErrMissingID
	}

	return u.DB.Transaction(func(tx *gorm.DB) error {
		// lookup the user
		var user User
		if err := tx.Where(&User{ID: req.Id}).Preload("Tokens").First(&user).Error; err == gorm.ErrRecordNotFound {
			return ErrNotFound
		} else if err != nil {
			logger.Errorf("Error reading from the database: %v", err)
			return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
		}

		// delete the tokens
		if err := tx.Delete(user.Tokens).Error; err != nil {
			logger.Errorf("Error deleting from the database: %v", err)
			return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
		}

		return nil
	})
}

// Validate a token, each time a token is validated it extends its lifetime for another week
func (u *Users) Validate(ctx context.Context, req *pb.ValidateRequest, rsp *pb.ValidateResponse) error {
	// validate the request
	if len(req.Token) == 0 {
		return ErrMissingToken
	}

	return u.DB.Transaction(func(tx *gorm.DB) error {
		// lookup the token
		var token Token
		if err := tx.Where(&Token{Key: req.Token}).Preload("User").First(&token).Error; err == gorm.ErrRecordNotFound {
			return ErrInvalidToken
		} else if err != nil {
			logger.Errorf("Error reading from the database: %v", err)
			return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
		}

		// ensure the token is valid
		if u.Time().After(token.ExpiresAt) {
			return ErrTokenExpired
		}

		// extend the token for another lifetime
		token.ExpiresAt = u.Time().Add(tokenTTL)
		if err := tx.Save(&token).Error; err != nil {
			logger.Errorf("Error writing to the database: %v", err)
			return errors.InternalServerError("DATABASE_ERROR", "Error connecting to the database")
		}

		// serialize the response
		rsp.User = token.User.Serialize()
		return nil
	})
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
