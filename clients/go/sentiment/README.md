package sentiment

import(
	"github.com/m3o/m3o-go/client"
)

func NewSentimentService(token string) *SentimentService {
	return &SentimentService{
		client: client.NewClient(&client.Options{
			Token: token,
		}),
	}
}

type SentimentService struct {
	client *client.Client
}


// Analyze and score a piece of text
func (t *SentimentService) Analyze(request *AnalyzeRequest) (*AnalyzeResponse, error) {
	rsp := &AnalyzeResponse{}
	return rsp, t.client.Call("sentiment", "Analyze", request, rsp)
}




type AnalyzeRequest struct {
  // The language. Defaults to english.
  Lang string `json:"lang"`
  // The text to analyze
  Text string `json:"text"`
}

type AnalyzeResponse struct {
  // The score of the text {positive is 1, negative is 0}
  Score float64 `json:"score"`
}

# { Sentiment

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Sentiment/api](https://m3o.com/Sentiment/api).

Endpoints:

#analyze

// Analyze and score a piece of text


[https://m3o.com/sentiment/api#analyze](https://m3o.com/sentiment/api#analyze)
