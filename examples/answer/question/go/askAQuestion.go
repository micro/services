package example

import (
	"fmt"
	"github.com/micro/services/clients/go/answer"
	"os"
)

// Ask a question and receive an instant answer
func AskAquestion() {
	answerService := answer.NewAnswerService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := answerService.Question(&answer.QuestionRequest{
		Query: "google",
	})
	fmt.Println(rsp, err)
}
