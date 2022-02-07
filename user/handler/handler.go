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

func NewUser(st store.Store, otp otp.OtpService) *User {
	return &User{
		domain: domain.New(st),
		Otp:    otp,
	}
}

// validatePostUserData checks userId, username, email post data are valid and in right format
func (s *User) validatePostUserData(ctx context.Context, userId, username, email string) error {
	username = strings.TrimSpace(strings.ToLower(username))
	email = strings.TrimSpace(strings.ToLower(email))

	if !emailFormat.MatchString(email) {
		return errors.BadRequest("create.email-format-check", "email has wrong format")
	}

	if userId == "" || username == "" || email == "" {
		return errors.BadRequest("valid-check", "missing id or username or email")
	}

	account, err := s.domain.SearchByUsername(ctx, username)
	if err != nil && err.Error() != domain.ErrNotFound.Error() {
		return err
	}

	if account.Id != "" && account.Username == username {
		return errors.BadRequest("username-check", "username already exists")
	}

	account, err = s.domain.SearchByEmail(ctx, email)
	if err != nil && err.Error() != domain.ErrNotFound.Error() {
		return err
	}

	if account.Id != "" && account.Email == email {
		return errors.BadRequest("email-check", "email already exists")
	}

	return nil
}

func (s *User) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	if len(req.Password) < 8 {
		return errors.InternalServerError("user.Create.Check", "Password is less than 8 characters")
	}

	if err := s.validatePostUserData(ctx, req.Id, req.Username, req.Email); err != nil {
		return err
	}

	salt := random(16)
	h, err := bcrypt.GenerateFromPassword([]byte(x+salt+req.Password), 10)
	if err != nil {
		return errors.InternalServerError("user.Create", err.Error())
	}
	pp := base64.StdEncoding.EncodeToString(h)
	if req.Id == "" {
		req.Id = uuid.New().String()
	}

	acc := &pb.Account{
		Id:       req.Id,
		Username: strings.ToLower(req.Username),
		Email:    strings.ToLower(req.Email),
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

	rsp.Account = account
	if err != nil {
		return err
	}

	return nil
}

func (s *User) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	if err := s.validatePostUserData(ctx, req.Id, req.Username, req.Email); err != nil {
		return err
	}

	return s.domain.Update(ctx, &pb.Account{
		Id:       req.Id,
		Username: strings.ToLower(req.Username),
		Email:    strings.ToLower(req.Email),
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
	if len(req.Email) == 0 {
		return errors.BadRequest("user.verifyemail", "missing email")
	}
	if len(req.Token) == 0 {
		return errors.BadRequest("user.verifyemail", "missing token")
	}

	// check the token exists
	email, err := s.domain.ReadToken(ctx, req.Email, req.Token)
	if err != nil {
		return err
	}

	// validate the code, e.g its an OTP token and hasn't expired
	resp, err := s.Otp.Validate(ctx, &otp.ValidateRequest{
		Id:   req.Email,
		Code: req.Token,
	})
	if err != nil {
		return err
	}

	// check if the code is actually valid
	if !resp.Success {
		return errors.BadRequest("user.resetpassword", "invalid code")
	}

	// mark user as verified
	user, err := s.domain.SearchByEmail(ctx, email)
	if err != nil {
		return err
	}

	user.Verified = true

	// update the user
	return s.domain.Update(ctx, user)
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

	// generate a new OTP code
	resp, err := s.Otp.Generate(ctx, &otp.GenerateRequest{
		Expiry: 900,
		Id:     req.Email,
	})

	if err != nil {
		return err
	}

	// generate/save a token for verification
	token, err := s.domain.CreateToken(ctx, req.Email, resp.Code)
	if err != nil {
		return err
	}

	return s.domain.SendEmail(req.FromName, req.Email, account.Username, req.Subject, req.TextContent, token, req.RedirectUrl, req.FailureRedirectUrl)
}

func (s *User) SendPasswordResetEmail(ctx context.Context, req *pb.SendPasswordResetEmailRequest, rsp *pb.SendPasswordResetEmailResponse) error {
	if len(req.Email) == 0 {
		return errors.BadRequest("user.sendpasswordresetemail", "missing email")
	}

	// look for an existing user
	account, err := s.domain.SearchByEmail(ctx, req.Email)
	if err != nil {
		return err
	}

	// generate a new OTP code
	resp, err := s.Otp.Generate(ctx, &otp.GenerateRequest{
		Expiry: 900,
		Id:     req.Email,
	})

	if err != nil {
		return err
	}

	// save the code in the database and then send via email
	return s.domain.SendPasswordResetEmail(ctx, account.Id, resp.Code, req.FromName, req.Email, account.Username, req.Subject, req.TextContent)
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

	// check if a request was made to reset the password, we should have saved it
	code, err := s.domain.ReadPasswordResetCode(ctx, account.Id, req.Code)
	if err != nil {
		return err
	}

	// validate the code, e.g its an OTP token and hasn't expired
	resp, err := s.Otp.Validate(ctx, &otp.ValidateRequest{
		Id:   req.Email,
		Code: req.Code,
	})
	if err != nil {
		return err
	}

	// check if the code is actually valid
	if !resp.Success {
		return errors.BadRequest("user.resetpassword", "invalid code")
	}

	// no error means it exists and not expired
	salt := random(16)
	h, err := bcrypt.GenerateFromPassword([]byte(x+salt+req.NewPassword), 10)
	if err != nil {
		return errors.InternalServerError("user.ResetPassword", err.Error())
	}
	pp := base64.StdEncoding.EncodeToString(h)

	// update the user password
	if err := s.domain.UpdatePassword(ctx, code.UserID, salt, pp); err != nil {
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
