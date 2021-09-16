import * as sentiment from "@m3o/services/sentiment";

// Analyze and score a piece of text
async function AnalyzeApieceOfText() {
  let sentimentService = new sentiment.SentimentService(
    process.env.MICRO_API_TOKEN
  );
  let rsp = await sentimentService.analyze({
    text: "this is amazing",
  });
  console.log(rsp);
}

await AnalyzeApieceOfText();
