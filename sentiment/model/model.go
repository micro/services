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

	// no words, just return whats scored
	if len(an.Words) == 0 {
		return float64(an.Score)
	}

	// take each word score then divide by num words
	var total float64

	for _, word := range an.Words {
		total += float64(word.Score)
	}

	// get the overall score
	return total / float64(len(an.Words))
}
