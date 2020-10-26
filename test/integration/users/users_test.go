// +build integration

package signup

import (
	"encoding/json"
	"errors"
	"math/rand"

	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/micro/micro/v3/test"
)

const (
	retryCount = 1
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func setupUsersTests(serv test.Server, t *test.T) {
	envToConfigKey := map[string][]string{}

	if err := test.Try("Set up config values", t, func() ([]byte, error) {
		for envKey, configKeys := range envToConfigKey {
			val := os.Getenv(envKey)
			if len(val) == 0 {
				t.Fatalf("'%v' flag is missing", envKey)
			}
			for _, configKey := range configKeys {
				outp, err := serv.Command().Exec("config", "set", configKey, val)
				if err != nil {
					return outp, err
				}
			}
		}
		return serv.Command().Exec("config", "set", "micro.billing.max_included_services", "3")
	}, 10*time.Second); err != nil {
		t.Fatal(err)
		return
	}

	services := []struct {
		envVar string
		deflt  string
	}{
		{envVar: "POSTS_SVC", deflt: "../../../users"},
	}

	for _, v := range services {
		outp, err := serv.Command().Exec("run", v.deflt)
		if err != nil {
			t.Fatal(string(outp))
			return
		}
	}

	if err := test.Try("Find posts and tags", t, func() ([]byte, error) {
		outp, err := serv.Command().Exec("services")
		if err != nil {
			return outp, err
		}
		list := []string{"users"}
		logOutp := []byte{}
		fail := false
		for _, s := range list {
			if !strings.Contains(string(outp), s) {
				o, _ := serv.Command().Exec("logs", s)
				logOutp = append(logOutp, o...)
				fail = true
			}
		}
		if fail {
			return append(outp, logOutp...), errors.New("Can't find required services in list")
		}
		return outp, err
	}, 180*time.Second); err != nil {
		return
	}

	// setup rules

	// Adjust rules before we signup into a non admin account
	outp, err := serv.Command().Exec("auth", "create", "rule", "--access=granted", "--scope=''", "--resource=\"service:users:*\"", "users")
	if err != nil {
		t.Fatalf("Error setting up rules: %v", outp)
		return
	}

	// copy the config with the admin logged in so we can use it for reading logs
	// we dont want to have an open access rule for logs as it's not how it works in live
	confPath := serv.Command().Config
	outp, err = exec.Command("cp", "-rf", confPath, confPath+".admin").CombinedOutput()
	if err != nil {
		t.Fatalf("Error copying config: %v", outp)
		return
	}
}

func TestUsersService(t *testing.T) {
	test.TrySuite(t, testUsers, retryCount)
}

func testUsers(t *test.T) {
	t.Parallel()

	serv := test.NewServer(t, test.WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	setupUsersTests(serv, t)

	cmd := serv.Command()

	email := "test@gmail.com"
	password := "testPassw"
	username := "john"
	id := "7"

	if err := test.Try("Save user", t, func() ([]byte, error) {
		// Attention! The content must be unquoted, don't add quotes.
		outp, err := cmd.Exec("users", "create", "--id="+id, "--email="+email, "--password="+password, "--username=john")
		if err != nil {
			outp1, _ := cmd.Exec("logs", "users")
			return append(outp, outp1...), err
		}
		return outp, err
	}, 15*time.Second); err != nil {
		return
	}

	outp, err := cmd.Exec("users", "read", "--id="+id)
	if err != nil {
		t.Fatal(string(outp), err)
		return
	}
	if !strings.Contains(string(outp), email) ||
		!strings.Contains(string(outp), username) ||
		!strings.Contains(string(outp), id) {
		t.Fatal(string(outp))
		return
	}

	// no password
	outp, err = cmd.Exec("users", "login", "--email="+email)
	if err == nil {
		t.Fatal(string(outp))
		return
	}

	// wrong password
	outp, err = cmd.Exec("users", "login", "--email="+email, "--password=somethingincorrect")
	if err == nil {
		t.Fatal(string(outp))
		return
	}

	outp, err = cmd.Exec("users", "login", "--username="+username, "--password="+password)
	if err != nil {
		t.Fatal(string(outp), err)
		return
	}
	loginRsp := map[string]interface{}{}
	err = json.Unmarshal(outp, &loginRsp)
	if err != nil {
		t.Fatal(err)
		return
	}
	session, ok := loginRsp["session"].(map[string]interface{})
	if !ok {
		t.Fatal(string(outp))
		return
	}
	sessionID := session["id"].(string)
	sessionUsername := session["username"].(string)
	if sessionUsername != username {
		t.Fatal(string(outp))
		return
	}

	if len(sessionID) == 0 {
		t.Fatal(string(outp))
		return
	}

	outp, err = cmd.Exec("users", "login", "--email="+email, "--password="+password)
	if err != nil {
		t.Fatal(string(outp), err)
		return
	}

	outp, err = cmd.Exec("users", "login", "--email="+email, "--password="+password)
	if err != nil {
		t.Fatal(string(outp), err)
		return
	}

	outp, err = cmd.Exec("users", "readSession", "--sessionId="+sessionID)
	if err != nil {
		t.Fatal(string(outp), err)
		return
	}

	loginRsp = map[string]interface{}{}
	err = json.Unmarshal(outp, &loginRsp)
	if err != nil {
		t.Fatal(err)
		return
	}
	session, ok = loginRsp["session"].(map[string]interface{})
	if !ok {
		t.Fatal(string(outp))
		return
	}
	sessionID = session["id"].(string)
	sessionUsername = session["username"].(string)
	if sessionUsername != username {
		t.Fatal(string(outp))
		return
	}
}
