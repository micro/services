package domain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	_struct "github.com/golang/protobuf/ptypes/struct"
	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/logger"
	db "github.com/micro/services/db/proto"
	user "github.com/micro/services/user/proto"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type pw struct {
	ID       string `json:"id"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

type verificationToken struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`
}

type Domain struct {
	db         db.DbService
	sengridKey string
}

func New(db db.DbService) *Domain {
	var key string
	cfg, err := config.Get("micro.user.sendgrid.api_key")
	if err == nil {
		key = cfg.String("")
	}
	if len(key) == 0 {
		logger.Info("No email key found")
	} else {
		logger.Info("Email key found")
	}
	return &Domain{
		sengridKey: key,
		db:         db,
	}
}

func (domain *Domain) SendEmail(fromName, toAddress, toUsername, subject, textContent, token, redirctUrl, failureRedirectUrl string) error {
	if domain.sengridKey == "" {
		return fmt.Errorf("empty email api key")
	}
	from := mail.NewEmail(fromName, "support@m3o.com")
	to := mail.NewEmail(toUsername, toAddress)
	textContent = strings.Replace(textContent, "$micro_verification_link", "https://angry-cori-854281.netlify.app?token="+token+"&redirectUrl="+url.QueryEscape(redirctUrl)+"&failureRedirectUrl="+url.QueryEscape(failureRedirectUrl), -1)
	message := mail.NewSingleEmail(from, subject, to, textContent, "")
	client := sendgrid.NewSendClient(domain.sengridKey)
	response, err := client.Send(message)
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
func (domain *Domain) ReadToken(ctx context.Context, tokenId string) (string, error) {
	if tokenId == "" {
		return "", errors.New("token id empty")
	}
	token := &verificationToken{}

	rsp, err := domain.db.Read(ctx, &db.ReadRequest{
		Table: "tokens",
		Query: fmt.Sprintf("id == '%v'", tokenId),
	})
	if err != nil {
		return "", err
	}
	if len(rsp.Records) == 0 {
		return "", errors.New("token not found")
	}
	m, _ := rsp.Records[0].MarshalJSON()
	json.Unmarshal(m, token)
	return token.UserID, nil
}

// CreateToken returns the created and saved token
func (domain *Domain) CreateToken(ctx context.Context, userId string) (string, error) {
	s := &_struct.Struct{}
	tokenId := uuid.New().String()
	jso, _ := json.Marshal(verificationToken{
		ID:     tokenId,
		UserID: userId,
	})
	err := s.UnmarshalJSON(jso)
	if err != nil {
		return "", err
	}
	_, err = domain.db.Create(ctx, &db.CreateRequest{
		Table:  "tokens",
		Record: s,
	})
	return tokenId, err
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
		return nil, errors.New("not found")
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
		return nil, errors.New("not found")
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
		return nil, errors.New("not found")
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
		return "", "", errors.New("not found")
	}
	m, _ := rsp.Records[0].MarshalJSON()
	json.Unmarshal(m, password)
	return password.Salt, password.Password, nil
}
