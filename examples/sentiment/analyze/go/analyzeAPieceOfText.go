package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/sentiment"
)

// Analyze and score a piece of text
func AnalyzeApieceOfText() {
	sentimentService := sentiment.NewSentimentService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := sentimentService.Analyze(&sentiment.AnalyzeRequest{
		Text: "this is amazing",
	})
	fmt.Println(rsp, err)
}
