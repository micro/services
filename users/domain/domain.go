package domain

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

type Domain struct {
	users     model.Model
	sessions  model.Model
	passwords model.Model

	nameIndex  model.Index
	emailIndex model.Index
	idIndex    model.Index
}

func New() *Domain {
	nameIndex := model.ByEquality("username")
	nameIndex.Unique = true
	nameIndex.Order.Type = model.OrderTypeUnordered

	emailIndex := model.ByEquality("email")
	emailIndex.Unique = true
	emailIndex.Order.Type = model.OrderTypeUnordered

	// @todo there should be a better way to get the default index from model
	// than recreating the options here
	idIndex := model.ByEquality("id")
	idIndex.Order.Type = model.OrderTypeUnordered

	return &Domain{
		users:      model.New(store.DefaultStore, "users", model.Indexes(nameIndex, emailIndex), nil),
		sessions:   model.New(store.DefaultStore, "sessions", nil, nil),
		passwords:  model.New(store.DefaultStore, "passwords", nil, nil),
		nameIndex:  nameIndex,
		emailIndex: emailIndex,
		idIndex:    idIndex,
	}
}

func (domain *Domain) CreateSession(sess *user.Session) error {
	if sess.Created == 0 {
		sess.Created = time.Now().Unix()
	}

	if sess.Expires == 0 {
		sess.Expires = time.Now().Add(time.Hour * 24 * 7).Unix()
	}

	return domain.sessions.Save(sess)
}

func (domain *Domain) DeleteSession(id string) error {
	return domain.sessions.Delete(domain.idIndex.ToQuery(id))
}

func (domain *Domain) ReadSession(id string) (*user.Session, error) {
	sess := &user.Session{}
	// @todo there should be a Read in the model to get rid of this pattern
	return sess, domain.sessions.Read(domain.idIndex.ToQuery(id), &sess)
}

func (domain *Domain) Create(user *user.User, salt string, password string) error {
	user.Created = time.Now().Unix()
	user.Updated = time.Now().Unix()
	err := domain.users.Save(user)
	if err != nil {
		return err
	}
	return domain.passwords.Save(pw{
		ID:       user.Id,
		Password: password,
		Salt:     salt,
	})
}

func (domain *Domain) Delete(id string) error {
	return domain.users.Delete(domain.idIndex.ToQuery(id))
}

func (domain *Domain) Update(user *user.User) error {
	user.Updated = time.Now().Unix()
	return domain.users.Save(user)
}

func (domain *Domain) Read(id string) (*user.User, error) {
	user := &user.User{}
	return user, domain.users.Read(domain.idIndex.ToQuery(id), user)
}

func (domain *Domain) Search(username, email string, limit, offset int64) ([]*user.User, error) {
	var query model.Query
	if len(username) > 0 {
		query = domain.nameIndex.ToQuery(username)
	} else if len(email) > 0 {
		query = domain.emailIndex.ToQuery(email)
	} else {
		return nil, errors.New("username and email cannot be blank")
	}

	users := []*user.User{}
	return users, domain.users.List(query, &users)
}

func (domain *Domain) UpdatePassword(id string, salt string, password string) error {
	return domain.passwords.Save(pw{
		ID:       id,
		Password: password,
		Salt:     salt,
	})
}

func (domain *Domain) SaltAndPassword(username, email string) (string, string, error) {
	var query model.Query
	if len(username) > 0 {
		query = domain.nameIndex.ToQuery(username)
	} else if len(email) > 0 {
		query = domain.emailIndex.ToQuery(email)
	} else {
		return "", "", errors.New("username and email cannot be blank")
	}

	user := &user.User{}
	err := domain.users.Read(query, &user)
	if err != nil {
		return "", "", err
	}

	query = model.Equals("id", user.Id)
	query.Order.Type = model.OrderTypeUnordered

	password := &pw{}
	err = domain.passwords.Read(query, password)
	if err != nil {
		return "", "", err
	}
	return password.Salt, password.Password, nil
}
