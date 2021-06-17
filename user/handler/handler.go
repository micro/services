package handler

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/errors"
	db "github.com/micro/services/db/proto"
	"github.com/micro/services/user/domain"
	pb "github.com/micro/services/user/proto"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

const (
	x = "cruft123"
)

var (
	alphanum = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
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
}

func NewUser(db db.DbService) *User {
	return &User{
		domain: domain.New(db),
	}
}

func (s *User) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	if len(req.Password) < 8 {
		return errors.InternalServerError("user.Create.Check", "Password is less than 8 characters")
	}
	req.Username = strings.ToLower(req.Username)
	req.Email = strings.ToLower(req.Email)

	salt := random(16)
	h, err := bcrypt.GenerateFromPassword([]byte(x+salt+req.Password), 10)
	if err != nil {
		return errors.InternalServerError("user.Create", err.Error())
	}
	pp := base64.StdEncoding.EncodeToString(h)
	if req.Id == "" {
		req.Id = uuid.New().String()
	}
	err = s.domain.Create(ctx, &pb.Account{
		Id:       req.Id,
		Username: req.Username,
		Email:    req.Email,
	}, salt, pp)
	if err != nil {
		return err
	}
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
		return errors.InternalServerError("user.updatepassword", "Passwords don't math")
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
	userId, err := s.domain.ReadToken(ctx, req.Token)
	if err != nil {
		return err
	}
	user, err := s.domain.Read(ctx, userId)
	user.Verified = true
	return s.domain.Update(ctx, user)
}

func (s *User) SendVerificationEmail(ctx context.Context, req *pb.SendVerificationEmailRequest, rsp *pb.SendVerificationEmailResponse) error {
	users, err := s.domain.Search(ctx, "", req.Email)
	if err != nil {
		return err
	}
	token, err := s.domain.CreateToken(ctx, users[0].Id)
	if err != nil {
		return err
	}
	return s.domain.SendEmail(req.Email, users[0].Username, req.Subject, req.TextContent, token)
}
