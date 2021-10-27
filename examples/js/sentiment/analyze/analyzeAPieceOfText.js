const { SentimentService } = require("m3o/sentiment");

// Analyze and score a piece of text
async function analyzeApieceOfText() {
  let sentimentService = new SentimentService(process.env.MICRO_API_TOKEN);
  let rsp = await sentimentService.analyze({
    text: "this is amazing",
  });
  console.log(rsp);
}

analyzeApieceOfText();
