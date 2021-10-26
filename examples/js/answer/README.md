# Answer

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Answer/api](https://m3o.com/Answer/api).

Endpoints:

## Question

Ask a question and receive an instant answer


[https://m3o.com/answer/api#Question](https://m3o.com/answer/api#Question)

```js
const { AnswerService } = require('m3o/answer');

// Ask a question and receive an instant answer
async function askAquestion() {
	let answerService = new AnswerService(process.env.MICRO_API_TOKEN)
	let rsp = await answerService.question({
  "query": "microsoft"
})
	console.log(rsp)
}

askAquestion()
```
