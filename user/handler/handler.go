package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	pauth "github.com/micro/services/pkg/auth"
	adminpb "github.com/micro/services/pkg/service/proto"
	"golang.org/x/crypto/bcrypt"

	otp "github.com/micro/services/otp/proto"
	"github.com/micro/services/user/domain"
	pb "github.com/micro/services/user/proto"
)

const (
	x = "cruft123"
)

var (
	alphanum    = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	emailFormat = regexp.MustCompile("^[\\w-\\.\\+]+@([\\w-]+\\.)+[\\w-]{2,4}$")
)

// random generate i length alphanum string
func random(i int) string {
	bytes := make([]byte, i)
	for {
		rand.Read(bytes)
		for i, b := range bytes {
			bytes[i] = alphanum[b%byte(len(alphanum))]
		}
		return string(bytes)
	}
	return "ughwhy?!!!"
}

type User struct {
	domain *domain.Domain
	Otp    otp.OtpService
}

type userPayload struct {
	Id       string
	Email    string
	Username string
}

func NewUser(st store.Store, otp otp.OtpService) *User {
	return &User{
		domain: domain.New(st),
		Otp:    otp,
	}
}

// validatePostUserData trims leading and trailing spaces of userId, username, email
// also, check for email format and make sure that values are not empty.
func (s *User) validatePostUserData(ctx context.Context, p *userPayload) error {
	p.Username = strings.TrimSpace(strings.ToLower(p.Username))
	p.Email = strings.TrimSpace(strings.ToLower(p.Email))
	p.Id = strings.TrimSpace(p.Id)

	// email format check
	if !emailFormat.MatchString(p.Email) {
		return errors.BadRequest("users-email-format-check", "email has wrong format")
	}

	if p.Id == "" || p.Username == "" || p.Email == "" {
		return errors.BadRequest("users-valid-check", "missing id or username or email")
	}

	return nil
}

func (s *User) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	check := func(err error) error {
		if err != nil && err.Error() != domain.ErrNotFound.Error() {
			return err
		}
		return nil
	}

	if len(req.Password) < 8 {
		return errors.InternalServerError("users-password-check", "Password is less than 8 characters")
	}

	// based on the docs Id is optional, hence the need to provide
	// one in case of absence
	if req.Id == "" {
		req.Id = uuid.New().String()
	}

	p := &userPayload{Id: req.Id, Email: req.Email, Username: req.Username}

	if err := s.validatePostUserData(ctx, p); err != nil {
		return err
	}

	// userId check
	account, err := s.domain.SearchByUserId(ctx, p.Id)
	if check(err) != nil {
		return err
	}

	if account.Id != "" && account.Id == p.Id {
		return errors.BadRequest("users-userId-check", "account already exists")
	}

	// email check
	account, err = s.domain.SearchByEmail(ctx, p.Email)
	if check(err) != nil {
		return err
	}

	if account.Id != "" && account.Email == p.Email {
		return errors.BadRequest("users-email-check", "email already exists")
	}

	// username check
	account, err = s.domain.SearchByUsername(ctx, p.Username)
	if check(err) != nil {
		return err
	}

	if account.Id != "" && account.Username == p.Username {
		return errors.BadRequest("users-username-check", "username already exists")
	}

	salt := random(16)
	h, err := bcrypt.GenerateFromPassword([]byte(x+salt+req.Password), 10)
	if err != nil {
		return errors.InternalServerError("users-Create", err.Error())
	}

	pp := base64.StdEncoding.EncodeToString(h)

	acc := &pb.Account{
		Id:       p.Id,
		Username: p.Username,
		Email:    p.Email,
		Profile:  req.Profile,
	}

	err = s.domain.Create(ctx, acc, salt, pp)
	if err != nil {
		return err
	}

	// return the account
	rsp.Account = acc

	return nil
}

func (s *User) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	var account = &pb.Account{}
	var err error

	switch {
	case req.Id != "":
		account, err = s.domain.Read(ctx, req.Id)
	case req.Username != "":
		account, err = s.domain.SearchByUsername(ctx, req.Username)
	case req.Email != "":
		account, err = s.domain.SearchByEmail(ctx, req.Email)
	}

	if account == nil {
		return errors.NotFound("user.read", "user not found")
	}

	rsp.Account = account
	if err != nil {
		return err
	}

	return nil
}

func (s *User) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {

	// based on the docs, Update allows user to update email or username.
	// here we handle three cases, when Id and Email are provided,
	// when Id and Username are provided and lastly, when Id, Email and Username are provided.

	check := func(err error) error {
		if err != nil && err.Error() != domain.ErrNotFound.Error() {
			return err
		}
		return nil
	}

	// fetch account
	account, err := s.domain.SearchByUserId(ctx, req.Id)
	if err != nil {
		return err
	}

	// check if req.Email is empty and replace it with account.Email
	// in case of absence, this is neccessary step to prevent validatePostUserData
	// form throwing an error
	if req.Email == "" {
		req.Email = account.Email

		p := &userPayload{Id: req.Id, Email: req.Email, Username: req.Username}

		if err := s.validatePostUserData(ctx, p); err != nil {
			return err
		}

		// check if the new Username is already exists in thge store
		account, err = s.domain.SearchByUsername(ctx, p.Username)
		if check(err) != nil {
			return err
		}

		if account.Id != "" && account.Username == p.Username {
			return errors.BadRequest("users-username-check", "username already exists")
		}

		return s.domain.Update(ctx, &pb.Account{
			Id:       p.Id,
			Username: p.Username,
			Email:    p.Email,
			Profile:  req.Profile,
		})
	}

	// check if req.Username is empty, same as above
	if req.Username == "" {
		req.Username = account.Username

		p := &userPayload{Id: req.Id, Email: req.Email, Username: req.Username}

		if err := s.validatePostUserData(ctx, p); err != nil {
			return err
		}

		// check if the new Email is already exists in the store
		account, err = s.domain.SearchByEmail(ctx, p.Email)
		if check(err) != nil {
			return err
		}

		if account.Id != "" && account.Email == p.Email {
			return errors.BadRequest("users-email-check", "email already exists")
		}

		return s.domain.Update(ctx, &pb.Account{
			Id:       p.Id,
			Username: p.Username,
			Email:    p.Email,
			Profile:  req.Profile,
		})
	}

	// if both new Email and new Username were provided
	p := &userPayload{Id: req.Id, Email: req.Email, Username: req.Username}

	if err := s.validatePostUserData(ctx, p); err != nil {
		return err
	}

	// check if the new Email is already exists in the store
	account, err = s.domain.SearchByEmail(ctx, p.Email)
	if check(err) != nil {
		return err
	}

	if account.Id != "" && account.Email == p.Email {
		return errors.BadRequest("users-email-check", "email already exists")
	}

	// check if the new Username is already exists in thge store
	account, err = s.domain.SearchByUsername(ctx, p.Username)
	if check(err) != nil {
		return err
	}

	if account.Id != "" && account.Username == p.Username {
		return errors.BadRequest("users-username-check", "username already exists")
	}

	return s.domain.Update(ctx, &pb.Account{
		Id:       p.Id,
		Username: p.Username,
		Email:    p.Email,
		Profile:  req.Profile,
	})
}

func (s *User) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	return s.domain.Delete(ctx, req.Id)
}

func (s *User) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordRequest, rsp *pb.UpdatePasswordResponse) error {
	usr, err := s.domain.Read(ctx, req.UserId)
	if err != nil {
		return errors.InternalServerError("user.updatepassword", err.Error())
	}
	if req.NewPassword != req.ConfirmPassword {
		return errors.InternalServerError("user.updatepassword", "Passwords don't match")
	}

	salt, hashed, err := s.domain.SaltAndPassword(ctx, usr.Id)
	if err != nil {
		return errors.InternalServerError("user.updatepassword", err.Error())
	}

	hh, err := base64.StdEncoding.DecodeString(hashed)
	if err != nil {
		return errors.InternalServerError("user.updatepassword", err.Error())
	}

	if err := bcrypt.CompareHashAndPassword(hh, []byte(x+salt+req.OldPassword)); err != nil {
		return errors.Unauthorized("user.updatepassword", err.Error())
	}

	salt = random(16)
	h, err := bcrypt.GenerateFromPassword([]byte(x+salt+req.NewPassword), 10)
	if err != nil {
		return errors.InternalServerError("user.updatepassword", err.Error())
	}
	pp := base64.StdEncoding.EncodeToString(h)

	if err := s.domain.UpdatePassword(ctx, req.UserId, salt, pp); err != nil {
		return errors.InternalServerError("user.updatepassword", err.Error())
	}
	return nil
}

func (s *User) Login(ctx context.Context, req *pb.LoginRequest, rsp *pb.LoginResponse) error {
	username := strings.ToLower(req.Username)
	email := strings.ToLower(req.Email)

	accounts, err := s.domain.Search(ctx, username, email)
	if err != nil {
		return err
	}
	if len(accounts) == 0 {
		return fmt.Errorf("account not found")
	}
	salt, hashed, err := s.domain.SaltAndPassword(ctx, accounts[0].Id)
	if err != nil {
		return err
	}

	hh, err := base64.StdEncoding.DecodeString(hashed)
	if err != nil {
		return errors.InternalServerError("user.Login", err.Error())
	}

	if err := bcrypt.CompareHashAndPassword(hh, []byte(x+salt+req.Password)); err != nil {
		return errors.Unauthorized("user.login", err.Error())
	}
	// save session
	sess := &pb.Session{
		Id:      random(128),
		Created: time.Now().Unix(),
		Expires: time.Now().Add(time.Hour * 24 * 7).Unix(),
		UserId:  accounts[0].Id,
	}

	if err := s.domain.CreateSession(ctx, sess); err != nil {
		return errors.InternalServerError("user.Login", err.Error())
	}
	rsp.Session = sess
	return nil
}

func (s *User) Logout(ctx context.Context, req *pb.LogoutRequest, rsp *pb.LogoutResponse) error {
	return s.domain.DeleteSession(ctx, req.SessionId)
}

func (s *User) ReadSession(ctx context.Context, req *pb.ReadSessionRequest, rsp *pb.ReadSessionResponse) error {
	sess, err := s.domain.ReadSession(ctx, req.SessionId)
	if err != nil {
		return err
	}
	rsp.Session = sess
	return nil
}

func (s *User) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest, rsp *pb.VerifyEmailResponse) error {
	if len(req.Token) == 0 {
		return errors.BadRequest("user.verifytoken", "missing token")
	}

	// check the token exists
	tenant, email, err := s.domain.ReadToken(ctx, req.Token)
	if err != nil {
		logger.Error("Failed to read token: %v", err)
		return err
	}

	// update the user
	err = s.domain.MarkVerified(ctx, tenant, email)
	if err != nil {
		logger.Error("Failed to mark email: %s for tenant: %s as verified: %v", email, tenant, err)
	}
	return err
}

func (s *User) SendVerificationEmail(ctx context.Context, req *pb.SendVerificationEmailRequest, rsp *pb.SendVerificationEmailResponse) error {
	if len(req.Email) == 0 {
		return errors.BadRequest("user.sendverificationemail", "missing email")
	}

	// search for the user
	account, err := s.domain.SearchByEmail(ctx, req.Email)
	if err != nil {
		return err
	}

	// generate random token
	token := random(256)

	// generate/save a token for verification
	err = s.domain.CreateToken(ctx, req.Email, token)
	if err != nil {
		return err
	}

	return s.domain.SendEmail(req.FromName, req.Email, account.Username, req.Subject, req.TextContent, token, req.RedirectUrl, req.FailureRedirectUrl)
}

func (s *User) SendPasswordResetEmail(ctx context.Context, req *pb.SendPasswordResetEmailRequest, rsp *pb.SendPasswordResetEmailResponse) error {
	if len(req.Email) == 0 {
		return errors.BadRequest("user.sendpasswordresetemail", "missing email")
	}
	if len(req.Subject) == 0 {
		return errors.BadRequest("user.sendpasswordresetemail", "missing subject")
	}
	if len(req.TextContent) == 0 {
		return errors.BadRequest("user.sendpasswordresetemail", "missing textContent")
	}

	// look for an existing user
	account, err := s.domain.SearchByEmail(ctx, req.Email)
	if err != nil {
		return err
	}

	var expiry int64 = 1800 // 1800 secs = 30 min
	if req.Expiration > 0 {
		expiry = req.Expiration
	}

	if err != nil {
		return err
	}
	code := random(8)

	// save the password reset code
	_, err = s.domain.SavePasswordResetCode(ctx, account.Id, code, time.Duration(expiry)*time.Second)
	if err != nil {
		return err
	}
	// save the code in the database and then send via email
	return s.domain.SendPasswordResetEmail(ctx, account.Id, code, req.FromName, req.Email, account.Username, req.Subject, req.TextContent)
}

func (s *User) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest, rsp *pb.ResetPasswordResponse) error {
	if len(req.Email) == 0 {
		return errors.BadRequest("user.resetpassword", "missing email")
	}
	if len(req.Code) == 0 {
		return errors.BadRequest("user.resetpassword", "missing code")
	}
	if len(req.ConfirmPassword) == 0 {
		return errors.BadRequest("user.resetpassword", "missing confirm password")
	}
	if len(req.NewPassword) == 0 {
		return errors.BadRequest("user.resetpassword", "missing new password")
	}
	if req.ConfirmPassword != req.NewPassword {
		return errors.BadRequest("user.resetpassword", "passwords do not match")
	}

	// look for an existing user
	account, err := s.domain.SearchByEmail(ctx, req.Email)
	if err != nil {
		return err
	}

	// validate the code
	_, err = s.domain.ReadPasswordResetCode(ctx, account.Id, req.Code)
	if err != nil {
		return err
	}

	// no error means it exists and not expired
	salt := random(16)
	h, err := bcrypt.GenerateFromPassword([]byte(x+salt+req.NewPassword), 10)
	if err != nil {
		return errors.InternalServerError("user.ResetPassword", err.Error())
	}
	pp := base64.StdEncoding.EncodeToString(h)

	// update the user password
	if err := s.domain.UpdatePassword(ctx, account.Id, salt, pp); err != nil {
		return errors.InternalServerError("user.resetpassword", err.Error())
	}

	// delete our saved code
	s.domain.DeletePasswordResetCode(ctx, account.Id, req.Code)

	return nil
}

func (s *User) List(ctx context.Context, request *pb.ListRequest, response *pb.ListResponse) error {
	accs, err := s.domain.List(ctx, request.Offset, request.Limit)
	if err != nil && err != domain.ErrNotFound {
		return errors.InternalServerError("user.List", "Error retrieving user list")
	}
	response.Users = make([]*pb.Account, len(accs))
	for i, v := range accs {
		response.Users[i] = v
	}
	return nil
}

func (s *User) SendMagicLink(ctx context.Context, req *pb.SendMagicLinkRequest, rsp *pb.SendMagicLinkResponse) error {
	// check if the email has the correct format
	if !emailFormat.MatchString(req.Email) {
		return errors.BadRequest("SendMagicLink.email-format-check", "email has wrong format")
	}

	// check if the email exist in the DB
	account, err := s.domain.SearchByEmail(ctx, req.Email)
	if err != nil && err.Error() == "not found" {
		return errors.BadRequest("SendMagicLink.email-check", "email doesn't exist")
	} else if err != nil {
		return errors.BadRequest("SendMagicLink.email-check", err.Error())
	}

	// create a token object
	token := random(128)

	// set ttl to 60 seconds
	ttl := 60

	// save token, so we can retrieve it later
	err = s.domain.CacheToken(ctx, token, req.Email, ttl)
	if err != nil {
		logger.Errorf("SendMagicLink.cacheToken error: %v", err)
		return errors.BadRequest("SendMagicLink.cacheToken", "Oooops something went wrong")
	}

	// send magic link to email address
	err = s.domain.SendMLE(req.FromName, req.Email, account.Username, req.Subject, req.TextContent, token, req.Address, req.Endpoint)
	if err != nil {
		logger.Errorf("SendMagicLink.cacheToken error: %v", err)
		return errors.BadRequest("SendMagicLink.sendEmail", "Oooops something went wrong")
	}

	return nil
}

func (s *User) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest, rsp *pb.VerifyTokenResponse) error {
	// extract token
	token := req.Token

	// check if token is valid
	email, err := s.domain.CacheReadToken(ctx, token)
	if err != nil && err.Error() == "token not found" {
		rsp.IsValid = false
		rsp.Message = err.Error()
		return nil
	} else if err != nil && err.Error() == "token expired" {
		rsp.IsValid = false
		rsp.Message = err.Error()
		return nil
	} else if err != nil && err.Error() == "token empty" {
		rsp.IsValid = false
		rsp.Message = err.Error()
		return nil
	} else if err != nil {
		rsp.IsValid = false
		rsp.Message = err.Error()
		return errors.BadRequest("VerifyToken.CacheReadToken", err.Error())
	}

	// save session
	account, err := s.domain.SearchByEmail(ctx, email)
	if err != nil {
		rsp.IsValid = false
		rsp.Message = "account not found"
		return err
	}

	sess := &pb.Session{
		Id:      random(128),
		Created: time.Now().Unix(),
		Expires: time.Now().Add(time.Hour * 24 * 7).Unix(),
		UserId:  account.Id,
	}

	if err := s.domain.CreateSession(ctx, sess); err != nil {
		rsp.IsValid = false
		rsp.Message = "Creation of a new session has failed"
		return errors.InternalServerError("VerifyToken.createSession", err.Error())
	}

	rsp.IsValid = true
	rsp.Session = sess

	return nil
}

func (s *User) DeleteData(ctx context.Context, request *adminpb.DeleteDataRequest, response *adminpb.DeleteDataResponse) error {
	if _, err := pauth.VerifyMicroAdmin(ctx, "user.DeleteData"); err != nil {
		return err
	}

	if len(request.TenantId) < 10 { // deliberate length check so we don't delete all the things
		return errors.BadRequest("user.DeleteData", "Missing tenant ID")
	}
	return s.domain.DeleteTenantData(request.TenantId)
}
