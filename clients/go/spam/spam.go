package spam

import (
	"github.com/micro/micro-go/client"
)

func NewSpamService(token string) *SpamService {
	return &SpamService{
		client: client.NewClient(&client.Options{
			Token: token,
		}),
	}
}

type SpamService struct {
	client *client.Client
}

// Check whether an email is likely to be spam based on its attributes
func (t *SpamService) Classify(request *ClassifyRequest) (*ClassifyResponse, error) {
	rsp := &ClassifyResponse{}
	return rsp, t.client.Call("spam", "Classify", request, rsp)
}

type ClassifyRequest struct {
	// The body of the email
	EmailBody string `json:"emailBody"`
	// The email address it has been sent from
	From string `json:"from"`
	// The subject of the email
	Subject string `json:"subject"`
	// The email address it is being sent to
	To string `json:"to"`
}

type ClassifyResponse struct {
	// The rules that have contributed to this score
	Details []string `json:"details"`
	// Is it spam? Returns true if its score is > 5
	IsSpam bool `json:"isSpam"`
	// The score evaluated for this email. A higher number means it is more likely to be spam
	Score float64 `json:"score"`
}
