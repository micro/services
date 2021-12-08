package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/micro/micro/v3/service/logger"

	"github.com/micro/services/translation/domain"
	pb "github.com/micro/services/translation/proto"
)

type translation struct {
	youdaoCfg domain.YoudaoConfig
}

func NewTranslation(youdaoCfg domain.YoudaoConfig) *translation {
	return &translation{
		youdaoCfg: youdaoCfg,
	}
}

// Youdao leverages https://ai.youdao.com APIs
func (t *translation) Youdao(_ context.Context, req *pb.YoudaoRequest, rsp *pb.YoudaoResponse) error {
	youdao := domain.NewYoudao(t.youdaoCfg, req)
	result, err := youdao.Translate()
	if err != nil {
		logger.Errorf("get youdaoCfg translation result error: %s", err)
	}

	if result.ErrorCode != "0" {
		logger.Errorf("youdao translation error code is not 0, response: %+v; For more information: https://bit.ly/3rLp4PH", result)
		return errors.New(fmt.Sprintf("youdao translation response error code is not 0; code=%s", result.ErrorCode))
	}

	rsp.Query = result.Query
	rsp.Translation = result.Translation
	rsp.Language = result.L
	rsp.TranslationSpeakUrl = result.TSpeakUrl
	rsp.WebDict = result.Webdict
	rsp.Dict = result.Dict
	rsp.QuerySpeakUrl = result.SpeakUrl
	rsp.IsWord = result.IsWord

	rsp.Basic = &pb.YoudaoBasic{
		ExamType:   result.Basic.ExamType,
		Explains:   result.Basic.Explains,
		Phonetic:   result.Basic.Phonetic,
		UkPhonetic: result.Basic.UkPhonetic,
		UkSpeech:   result.Basic.UkSpeech,
		UsPhonetic: result.Basic.UsPhonetic,
		UsSpeech:   result.Basic.UsSpeech,
	}

	for _, v := range result.Basic.Wfs {
		rsp.Basic.Variants = append(rsp.Basic.Variants, &pb.YoudaoVariant{
			Name:  v.Wf.Name,
			Value: v.Wf.Value,
		})
	}

	for _, v := range result.Web {
		rsp.Web = append(rsp.Web, &pb.YoudaoWeb{
			Key:   v.Key,
			Value: v.Value,
		})
	}

	return nil
}
