import * as answer from "@m3o/services/answer";

// Ask a question and receive an instant answer
async function AskAquestion() {
  let answerService = new answer.AnswerService(process.env.MICRO_API_TOKEN);
  let rsp = await answerService.question({
    query: "google",
  });
  console.log(rsp);
}

await AskAquestion();
