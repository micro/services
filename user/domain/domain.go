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

func generatePasswordResetCodeStoreKey(userId, code string) string {
	return fmt.Sprintf("user/password-reset-codes/%s-%s", userId, code)
}

func generateSessionStoreKey(sessionId string) string {
	return fmt.Sprintf("user/session/%s", sessionId)
}

func generateVerificationsTokenStoreKey(userId, token string) string {
	return fmt.Sprintf("user/verifycation-token/%s-%s", userId, token)
}

func generateAccountStoreKey(userId string) string {
	return fmt.Sprintf("user/account/%s", userId)
}

func generatePasswordStoreKey(userId string) string {
	return fmt.Sprintf("user/password/%s", userId)
}

type Domain struct {
	store store.Store
	//db         db.DbService
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
	key := generatePasswordResetCodeStoreKey(userId, code)

	records, err := domain.store.Read(key)
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

func (domain *Domain) Create(ctx context.Context, user *user.Account, salt string, password string) error {
	user.Created = time.Now().Unix()
	user.Updated = time.Now().Unix()

	records := make([]*store.Record, 2)

	// user account record
	val, err := json.Marshal(user)
	if err != nil {
		return err
	}
	records = append(records, &store.Record{
		Key:   generateAccountStoreKey(user.Id),
		Value: val,
	})

	// password record
	val, err = json.Marshal(pw{
		Password: password,
		Salt:     salt,
	})
	if err != nil {
		return err
	}
	records = append(records, &store.Record{
		Key:   generatePasswordStoreKey(user.Id),
		Value: val,
	})

	wg := sync.WaitGroup{}
	errs := make([]error, 0)
	for _, v := range records {
		wg.Add(1)
		go func(r *store.Record) {
			defer wg.Done()
			if err := store.Write(r); err != nil {
				errs = append(errs, err)
			}
		}(v)
	}
	wg.Wait()

	if len(errs) != 0 {
		return errs[0]
	}

	return nil
}

func (domain *Domain) Delete(_ context.Context, id string) error {
	return domain.store.Delete(generateAccountStoreKey(id))
}

func (domain *Domain) Update(ctx context.Context, user *user.Account) error {
	user.Updated = time.Now().Unix()
	val, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return domain.store.Write(&store.Record{
		Key:   generateAccountStoreKey(user.Id),
		Value: val,
	})
}

func (domain *Domain) Read(_ context.Context, userId string) (*user.Account, error) {
	records, err := domain.store.Read(generateAccountStoreKey(userId))
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, ErrNotFound
	}
	result := &user.Account{}
	err = json.Unmarshal(records[0].Value, result)
	return result, err
}

// TODO: search
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

func (domain *Domain) UpdatePassword(ctx context.Context, userId string, salt string, password string) error {
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

func (domain *Domain) List(_ context.Context, o, l uint32) ([]*user.Account, error) {
	records, err := store.Read("user/account/", store.ReadLimit(uint(l)), store.ReadLimit(uint(o)))

	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, ErrNotFound
	}

	ret := make([]*user.Account, len(records))

	for i, v := range records {
		account := &user.Account{}
		json.Unmarshal(v.Value, account)
		ret[i] = account
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
		return "", "", ErrNotFound
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
