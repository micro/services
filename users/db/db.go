package db

import (
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/micro/dev/model"
	user "github.com/micro/services/users/proto"
)

var (
	Accounts  model.DB
	Sessions  model.DB
	Passwords model.DB
)

type pw struct {
	ID       string `json:"id"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

func CreateSession(sess *user.Session) error {
	if sess.Created == 0 {
		sess.Created = time.Now().Unix()
	}

	if sess.Expires == 0 {
		sess.Expires = time.Now().Add(time.Hour * 24 * 7).Unix()
	}

	return Sessions.Save(sess)
}

func DeleteSession(id string) error {
	return Sessions.Delete(model.Equals("id", id))
}

func ReadSession(id string) (*user.Session, error) {
	sess := []*user.Session{}
	// @todo there should be a Read in the model to get rid of this pattern
	err := Sessions.List(model.Equals("id", id), &sess)
	if err != nil {
		return nil, err
	}
	if len(sess) == 0 {
		return nil, errors.New("Not found")
	}
	return sess[0], nil
}

func Create(user *user.User, salt string, password string) error {
	user.Created = time.Now().Unix()
	user.Updated = time.Now().Unix()
	// @todo what to do with salt and password?
	err := Accounts.Save(user)
	if err != nil {
		return err
	}
	return Passwords.Save(pw{
		ID:       user.Id,
		Password: password,
		Salt:     salt,
	})
}

func Delete(id string) error {
	return Accounts.Delete(model.Equals("id", id))
}

func Update(user *user.User) error {
	user.Updated = time.Now().Unix()
	return Accounts.Save(user)
}

func Read(id string) (*user.User, error) {
	users := []*user.User{}
	// @todo there should be a Read in the model to get rid of this pattern
	err := Accounts.List(model.Equals("id", id), &users)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("Not found")
	}
	return users[0], nil
}

func Search(username, email string, limit, offset int64) ([]*user.User, error) {
	var query model.Query
	if len(username) > 0 {
		query = model.Equals("name", username)
	} else if len(email) > 0 {
		query = model.Equals("email", email)
	} else {
		return nil, errors.New("username and email cannot be blank")
	}

	users := []*user.User{}
	return users, Accounts.List(query, &users)
}

func UpdatePassword(id string, salt string, password string) error {
	return Passwords.Save(pw{
		ID:       id,
		Password: password,
		Salt:     salt,
	})
}

func SaltAndPassword(username, email string) (string, string, error) {
	var query model.Query
	if len(username) > 0 {
		query = model.Equals("name", username)
	} else if len(email) > 0 {
		query = model.Equals("email", email)
	} else {
		return "", "", errors.New("username and email cannot be blank")
	}

	users := []*user.User{}
	err := Accounts.List(query, &users)
	if err != nil {
		return "", "", err
	}
	if len(users) == 0 {
		return "", "", errors.New("not found")
	}

	query = model.Equals("id", users[0].Id)
	query.Order.Type = model.OrderTypeUnordered

	passwords := []pw{}
	err = Passwords.List(query, &passwords)
	if err != nil {
		return "", "", err
	}
	if len(passwords) == 0 {
		return "", "", errors.New("not found")
	}
	return passwords[0].Salt, passwords[0].Password, nil
}
