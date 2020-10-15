// +build blog

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
	p "github.com/micro/services/blog/posts/handler"
)

const (
	retryCount          = 1
	signupSuccessString = "Signup complete"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func setupBlogTests(serv test.Server, t *test.T) {
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
		{envVar: "POSTS_SVC", deflt: "../../../blog/posts"},
		{envVar: "TAGS_SVC", deflt: "../../../blog/tags"},
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
		list := []string{"posts", "tags"}
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
	outp, err := serv.Command().Exec("auth", "create", "rule", "--access=granted", "--scope=''", "--resource=\"service:posts:*\"", "posts")
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

func TestPostsService(t *testing.T) {
	test.TrySuite(t, testPosts, retryCount)
}

func testPosts(t *test.T) {
	t.Parallel()

	serv := test.NewServer(t, test.WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	setupBlogTests(serv, t)

	time.Sleep(5 * time.Second)
	cmd := serv.Command()

	if err := test.Try("Save post", t, func() ([]byte, error) {
		// Attention! The content must be unquoted, don't add quotes.
		outp, err := cmd.Exec("posts", "--id=1", "--title=Hi", "--content=Hi there", "--tags=a,b", "save")
		if err != nil {
			outp1, _ := cmd.Exec("logs", "posts")
			return append(outp, outp1...), err
		}
		return outp, err
	}, 15*time.Second); err != nil {
		return
	}

	outp, err := cmd.Exec("posts", "query")
	if err != nil {
		t.Fatal(string(outp))
	}

	expected := []p.Post{
		{
			ID:      "1",
			Title:   "Hi",
			Content: "Hi there",
			Tags:    []string{"a", "b"},
		},
	}
	type rsp struct {
		Posts []p.Post `json:"posts"`
	}
	var actual rsp
	json.Unmarshal(outp, &actual)
	if len(actual.Posts) == 0 {
		t.Fatal(string(outp))
		return
	}

	if expected[0].ID != actual.Posts[0].ID ||
		expected[0].Title != actual.Posts[0].Title ||
		expected[0].Content != actual.Posts[0].Content ||
		len(expected[0].Tags) != len(actual.Posts[0].Tags) {
		t.Fatal(expected[0], actual.Posts[0])
	}

	if err := test.Try("Save post", t, func() ([]byte, error) {
		outp, err = cmd.Exec("tags", "list", "--type=post-tag")
		type tagsRsp struct {
			Tags []p.Post `json:"tags"`
		}
		var tagsActual tagsRsp
		json.Unmarshal(outp, &tagsActual)
		if len(tagsActual.Tags) == 0 {
			outp1, _ := cmd.Exec("logs", "tags")
			return append(outp, outp1...), errors.New("Unexpected output")
		}
		if len(tagsActual.Tags) != 2 {
			return outp, errors.New("Unexpected output")
		}
		return outp, err
	}, 15*time.Second); err != nil {
		return
	}
}
