import * as answer from "m3o/answer";

// Ask a question and receive an instant answer
async function AskAquestion() {
  let answerService = new answer.AnswerService(process.env.MICRO_API_TOKEN);
  let rsp = await answerService.question({
    query: "microsoft",
  });
  console.log(rsp);
}

await AskAquestion();
