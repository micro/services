package handler

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
	"time"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/services/users/domain"
	pb "github.com/micro/services/users/proto"
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

type Users struct {
	domain *domain.Domain
}

func NewUsers() *Users {
	return &Users{
		domain: domain.New(),
	}
}

func (s *Users) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	if len(req.Password) < 8 {
		return errors.InternalServerError("users.Create.Check", "Password is less than 8 characters")
	}
	salt := random(16)
	h, err := bcrypt.GenerateFromPassword([]byte(x+salt+req.Password), 10)
	if err != nil {
		return errors.InternalServerError("users.Create", err.Error())
	}
	pp := base64.StdEncoding.EncodeToString(h)

	return s.domain.Create(&pb.User{
		Id:       req.Id,
		Username: strings.ToLower(req.Username),
		Email:    strings.ToLower(req.Email),
	}, salt, pp)
}

func (s *Users) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	user, err := s.domain.Read(req.Id)
	if err != nil {
		return err
	}
	rsp.User = user
	return nil
}

func (s *Users) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	return s.domain.Update(&pb.User{
		Id:       req.Id,
		Username: strings.ToLower(req.Username),
		Email:    strings.ToLower(req.Email),
	})
}

func (s *Users) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	return s.domain.Delete(req.Id)
}

func (s *Users) Search(ctx context.Context, req *pb.SearchRequest, rsp *pb.SearchResponse) error {
	users, err := s.domain.Search(req.Username, req.Email, req.Limit, req.Offset)
	if err != nil {
		return err
	}
	rsp.Users = users
	return nil
}

func (s *Users) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordRequest, rsp *pb.UpdatePasswordResponse) error {
	usr, err := s.domain.Read(req.UserId)
	if err != nil {
		return errors.InternalServerError("users.updatepassword", err.Error())
	}
	if req.NewPassword != req.ConfirmPassword {
		return errors.InternalServerError("users.updatepassword", "Passwords don't math")
	}

	salt, hashed, err := s.domain.SaltAndPassword(usr.Username, usr.Email)
	if err != nil {
		return errors.InternalServerError("users.updatepassword", err.Error())
	}

	hh, err := base64.StdEncoding.DecodeString(hashed)
	if err != nil {
		return errors.InternalServerError("users.updatepassword", err.Error())
	}

	if err := bcrypt.CompareHashAndPassword(hh, []byte(x+salt+req.OldPassword)); err != nil {
		return errors.Unauthorized("users.updatepassword", err.Error())
	}

	salt = random(16)
	h, err := bcrypt.GenerateFromPassword([]byte(x+salt+req.NewPassword), 10)
	if err != nil {
		return errors.InternalServerError("users.updatepassword", err.Error())
	}
	pp := base64.StdEncoding.EncodeToString(h)

	if err := s.domain.UpdatePassword(req.UserId, salt, pp); err != nil {
		return errors.InternalServerError("users.updatepassword", err.Error())
	}
	return nil
}

func (s *Users) Login(ctx context.Context, req *pb.LoginRequest, rsp *pb.LoginResponse) error {
	username := strings.ToLower(req.Username)
	email := strings.ToLower(req.Email)

	salt, hashed, err := s.domain.SaltAndPassword(username, email)
	if err != nil {
		return err
	}

	hh, err := base64.StdEncoding.DecodeString(hashed)
	if err != nil {
		return errors.InternalServerError("users.Login", err.Error())
	}

	if err := bcrypt.CompareHashAndPassword(hh, []byte(x+salt+req.Password)); err != nil {
		return errors.Unauthorized("users.login", err.Error())
	}
	// save session
	sess := &pb.Session{
		Id:       random(128),
		Username: username,
		Email:    email,
		Created:  time.Now().Unix(),
		Expires:  time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	if err := s.domain.CreateSession(sess); err != nil {
		return errors.InternalServerError("users.Login", err.Error())
	}
	rsp.Session = sess
	return nil
}

func (s *Users) Logout(ctx context.Context, req *pb.LogoutRequest, rsp *pb.LogoutResponse) error {
	return s.domain.DeleteSession(req.SessionId)
}

func (s *Users) ReadSession(ctx context.Context, req *pb.ReadSessionRequest, rsp *pb.ReadSessionResponse) error {
	sess, err := s.domain.ReadSession(req.SessionId)
	if err != nil {
		return err
	}
	rsp.Session = sess
	return nil
}
