# Sentiment

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Sentiment/api](https://m3o.com/Sentiment/api).

Endpoints:

## Analyze

Analyze and score a piece of text


[https://m3o.com/sentiment/api#Analyze](https://m3o.com/sentiment/api#Analyze)

```go
package example

import(
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
```
