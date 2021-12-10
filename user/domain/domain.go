package domain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"

	_struct "github.com/golang/protobuf/ptypes/struct"
	"github.com/micro/micro/v3/service/config"
	microerr "github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	db "github.com/micro/services/db/proto"
	"github.com/micro/services/pkg/cache"
	user "github.com/micro/services/user/proto"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var (
	ErrNotFound = errors.New("not found")
)

type pw struct {
	ID       string `json:"id"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

type verificationToken struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`
	Token  string `json:"token"`
}

type passwordResetCode struct {
	ID      string    `json:"id"`
	Expires time.Time `json:"expires"`
	UserID  string    `json:"userId"`
	Code    string    `json:"code"`
}

type Domain struct {
	db         db.DbService
	sengridKey string
	fromEmail  string
}

var (
	// TODO: use the config to drive this value
	defaultSender = "noreply@email.m3ocontent.com"
)

func New(db db.DbService) *Domain {
	var key, email string
	cfg, err := config.Get("micro.user.sendgrid.api_key")
	if err == nil {
		key = cfg.String("")
	}
	cfg, err = config.Get("micro.user.sendgrid.from_email")
	if err == nil {
		email = cfg.String(defaultSender)
	}
	if len(key) == 0 {
		logger.Info("No email key found")
	} else {
		logger.Info("Email key found")
	}
	return &Domain{
		sengridKey: key,
		db:         db,
		fromEmail:  email,
	}
}

func (domain *Domain) SendEmail(fromName, toAddress, toUsername, subject, textContent, token, redirctUrl, failureRedirectUrl string) error {
	if domain.sengridKey == "" {
		return fmt.Errorf("empty email api key")
	}
	from := mail.NewEmail(fromName, domain.fromEmail)
	to := mail.NewEmail(toUsername, toAddress)

	// set the text content
	textContent = strings.Replace(textContent, "$micro_verification_link", "https://user.m3o.com?token="+token+"&redirectUrl="+url.QueryEscape(redirctUrl)+"&failureRedirectUrl="+url.QueryEscape(failureRedirectUrl), -1)
	message := mail.NewSingleEmail(from, subject, to, textContent, "")

	// send the email
	client := sendgrid.NewSendClient(domain.sengridKey)
	response, err := client.Send(message)
	logger.Info(response)

	return err
}

func (domain *Domain) SavePasswordResetCode(ctx context.Context, userID, code string) (*passwordResetCode, error) {
	pwcode := passwordResetCode{
		ID:      userID + "-" + code,
		Expires: time.Now().Add(24 * time.Hour),
		UserID:  userID,
		Code:    code,
	}

	s := &_struct.Struct{}
	jso, _ := json.Marshal(pwcode)
	err := s.UnmarshalJSON(jso)
	if err != nil {
		return nil, err
	}

	_, err = domain.db.Create(ctx, &db.CreateRequest{
		Table:  "password-reset-codes",
		Record: s,
	})

	return &pwcode, err
}

func (domain *Domain) DeletePasswordResetCode(ctx context.Context, userId, code string) error {
	_, err := domain.db.Delete(ctx, &db.DeleteRequest{
		Table: "password-reset-codes",
		Id:    userId + "-" + code,
	})
	return err
}

// ReadToken returns the user id
func (domain *Domain) ReadPasswordResetCode(ctx context.Context, userId, code string) (*passwordResetCode, error) {
	// generate the key
	id := userId + "-" + code

	if id == "" {
		return nil, errors.New("password reset code id is empty")
	}
	token := &passwordResetCode{}

	rsp, err := domain.db.Read(ctx, &db.ReadRequest{
		Table: "password-reset-codes",
		Query: fmt.Sprintf("id == '%v'", id),
	})
	if err != nil {
		return nil, err
	}
	if len(rsp.Records) == 0 {
		return nil, errors.New("password reset code not found")
	}
	m, _ := rsp.Records[0].MarshalJSON()
	json.Unmarshal(m, token)

	// check the expiry
	if token.Expires.Before(time.Now()) {
		return nil, errors.New("password reset code expired")
	}

	return token, nil
}

func (domain *Domain) SendPasswordResetEmail(ctx context.Context, userId, codeStr, fromName, toAddress, toUsername, subject, textContent string) error {
	if domain.sengridKey == "" {
		return fmt.Errorf("empty email api key")
	}

	from := mail.NewEmail(fromName, domain.fromEmail)
	to := mail.NewEmail(toUsername, toAddress)

	// save the password reset code
	pw, err := domain.SavePasswordResetCode(ctx, userId, codeStr)
	if err != nil {
		return err
	}

	// set the code in the text content
	textContent = strings.Replace(textContent, "$code", pw.Code, -1)
	message := mail.NewSingleEmail(from, subject, to, textContent, "")

	// send the email
	client := sendgrid.NewSendClient(domain.sengridKey)
	response, err := client.Send(message)

	// log the response
	logger.Info(response)

	return err
}

func (domain *Domain) CreateSession(ctx context.Context, sess *user.Session) error {
	if sess.Created == 0 {
		sess.Created = time.Now().Unix()
	}

	if sess.Expires == 0 {
		sess.Expires = time.Now().Add(time.Hour * 24 * 7).Unix()
	}

	s := &_struct.Struct{}
	jso, _ := json.Marshal(sess)
	err := s.UnmarshalJSON(jso)
	if err != nil {
		return err
	}
	_, err = domain.db.Create(ctx, &db.CreateRequest{
		Table:  "sessions",
		Record: s,
	})
	return err
}

func (domain *Domain) DeleteSession(ctx context.Context, id string) error {
	_, err := domain.db.Delete(ctx, &db.DeleteRequest{
		Table: "sessions",
		Id:    id,
	})
	return err
}

// ReadToken returns the user id
func (domain *Domain) ReadToken(ctx context.Context, userId, token string) (string, error) {
	id := userId + "-" + token

	if token == "" {
		return "", errors.New("token id empty")
	}

	tk := &verificationToken{}

	rsp, err := domain.db.Read(ctx, &db.ReadRequest{
		Table: "tokens",
		Query: fmt.Sprintf("id == '%v'", id),
	})
	if err != nil {
		return "", err
	}

	if len(rsp.Records) == 0 {
		return "", errors.New("token not found")
	}

	m, _ := rsp.Records[0].MarshalJSON()
	json.Unmarshal(m, tk)

	return tk.UserID, nil
}

// CreateToken returns the created and saved token
func (domain *Domain) CreateToken(ctx context.Context, userId, token string) (string, error) {
	s := &_struct.Struct{}
	jso, _ := json.Marshal(verificationToken{
		ID:     userId + "-" + token,
		UserID: userId,
		Token:  token,
	})

	err := s.UnmarshalJSON(jso)
	if err != nil {
		return "", err
	}

	_, err = domain.db.Create(ctx, &db.CreateRequest{
		Table:  "tokens",
		Record: s,
	})

	return token, err
}

func (domain *Domain) ReadSession(ctx context.Context, id string) (*user.Session, error) {
	sess := &user.Session{}
	if len(id) == 0 {
		return nil, fmt.Errorf("no id provided")
	}
	q := fmt.Sprintf("id == '%v'", id)
	logger.Infof("Running query: %v", q)

	rsp, err := domain.db.Read(ctx, &db.ReadRequest{
		Table: "sessions",
		Query: q,
	})
	if err != nil {
		return nil, err
	}
	if len(rsp.Records) == 0 {
		return nil, ErrNotFound
	}
	m, _ := rsp.Records[0].MarshalJSON()
	json.Unmarshal(m, sess)
	return sess, nil
}

func (domain *Domain) Create(ctx context.Context, user *user.Account, salt string, password string) error {
	user.Created = time.Now().Unix()
	user.Updated = time.Now().Unix()

	s := &_struct.Struct{}
	jso, _ := json.Marshal(user)
	err := s.UnmarshalJSON(jso)
	if err != nil {
		return err
	}
	_, err = domain.db.Create(ctx, &db.CreateRequest{
		Table:  "users",
		Record: s,
	})
	if err != nil {
		return err
	}

	pass := pw{
		ID:       user.Id,
		Password: password,
		Salt:     salt,
	}
	s = &_struct.Struct{}
	jso, _ = json.Marshal(pass)
	err = s.UnmarshalJSON(jso)
	if err != nil {
		return err
	}
	_, err = domain.db.Create(ctx, &db.CreateRequest{
		Table:  "passwords",
		Record: s,
	})

	return err
}

func (domain *Domain) Delete(ctx context.Context, id string) error {
	_, err := domain.db.Delete(ctx, &db.DeleteRequest{
		Table: "users",
		Id:    id,
	})
	return err
}

func (domain *Domain) Update(ctx context.Context, user *user.Account) error {
	user.Updated = time.Now().Unix()

	s := &_struct.Struct{}
	jso, _ := json.Marshal(user)
	err := s.UnmarshalJSON(jso)
	if err != nil {
		return err
	}
	_, err = domain.db.Update(ctx, &db.UpdateRequest{
		Table:  "users",
		Record: s,
	})
	return err
}

func (domain *Domain) Read(ctx context.Context, userId string) (*user.Account, error) {
	user := &user.Account{}
	if len(userId) == 0 {
		return nil, fmt.Errorf("no id provided")
	}
	q := fmt.Sprintf("id == '%v'", userId)
	logger.Infof("Running query: %v", q)
	rsp, err := domain.db.Read(ctx, &db.ReadRequest{
		Table: "users",
		Query: q,
	})
	if err != nil {
		return nil, err
	}
	if len(rsp.Records) == 0 {
		return nil, ErrNotFound
	}
	m, _ := rsp.Records[0].MarshalJSON()
	json.Unmarshal(m, user)
	return user, nil
}

func (domain *Domain) Search(ctx context.Context, username, email string) ([]*user.Account, error) {
	var query string
	if len(username) > 0 {
		query = fmt.Sprintf("username == '%v'", username)
	} else if len(email) > 0 {
		query = fmt.Sprintf("email == '%v'", email)
	} else {
		return nil, errors.New("username and email cannot be blank")
	}

	usr := &user.Account{}

	rsp, err := domain.db.Read(ctx, &db.ReadRequest{
		Table: "users",
		Query: query,
	})
	if err != nil {
		return nil, err
	}
	if len(rsp.Records) == 0 {
		return nil, ErrNotFound
	}
	m, _ := rsp.Records[0].MarshalJSON()
	json.Unmarshal(m, usr)
	return []*user.Account{usr}, nil
}

func (domain *Domain) UpdatePassword(ctx context.Context, id string, salt string, password string) error {
	pass := pw{
		ID:       id,
		Password: password,
		Salt:     salt,
	}
	s := &_struct.Struct{}
	jso, _ := json.Marshal(pass)
	err := s.UnmarshalJSON(jso)
	if err != nil {
		return err
	}
	_, err = domain.db.Update(ctx, &db.UpdateRequest{
		Table:  "passwords",
		Record: s,
	})
	return err
}

func (domain *Domain) SaltAndPassword(ctx context.Context, userId string) (string, string, error) {
	password := &pw{}

	rsp, err := domain.db.Read(ctx, &db.ReadRequest{
		Table: "passwords",
		Query: fmt.Sprintf("id == '%v'", userId),
	})
	if err != nil {
		return "", "", err
	}
	if len(rsp.Records) == 0 {
		return "", "", ErrNotFound
	}
	m, _ := rsp.Records[0].MarshalJSON()
	json.Unmarshal(m, password)
	return password.Salt, password.Password, nil
}

func (domain *Domain) List(ctx context.Context, o, l int32) ([]*user.Account, error) {
	var limit int32 = 25
	var offset int32 = 0
	if l > 0 {
		limit = l
	}
	if o > 0 {
		offset = o
	}
	rsp, err := domain.db.Read(ctx, &db.ReadRequest{
		Table:  "users",
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	if len(rsp.Records) == 0 {
		return nil, ErrNotFound
	}
	ret := make([]*user.Account, len(rsp.Records))
	for i, v := range rsp.Records {
		m, _ := v.MarshalJSON()
		var user user.Account
		json.Unmarshal(m, &user)
		ret[i] = &user
	}
	return ret, nil
}

func (domain *Domain) CacheToken(ctx context.Context, token, id, email string, ttl int) error {
	obj := &tokenObject{
		Id:    id,
		Email: email,
	}

	expires := time.Now().Add(time.Duration(ttl) * time.Second)

	err := cache.Context(ctx).Set(token, obj, expires)

	return err
}

func (domain *Domain) SendMLE(fromName, toAddress, toUsername, subject, textContent, token, address, endpoint string) error {
	if domain.sengridKey == "" {
		return fmt.Errorf("empty email api key")
	}
	from := mail.NewEmail(fromName, "support@m3o.com")
	to := mail.NewEmail(toUsername, toAddress)
	textContent = strings.Replace(textContent, "$micro_verification_link", fmt.Sprint("https://", path.Join(address, endpoint, token)), -1)
	message := mail.NewSingleEmail(from, subject, to, textContent, "")
	client := sendgrid.NewSendClient(domain.sengridKey)
	response, err := client.Send(message)
	logger.Info(response)

	return err
}

func (domain *Domain) CacheReadToken(ctx context.Context, token string) (string, string, error) {

	if token == "" {
		return "", "", errors.New("token empty")
	}

	var obj tokenObject

	expires, err := cache.Context(ctx).Get(token, obj)

	if err == cache.ErrNotFound {
		return "", "", errors.New("token not found")
	} else if time.Until(expires).Seconds() < 0 {
		return "", "", errors.New("token expired")
	} else if err != nil {
		return "", "", microerr.InternalServerError("CacheReadToken", err.Error())
	}

	return obj.Id, obj.Email, nil
}

type tokenObject struct {
	Id    string
	Email string
}
