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
	p "github.com/micro/services/blog/posts/proto"
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

// count is a string in responses...
type protoTag struct {
	Title string `json:"title"`
	Slug  string `json:"slug"`
	Type  string `json:"type"`
	Count string `json:"count"`
}

func testPosts(t *test.T) {
	t.Parallel()

	serv := test.NewServer(t, test.WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	setupBlogTests(serv, t)

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
			Id:      "1",
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

	if expected[0].Id != actual.Posts[0].Id ||
		expected[0].Title != actual.Posts[0].Title ||
		expected[0].Content != actual.Posts[0].Content ||
		len(expected[0].Tags) != len(actual.Posts[0].Tags) {
		t.Fatal(expected[0], actual.Posts[0])
	}

	outp, err = cmd.Exec("tags", "list", "--type=post-tag")
	type tagsRsp struct {
		Tags []protoTag `json:"tags"`
	}
	var tagsActual tagsRsp
	json.Unmarshal(outp, &tagsActual)
	if len(tagsActual.Tags) == 0 {
		outp1, _ := cmd.Exec("logs", "tags")
		t.Fatal(string(append(outp, outp1...)))
		return
	}
	if len(tagsActual.Tags) != 2 {
		t.Fatal(tagsActual.Tags)
		return
	}

	if tagsActual.Tags[0].Count != "1" {
		t.Fatal(tagsActual.Tags[0])
		return
	}
	if tagsActual.Tags[1].Count != "1" {
		t.Fatal(tagsActual.Tags[1])
		return
	}

	time.Sleep(5 * time.Second)
	// Inserting an other post so tag counts increase
	outp, err = cmd.Exec("posts", "--id=2", "--title=Hi1", "--content=Hi there1", "--tags=a,b", "save")
	if err != nil {
		t.Fatal(string(outp))
		return
	}

	outp, err = cmd.Exec("tags", "list", "--type=post-tag")
	json.Unmarshal(outp, &tagsActual)
	if len(tagsActual.Tags) == 0 {
		outp1, _ := cmd.Exec("logs", "tags")
		t.Fatal(string(append(outp, outp1...)))
		return
	}
	if len(tagsActual.Tags) != 2 {
		t.Fatal(tagsActual.Tags)
		return
	}

	if tagsActual.Tags[0].Count != "2" {
		outp1, _ := cmd.Exec("store", "list", "--table=tags")
		outp2, _ := cmd.Exec("store", "list", "--table=posts")
		t.Fatal(tagsActual.Tags[0], string(outp1), string(outp2))
		return
	}
	if tagsActual.Tags[1].Count != "2" {
		outp1, _ := cmd.Exec("store", "list", "--table=tags")
		outp2, _ := cmd.Exec("store", "list", "--table=posts")
		t.Fatal(tagsActual.Tags[1], string(outp1), string(outp2))
		return
	}

	// test updating fields fields and removing tags
	outp, err = cmd.Exec("posts", "--id=2", "--title=Hi2", "--tags=a", "save")
	if err != nil {
		t.Fatal(string(outp))
		return
	}

	outp, err = cmd.Exec("tags", "list", "--type=post-tag")
	json.Unmarshal(outp, &tagsActual)
	if len(tagsActual.Tags) == 0 {
		outp1, _ := cmd.Exec("logs", "tags")
		t.Fatal(string(append(outp, outp1...)))
		return
	}
	if len(tagsActual.Tags) != 2 {
		t.Fatal(tagsActual.Tags)
		return
	}
	for _, tag := range tagsActual.Tags {
		if tag.Title == "b" {
			if tag.Count != "1" {
				t.Fatal("Tag b should have a count 1")
				return
			}
		}
		if tag.Title == "a" {
			if tag.Count != "2" {
				t.Fatal("Tag b should have a count 2")
				return
			}
		}
	}

	outp, err = cmd.Exec("posts", "--id=2", "query")
	if err != nil {
		t.Fatal(string(outp))
		return
	}
	json.Unmarshal(outp, &actual)
	if len(actual.Posts) == 0 {
		t.Fatal(string(outp))
		return
	}
	if actual.Posts[0].Title != "Hi2" ||
		actual.Posts[0].Content != "Hi there1" ||
		actual.Posts[0].Slug != "hi2" || len(actual.Posts[0].Tags) != 1 {
		t.Fatal(actual)
		return
	}
}
