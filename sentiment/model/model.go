package model

import (
	"github.com/cdipaolo/sentiment"
	"github.com/micro/micro/v3/service/logger"
)

var (
	model *sentiment.Models
)

func init() {
	// load sentiment analysis tool
	md, err := sentiment.Restore()
	if err != nil {
		logger.Fatal(err)
	}
	model = &md
}

func Analyze(text string) float64 {
	an := model.SentimentAnalysis(text, sentiment.English)
	return float64(an.Score)
}
