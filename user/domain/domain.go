package domain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/micro/micro/v3/service/config"
	microerr "github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"

	"github.com/micro/services/pkg/cache"
	user "github.com/micro/services/user/proto"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var (
	ErrNotFound = errors.New("not found")
)

type pw struct {
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

type verificationToken struct {
	UserID string `json:"userId"`
	Token  string `json:"token"`
}

type passwordResetCode struct {
	Expires time.Time `json:"expires"`
	UserID  string    `json:"userId"`
	Code    string    `json:"code"`
}

type Domain struct {
	store      store.Store
	sengridKey string
	fromEmail  string
}

var (
	// TODO: use the config to drive this value
	defaultSender = "noreply@email.m3ocontent.com"
)

func New(st store.Store) *Domain {
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
		store:      st,
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

func (domain *Domain) SavePasswordResetCode(_ context.Context, userId, code string) (*passwordResetCode, error) {
	pwcode := passwordResetCode{
		Expires: time.Now().Add(24 * time.Hour),
		UserID:  userId,
		Code:    code,
	}

	val, err := json.Marshal(pwcode)
	if err != nil {
		return nil, err
	}

	record := store.NewRecord(generatePasswordResetCodeStoreKey(userId, code), val)
	err = domain.store.Write(record)

	return &pwcode, err
}

func (domain *Domain) DeletePasswordResetCode(_ context.Context, userId, code string) error {
	return domain.store.Delete(generatePasswordResetCodeStoreKey(userId, code))
}

// ReadPasswordResetCode returns the user reset code
func (domain *Domain) ReadPasswordResetCode(_ context.Context, userId, code string) (*passwordResetCode, error) {
	records, err := domain.store.Read(generatePasswordResetCodeStoreKey(userId, code))
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, errors.New("password reset code not found")
	}

	resetCode := &passwordResetCode{}
	if err := json.Unmarshal(records[0].Value, resetCode); err != nil {
		return nil, err
	}

	// check the expiry
	if resetCode.Expires.Before(time.Now()) {
		return nil, errors.New("password reset code expired")
	}

	return resetCode, nil
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

func (domain *Domain) CreateSession(_ context.Context, sess *user.Session) error {
	if sess.Created == 0 {
		sess.Created = time.Now().Unix()
	}

	if sess.Expires == 0 {
		sess.Expires = time.Now().Add(time.Hour * 24 * 7).Unix()
	}

	val, err := json.Marshal(sess)
	if err != nil {
		return err
	}
	record := &store.Record{
		Key:   generateSessionStoreKey(sess.Id),
		Value: val,
	}

	return domain.store.Write(record)
}

func (domain *Domain) DeleteSession(_ context.Context, id string) error {
	return domain.store.Delete(generateSessionStoreKey(id))
}

// ReadToken returns the user id
func (domain *Domain) ReadToken(_ context.Context, userId, token string) (string, error) {
	if token == "" {
		return "", errors.New("token id empty")
	}

	records, err := domain.store.Read(generateVerificationsTokenStoreKey(userId, token))
	if err != nil {
		return "", err
	}

	if len(records) == 0 {
		return "", errors.New("token not found")
	}

	tk := &verificationToken{}
	err = json.Unmarshal(records[0].Value, tk)
	if err != nil {
		return "", err
	}
	return tk.UserID, nil
}

// CreateToken returns the created and saved token
func (domain *Domain) CreateToken(ctx context.Context, userId, token string) (string, error) {
	tk, err := json.Marshal(verificationToken{
		UserID: userId,
		Token:  token,
	})

	if err != nil {
		return "", err
	}

	record := &store.Record{
		Key:   generateVerificationsTokenStoreKey(userId, token),
		Value: tk,
	}
	err = domain.store.Write(record)
	if err != nil {
		return "", err
	}

	return token, err
}

func (domain *Domain) ReadSession(ctx context.Context, id string) (*user.Session, error) {
	records, err := domain.store.Read(generateSessionStoreKey(id))
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, ErrNotFound
	}

	sess := &user.Session{}
	err = json.Unmarshal(records[0].Value, sess)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

// batchWrite write multiple records in batches
func (domain *Domain) batchWrite(records []*store.Record) error {
	if len(records) == 0 {
		return nil
	}

	wg := sync.WaitGroup{}
	errs := make([]string, 0)
	for _, v := range records {
		wg.Add(1)
		go func(r *store.Record) {
			defer wg.Done()
			if err := domain.store.Write(r); err != nil {
				errs = append(errs, err.Error())
			}
		}(v)
	}
	wg.Wait()

	if len(errs) != 0 {
		return errors.New(strings.Join(errs, ";"))
	}

	return nil
}

func (domain *Domain) Create(ctx context.Context, user *user.Account, salt string, password string) error {
	user.Created = time.Now().Unix()
	user.Updated = time.Now().Unix()

	// user account record
	accountVal, err := json.Marshal(user)
	if err != nil {
		return err
	}

	// password record
	passwordVal, err := json.Marshal(pw{
		Password: password,
		Salt:     salt,
	})
	if err != nil {
		return err
	}

	records := []*store.Record{
		{Key: generateAccountStoreKey(user.Id), Value: accountVal},
		{Key: generateAccountUsernameStoreKey(user.Username), Value: accountVal},
		{Key: generateAccountEmailStoreKey(user.Email), Value: accountVal},
		{Key: generatePasswordStoreKey(user.Id), Value: passwordVal},
	}

	return domain.batchWrite(records)
}

// batchDelete deletes the keys in batches
func (domain *Domain) batchDelete(keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	wg := sync.WaitGroup{}
	errs := make([]string, 0)
	for _, key := range keys {
		wg.Add(1)
		go func(keyToDel string) {
			defer wg.Done()
			if err := domain.store.Delete(keyToDel); err != nil {
				errs = append(errs, err.Error())
			}
		}(key)

	}
	wg.Wait()

	if len(errs) != 0 {
		return errors.New(strings.Join(errs, ";"))
	}

	return nil
}

func (domain *Domain) Delete(ctx context.Context, userId string) error {
	account, err := domain.Read(ctx, userId)
	if err != nil {
		return err
	}

	keys := []string{
		generateAccountStoreKey(userId),
		generateAccountEmailStoreKey(account.Email),
		generateAccountUsernameStoreKey(account.Username),
	}

	return domain.batchDelete(keys)
}

func (domain *Domain) Update(ctx context.Context, user *user.Account) error {
	// get old information of the user
	old, err := domain.Read(ctx, user.Id)
	if err != nil {
		return err
	}

	keysToDelete := make([]string, 0)
	if old.Email != user.Email {
		keysToDelete = append(keysToDelete, generateAccountEmailStoreKey(old.Email))
	}

	if old.Username != user.Username {
		keysToDelete = append(keysToDelete, generateAccountUsernameStoreKey(old.Username))
	}

	// update user
	user.Created = old.Created
	user.Updated = time.Now().Unix()
	val, err := json.Marshal(user)
	if err != nil {
		return err
	}

	records := []*store.Record{
		{Key: generateAccountStoreKey(user.Id), Value: val},
		{Key: generateAccountUsernameStoreKey(user.Username), Value: val},
		{Key: generateAccountEmailStoreKey(user.Email), Value: val},
	}

	// update
	if err := domain.batchWrite(records); err != nil {
		return err
	}

	// delete
	if err := domain.batchDelete(keysToDelete); err != nil {
		return err
	}

	return nil
}

// readUserByKey read user account in store by key
func (domain *Domain) readUserByKey(_ context.Context, key string) (*user.Account, error) {
	var result = &user.Account{}
	records, err := domain.store.Read(key)
	if err != nil {
		return result, err
	}

	if len(records) == 0 {
		return result, ErrNotFound
	}

	err = json.Unmarshal(records[0].Value, result)
	return result, err
}

func (domain *Domain) Read(ctx context.Context, userId string) (*user.Account, error) {
	return domain.readUserByKey(ctx, generateAccountStoreKey(userId))
}

func (domain *Domain) SearchByUsername(ctx context.Context, username string) (*user.Account, error) {
	return domain.readUserByKey(ctx, generateAccountUsernameStoreKey(username))
}

func (domain *Domain) SearchByEmail(ctx context.Context, email string) (*user.Account, error) {
	return domain.readUserByKey(ctx, generateAccountEmailStoreKey(email))
}

func (domain *Domain) Search(ctx context.Context, username, email string) ([]*user.Account, error) {
	var account = &user.Account{}
	var err error

	switch {
	case username != "":
		account, err = domain.SearchByUsername(ctx, username)
	case email != "":
		account, err = domain.SearchByEmail(ctx, email)
	}

	if err != nil {
		return []*user.Account{}, err
	}

	return []*user.Account{account}, nil
}

func (domain *Domain) UpdatePassword(_ context.Context, userId string, salt string, password string) error {
	val, err := json.Marshal(pw{
		Password: password,
		Salt:     salt,
	})

	if err != nil {
		return err
	}

	record := &store.Record{
		Key:   generatePasswordStoreKey(userId),
		Value: val,
	}

	return domain.store.Write(record)
}

func (domain *Domain) SaltAndPassword(_ context.Context, userId string) (string, string, error) {
	records, err := domain.store.Read(generatePasswordStoreKey(userId))
	if err != nil {
		return "", "", err
	}
	if len(records) == 0 {
		return "", "", ErrNotFound
	}

	password := &pw{}
	if err := json.Unmarshal(records[0].Value, password); err != nil {
		return "", "", err
	}

	return password.Salt, password.Password, nil
}

func (domain *Domain) List(_ context.Context, o, l uint32) (result []*user.Account, err error) {
	records, err := domain.store.Read(generateAccountStoreKey(""),
		store.ReadPrefix(),
		store.ReadLimit(uint(l)),
		store.ReadLimit(uint(o)))

	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return result, ErrNotFound
	}

	ret := make([]*user.Account, len(records))

	for i, v := range records {
		account := user.Account{}
		json.Unmarshal(v.Value, &account)
		ret[i] = &account
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

type tokenObject struct {
	Id    string
	Email string
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
