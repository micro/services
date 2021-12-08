package domain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/logger"
	"github.com/pkg/errors"

	pb "github.com/micro/services/translation/proto"
)

type YoudaoVariants struct {
	Wf struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"wf"`
}

type YoudaoBasic struct {
	ExamType   []string         `json:"exam_type"`
	Explains   []string         `json:"explains"`
	Phonetic   string           `json:"phonetic"`
	UkPhonetic string           `json:"uk-phonetic"`
	UkSpeech   string           `json:"uk-speech"`
	UsPhonetic string           `json:"us-phonetic"`
	UsSpeech   string           `json:"us-speech"`
	Wfs        []YoudaoVariants `json:"wfs"`
}

type YoudaoWeb struct {
	Key   string
	Value []string
}

type YoudaoResponse struct {
	ErrorCode   string            `json:"errorCode"`
	RequestId   string            `json:"requestId"`
	Query       string            `json:"query"`
	Translation []string          `json:"translation"`
	Dict        map[string]string `json:"dict"`
	Webdict     map[string]string `json:"webdict"`
	Web         []YoudaoWeb       `json:"web"`
	Basic       YoudaoBasic       `json:"basic"`
	L           string            `json:"l"`
	TSpeakUrl   string            `json:"tSpeakUrl"`
	SpeakUrl    string            `json:"speakUrl"`
	IsWord      bool              `json:"isWord"`
}

type YoudaoConfig struct {
	Api    string
	AppKey string
	Secret string
}

type youdao struct {
	config YoudaoConfig
	req    *pb.YoudaoRequest
}

func (y *youdao) generateV3Sign(input, salt string, timestamp int64) string {
	decode := []byte(fmt.Sprintf("%s%s%s%d%s", y.config.AppKey, input, salt, timestamp, y.config.Secret))
	return fmt.Sprintf("%x", sha256.Sum256(decode))
}

func NewYoudao(cfg YoudaoConfig, req *pb.YoudaoRequest) *youdao {
	return &youdao{
		config: cfg,
		req:    req,
	}
}

// Translate the word
func (y *youdao) Translate() (result YoudaoResponse, err error) {

	from, to := y.req.From, y.req.To
	if from == "" {
		from = "auto"
	}

	if to == "" {
		to = "auto"
	}

	voice := y.req.Voice
	if voice == "" {
		voice = "0"
	}

	vals := url.Values{}
	vals.Set("q", y.req.Q)
	vals.Set("from", from)
	vals.Set("to", to)
	vals.Set("appKey", y.config.AppKey)

	input := func() string {
		lenQ := len(y.req.Q)
		if lenQ > 20 {
			return fmt.Sprintf("%s%d%s", y.req.Q[0:10], lenQ, y.req.Q[lenQ-10:])
		}

		return y.req.Q
	}()
	ts := time.Now().Unix()
	salt, _ := uuid.NewUUID()
	vals.Set("salt", salt.String())
	vals.Set("sign", y.generateV3Sign(input, salt.String(), ts))
	vals.Set("signType", "v3")

	vals.Set("curtime", fmt.Sprintf("%d", ts))
	vals.Set("ext", "mp3")
	vals.Set("voice", voice)
	vals.Set("strict", fmt.Sprintf("%t", y.req.Strict))

	client := http.Client{Timeout: 10 * time.Second}
	api := fmt.Sprintf("%s?%s", y.config.Api, vals.Encode())
	resp, err := client.Get(api)
	if err != nil {
		return result, errors.Wrap(err, "youdao translation api error: "+err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return result, errors.Wrap(err, "youdao translation api response status is not 200: "+resp.Status)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	b, _ := ioutil.ReadAll(resp.Body)

	logger.Infof("youdao translation api: %s, response: %s", api, string(b))

	if err := json.Unmarshal(b, &result); err != nil {
		return result, errors.Wrap(err, "youdao translation response json unmarshal error: "+err.Error())
	}

	return result, nil
}
