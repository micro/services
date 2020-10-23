package handler

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
	"time"

	"github.com/micro/dev/model"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/store"
	"github.com/micro/services/users/db"
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

type Users struct{}

func NewUsers() *Users {
	nameIndex := model.ByEquality("name")
	nameIndex.Unique = true
	emailIndex := model.ByEquality("email")
	emailIndex.Unique = true

	db.Accounts = model.NewDB(store.DefaultStore, "users", model.Indexes(nameIndex, emailIndex), nil)
	db.Sessions = model.NewDB(store.DefaultStore, "sessions", nil, nil)
	db.Passwords = model.NewDB(store.DefaultStore, "passwords", nil, nil)
	return &Users{}

}

func (s *Users) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	salt := random(16)
	h, err := bcrypt.GenerateFromPassword([]byte(x+salt+req.Password), 10)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.user.Create", err.Error())
	}
	pp := base64.StdEncoding.EncodeToString(h)

	req.User.Username = strings.ToLower(req.User.Username)
	req.User.Email = strings.ToLower(req.User.Email)
	return db.Create(req.User, salt, pp)
}

func (s *Users) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	user, err := db.Read(req.Id)
	if err != nil {
		return err
	}
	rsp.User = user
	return nil
}

func (s *Users) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	req.User.Username = strings.ToLower(req.User.Username)
	req.User.Email = strings.ToLower(req.User.Email)
	return db.Update(req.User)
}

func (s *Users) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	return db.Delete(req.Id)
}

func (s *Users) Search(ctx context.Context, req *pb.SearchRequest, rsp *pb.SearchResponse) error {
	users, err := db.Search(req.Username, req.Email, req.Limit, req.Offset)
	if err != nil {
		return err
	}
	rsp.Users = users
	return nil
}

func (s *Users) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordRequest, rsp *pb.UpdatePasswordResponse) error {
	usr, err := db.Read(req.UserId)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.user.updatepassword", err.Error())
	}

	salt, hashed, err := db.SaltAndPassword(usr.Username, usr.Email)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.user.updatepassword", err.Error())
	}

	hh, err := base64.StdEncoding.DecodeString(hashed)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.user.updatepassword", err.Error())
	}

	if err := bcrypt.CompareHashAndPassword(hh, []byte(x+salt+req.OldPassword)); err != nil {
		return errors.Unauthorized("go.micro.srv.user.updatepassword", err.Error())
	}

	salt = random(16)
	h, err := bcrypt.GenerateFromPassword([]byte(x+salt+req.NewPassword), 10)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.user.updatepassword", err.Error())
	}
	pp := base64.StdEncoding.EncodeToString(h)

	if err := db.UpdatePassword(req.UserId, salt, pp); err != nil {
		return errors.InternalServerError("go.micro.srv.user.updatepassword", err.Error())
	}
	return nil
}

func (s *Users) Login(ctx context.Context, req *pb.LoginRequest, rsp *pb.LoginResponse) error {
	username := strings.ToLower(req.Username)
	email := strings.ToLower(req.Email)

	salt, hashed, err := db.SaltAndPassword(username, email)
	if err != nil {
		return err
	}

	hh, err := base64.StdEncoding.DecodeString(hashed)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.user.Login", err.Error())
	}

	if err := bcrypt.CompareHashAndPassword(hh, []byte(x+salt+req.Password)); err != nil {
		return errors.Unauthorized("go.micro.srv.user.login", err.Error())
	}
	// save session
	sess := &pb.Session{
		Id:       random(128),
		Username: username,
		Created:  time.Now().Unix(),
		Expires:  time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	if err := db.CreateSession(sess); err != nil {
		return errors.InternalServerError("go.micro.srv.user.Login", err.Error())
	}
	rsp.Session = sess
	return nil
}

func (s *Users) Logout(ctx context.Context, req *pb.LogoutRequest, rsp *pb.LogoutResponse) error {
	return db.DeleteSession(req.SessionId)
}

func (s *Users) ReadSession(ctx context.Context, req *pb.ReadSessionRequest, rsp *pb.ReadSessionResponse) error {
	sess, err := db.ReadSession(req.SessionId)
	if err != nil {
		return err
	}
	rsp.Session = sess
	return nil
}
