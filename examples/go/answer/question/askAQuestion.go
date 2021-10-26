package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/answer"
)

// Ask a question and receive an instant answer
func AskAquestion() {
	answerService := answer.NewAnswerService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := answerService.Question(&answer.QuestionRequest{
		Query: "microsoft",
	})
	fmt.Println(rsp, err)
}
