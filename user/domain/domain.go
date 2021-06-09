package domain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	_struct "github.com/golang/protobuf/ptypes/struct"
	db "github.com/micro/services/db/proto"
	user "github.com/micro/services/user/proto"
)

type pw struct {
	ID       string `json:"id"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

type Domain struct {
	db db.DbService
}

func New(db db.DbService) *Domain {
	return &Domain{
		db: db,
	}
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

func (domain *Domain) ReadSession(ctx context.Context, id string) (*user.Session, error) {
	sess := &user.Session{}

	rsp, err := domain.db.Read(ctx, &db.ReadRequest{
		Table: "sessions",
		Query: fmt.Sprintf("id == '%v'", id),
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

func (domain *Domain) Read(ctx context.Context, id string) (*user.Account, error) {
	user := &user.Account{}

	rsp, err := domain.db.Read(ctx, &db.ReadRequest{
		Table: "users",
		Query: fmt.Sprintf("id == '%v'", id),
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

func (domain *Domain) Search(ctx context.Context, username, email string, limit, offset int64) ([]*user.Account, error) {
	var query string
	if len(username) > 0 {
		query = fmt.Sprint("userName == '%v'")
	} else if len(email) > 0 {
		query = fmt.Sprintf("email == '%v'")
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

func (domain *Domain) SaltAndPassword(ctx context.Context, username, email string) (string, string, error) {
	var query string
	if len(username) > 0 {
		query = fmt.Sprint("userName == '%v'")
	} else if len(email) > 0 {
		query = fmt.Sprintf("email == '%v'")
	} else {
		return "", "", errors.New("username and email cannot be blank")
	}

	password := &pw{}

	rsp, err := domain.db.Read(ctx, &db.ReadRequest{
		Table: "passwords",
		Query: query,
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
