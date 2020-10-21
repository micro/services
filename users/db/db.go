package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	user "github.com/micro/services/users/proto"
)

var (
	Url      = "root:root@tcp(127.0.0.1:3306)/user"
	database string
	db       *sql.DB

	q = map[string]string{}

	accountQ = map[string]string{
		"delete": "DELETE from %s.%s where id = ?",
		"create": `INSERT into %s.%s (
				id, username, email, salt, password, created, updated) 
				values (?, ?, ?, ?, ?, ?, ?)`,
		"update":                 "UPDATE %s.%s set username = ?, email = ?, updated = ? where id = ?",
		"updatePassword":         "UPDATE %s.%s set salt = ?, password = ?, updated = ? where id = ?",
		"read":                   "SELECT id, username, email, salt, password, created, updated from %s.%s where id = ?",
		"list":                   "SELECT id, username, email, salt, password, created, updated from %s.%s limit ? offset ?",
		"searchUsername":         "SELECT id, username, email, salt, password, created, updated from %s.%s where username = ? limit ? offset ?",
		"searchEmail":            "SELECT id, username, email, salt, password, created, updated from %s.%s where email = ? limit ? offset ?",
		"searchUsernameAndEmail": "SELECT id, username, email, salt, password, created, updated from %s.%s where username = ? and email = ? limit ? offset ?",
	}

	sessionQ = map[string]string{
		"createSession": "INSERT into %s.%s (id, username, created, expires) values (?, ?, ?, ?)",
		"deleteSession": "DELETE from %s.%s where id = ?",
		"readSession":   "SELECT id, username, created, expires from %s.%s where id = ?",
	}

	st = map[string]*sql.Stmt{}
)

func Init() {
	var d *sql.DB
	var err error

	parts := strings.Split(Url, "/")
	if len(parts) != 2 {
		panic("Invalid database url")
	}

	if len(parts[1]) == 0 {
		panic("Invalid database name")
	}

	url := parts[0]
	database = parts[1]

	if d, err = sql.Open("mysql", url+"/"); err != nil {
		log.Fatal(err)
	}
	if _, err := d.Exec("CREATE DATABASE IF NOT EXISTS " + database); err != nil {
		log.Fatal(err)
	}
	d.Close()
	if d, err = sql.Open("mysql", Url); err != nil {
		log.Fatal(err)
	}
	if _, err = d.Exec(accountSchema); err != nil {
		log.Fatal(err)
	}
	if _, err = d.Exec(sessionSchema); err != nil {
		log.Fatal(err)
	}

	db = d

	for query, statement := range accountQ {
		prepared, err := db.Prepare(fmt.Sprintf(statement, database, "accounts"))
		if err != nil {
			log.Fatal(err)
		}
		st[query] = prepared
	}

	for query, statement := range sessionQ {
		prepared, err := db.Prepare(fmt.Sprintf(statement, database, "sessions"))
		if err != nil {
			log.Fatal(err)
		}
		st[query] = prepared
	}
}

func CreateSession(sess *user.Session) error {
	if sess.Created == 0 {
		sess.Created = time.Now().Unix()
	}

	if sess.Expires == 0 {
		sess.Expires = time.Now().Add(time.Hour * 24 * 7).Unix()
	}

	_, err := st["createSession"].Exec(sess.Id, sess.Username, sess.Created, sess.Expires)
	return err
}

func DeleteSession(id string) error {
	_, err := st["deleteSession"].Exec(id)
	return err
}

func ReadSession(id string) (*user.Session, error) {
	sess := &user.Session{}

	r := st["readSession"].QueryRow(id)
	if err := r.Scan(&sess.Id, &sess.Username, &sess.Created, &sess.Expires); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("not found")
		}
		return nil, err
	}

	return sess, nil
}

func Create(user *user.User, salt string, password string) error {
	user.Created = time.Now().Unix()
	user.Updated = time.Now().Unix()
	_, err := st["create"].Exec(user.Id, user.Username, user.Email, salt, password, user.Created, user.Updated)
	return err
}

func Delete(id string) error {
	_, err := st["delete"].Exec(id)
	return err
}

func Update(user *user.User) error {
	user.Updated = time.Now().Unix()
	_, err := st["update"].Exec(user.Username, user.Email, user.Updated, user.Id)
	return err
}

func Read(id string) (*user.User, error) {
	user := &user.User{}

	r := st["read"].QueryRow(id)
	var s, p string
	if err := r.Scan(&user.Id, &user.Username, &user.Email, &s, &p, &user.Created, &user.Updated); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("not found")
		}
		return nil, err
	}

	return user, nil
}

func Search(username, email string, limit, offset int64) ([]*user.User, error) {
	var r *sql.Rows
	var err error

	if len(username) > 0 && len(email) > 0 {
		r, err = st["searchUsernameAndEmail"].Query(username, email, limit, offset)
	} else if len(username) > 0 {
		r, err = st["searchUsername"].Query(username, limit, offset)
	} else if len(email) > 0 {
		r, err = st["searchEmail"].Query(email, limit, offset)
	} else {
		r, err = st["list"].Query(limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer r.Close()

	var users []*user.User

	for r.Next() {
		user := &user.User{}
		var s, p string
		if err := r.Scan(&user.Id, &user.Username, &user.Email, &s, &p, &user.Created, &user.Updated); err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("not found")
			}
			return nil, err
		}
		users = append(users, user)

	}
	if r.Err() != nil {
		return nil, err
	}

	return users, nil
}

func UpdatePassword(id string, salt string, password string) error {
	_, err := st["updatePassword"].Exec(salt, password, time.Now().Unix(), id)
	return err
}

func SaltAndPassword(username, email string) (string, string, error) {
	var r *sql.Rows
	var err error

	if len(username) > 0 && len(email) > 0 {
		r, err = st["searchUsernameAndEmail"].Query(username, email, 1, 0)
	} else if len(username) > 0 {
		r, err = st["searchUsername"].Query(username, 1, 0)
	} else if len(email) > 0 {
		r, err = st["searchEmail"].Query(email, 1, 0)
	} else {
		return "", "", errors.New("username and email cannot be blank")
	}

	if err != nil {
		return "", "", err
	}
	defer r.Close()

	if !r.Next() {
		return "", "", errors.New("not found")
	}

	var salt, pass string
	user := &user.User{}
	if err := r.Scan(&user.Id, &user.Username, &user.Email, &salt, &pass, &user.Created, &user.Updated); err != nil {
		if err == sql.ErrNoRows {
			return "", "", errors.New("not found")
		}
		return "", "", err
	}
	if r.Err() != nil {
		return "", "", err
	}

	return salt, pass, nil
}
