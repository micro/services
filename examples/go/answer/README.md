# Answer

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Answer/api](https://m3o.com/Answer/api).

Endpoints:

## Question

Ask a question and receive an instant answer


[https://m3o.com/answer/api#Question](https://m3o.com/answer/api#Question)

```go
package example

import(
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
```
