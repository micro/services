package dao

import (
	"errors"
	"time"

	"github.com/micro/dev/model"
	"github.com/micro/micro/v3/service/store"
	user "github.com/micro/services/users/proto"
)

type pw struct {
	ID       string `json:"id"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

type Dao struct {
	users     model.Table
	sessions  model.Table
	passwords model.Table
}

func New() *Dao {
	nameIndex := model.ByEquality("username")
	nameIndex.Unique = true
	emailIndex := model.ByEquality("email")
	emailIndex.Unique = true

	return &Dao{
		users:     model.NewTable(store.DefaultStore, "users", model.Indexes(nameIndex, emailIndex), nil),
		sessions:  model.NewTable(store.DefaultStore, "sessions", nil, nil),
		passwords: model.NewTable(store.DefaultStore, "passwords", nil, nil),
	}
}

func (dao *Dao) CreateSession(sess *user.Session) error {
	if sess.Created == 0 {
		sess.Created = time.Now().Unix()
	}

	if sess.Expires == 0 {
		sess.Expires = time.Now().Add(time.Hour * 24 * 7).Unix()
	}

	return dao.sessions.Save(sess)
}

func (dao *Dao) DeleteSession(id string) error {
	return dao.sessions.Delete(model.Equals("id", id))
}

func (dao *Dao) ReadSession(id string) (*user.Session, error) {
	sess := &user.Session{}
	// @todo there should be a Read in the model to get rid of this pattern
	return sess, dao.sessions.Read(model.Equals("id", id), &sess)
}

func (dao *Dao) Create(user *user.User, salt string, password string) error {
	user.Created = time.Now().Unix()
	user.Updated = time.Now().Unix()
	err := dao.users.Save(user)
	if err != nil {
		return err
	}
	return dao.passwords.Save(pw{
		ID:       user.Id,
		Password: password,
		Salt:     salt,
	})
}

func (dao *Dao) Delete(id string) error {
	return dao.users.Delete(model.Equals("id", id))
}

func (dao *Dao) Update(user *user.User) error {
	user.Updated = time.Now().Unix()
	return dao.users.Save(user)
}

func (dao *Dao) Read(id string) (*user.User, error) {
	user := &user.User{}
	q := model.Equals("id", id)
	q.Order.Type = model.OrderTypeUnordered
	return user, dao.users.Read(q, user)
}

func (dao *Dao) Search(username, email string, limit, offset int64) ([]*user.User, error) {
	var query model.Query
	if len(username) > 0 {
		query = model.Equals("username", username)
	} else if len(email) > 0 {
		query = model.Equals("email", email)
	} else {
		return nil, errors.New("username and email cannot be blank")
	}

	users := []*user.User{}
	return users, dao.users.List(query, &users)
}

func (dao *Dao) UpdatePassword(id string, salt string, password string) error {
	return dao.passwords.Save(pw{
		ID:       id,
		Password: password,
		Salt:     salt,
	})
}

func (dao *Dao) SaltAndPassword(username, email string) (string, string, error) {
	var query model.Query
	if len(username) > 0 {
		query = model.Equals("name", username)
	} else if len(email) > 0 {
		query = model.Equals("email", email)
	} else {
		return "", "", errors.New("username and email cannot be blank")
	}

	user := &user.User{}
	err := dao.users.Read(query, &user)
	if err != nil {
		return "", "", err
	}

	query = model.Equals("id", user.Id)
	query.Order.Type = model.OrderTypeUnordered

	password := &pw{}
	err = dao.passwords.Read(query, password)
	if err != nil {
		return "", "", err
	}
	return password.Salt, password.Password, nil
}
