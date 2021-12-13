package handler

import (
	goctx "context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/errors"
	db "github.com/micro/services/db/proto"
	otp "github.com/micro/services/otp/proto"
	"github.com/micro/services/user/domain"
	pb "github.com/micro/services/user/proto"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

const (
	x = "cruft123"
)

var (
	alphanum    = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	emailFormat = regexp.MustCompile("^[\\w-\\.\\+]+@([\\w-]+\\.)+[\\w-]{2,4}$")
)

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

func NewUser(db db.DbService, otp otp.OtpService) *User {
	return &User{
		domain: domain.New(db),
		Otp:    otp,
	}
}

func (s *User) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	if !emailFormat.MatchString(req.Email) {
		return errors.BadRequest("create.email-format-check", "email has wrong format")
	}
	if len(req.Password) < 8 {
		return errors.InternalServerError("user.Create.Check", "Password is less than 8 characters")
	}
	req.Username = strings.ToLower(req.Username)
	req.Email = strings.ToLower(req.Email)
	usernames, err := s.domain.Search(ctx, req.Username, "")
	if err != nil && err.Error() != "not found" {
		return err
	}
	if len(usernames) > 0 {
		return errors.BadRequest("create.username-check", "username already exists")
	}

	// TODO: don't error out here
	emails, err := s.domain.Search(ctx, "", req.Email)
	if err != nil && err.Error() != "not found" {
		return err
	}
	if len(emails) > 0 {
		return errors.BadRequest("create.email-check", "email already exists")
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
		Username: req.Username,
		Email:    req.Email,
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
	switch {
	case req.Id != "":
		account, err := s.domain.Read(ctx, req.Id)
		if err != nil {
			return err
		}
		rsp.Account = account
		return nil
	case req.Username != "" || req.Email != "":
		accounts, err := s.domain.Search(ctx, req.Username, req.Email)
		if err != nil {
			return err
		}
		rsp.Account = accounts[0]
		return nil
	}
	return nil
}

func (s *User) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
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
	userId, err := s.domain.ReadToken(ctx, req.Email, req.Token)
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
	user, err := s.domain.Read(ctx, userId)
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
	users, err := s.domain.Search(ctx, "", req.Email)
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

	return s.domain.SendEmail(req.FromName, req.Email, users[0].Username, req.Subject, req.TextContent, token, req.RedirectUrl, req.FailureRedirectUrl)
}

func (s *User) SendPasswordResetEmail(ctx context.Context, req *pb.SendPasswordResetEmailRequest, rsp *pb.SendPasswordResetEmailResponse) error {
	if len(req.Email) == 0 {
		return errors.BadRequest("user.sendpasswordresetemail", "missing email")
	}

	// look for an existing user
	users, err := s.domain.Search(ctx, "", req.Email)
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
	return s.domain.SendPasswordResetEmail(ctx, users[0].Id, resp.Code, req.FromName, req.Email, users[0].Username, req.Subject, req.TextContent)
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
	users, err := s.domain.Search(ctx, "", req.Email)
	if err != nil {
		return err
	}

	// check if a request was made to reset the password, we should have saved it
	code, err := s.domain.ReadPasswordResetCode(ctx, users[0].Id, req.Code)
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
	s.domain.DeletePasswordResetCode(ctx, users[0].Id, req.Code)

	return nil
}

func (s *User) List(ctx goctx.Context, request *pb.ListRequest, response *pb.ListResponse) error {
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
	users, err := s.domain.Search(ctx, "", req.Email)
	if err.Error() == "not found" {
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
		return errors.BadRequest("SendMagicLink.cacheToken", "Oooops something went wrong")
	}

	// send magic link to email address
	err = s.domain.SendMLE(req.FromName, req.Email, users[0].Username, req.Subject, req.TextContent, token, req.Address, req.Endpoint)
	if err != nil {
		return errors.BadRequest("SendMagicLink.sendEmail", "Oooops something went wrong")
	}

	return nil
}

func (s *User) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest, rsp *pb.VerifyTokenResponse) error {
	// extract token
	token := req.Token

	// check if token is valid
	email, err := s.domain.CacheReadToken(ctx, token)
	if err.Error() == "token not found" {
		rsp.IsValid = false
		rsp.Message = err.Error()
		return nil
	} else if err.Error() == "token expired" {
		rsp.IsValid = false
		rsp.Message = err.Error()
		return nil
	} else if err.Error() == "token empty" {
		rsp.IsValid = false
		rsp.Message = err.Error()
		return nil
	} else if err != nil {
		rsp.IsValid = false
		rsp.Message = err.Error()
		return errors.BadRequest("VerifyToken.CacheReadToken", err.Error())
	}

	// save session
	accounts, err := s.domain.Search(ctx, "", email)
	if err != nil {
		rsp.IsValid = false
		rsp.Message = "account not found"
		return err
	}
	if len(accounts) == 0 {
		rsp.IsValid = false
		rsp.Message = "account not found"
		return nil
	}

	sess := &pb.Session{
		Id:      random(128),
		Created: time.Now().Unix(),
		Expires: time.Now().Add(time.Hour * 24 * 7).Unix(),
		UserId:  accounts[0].Id,
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
