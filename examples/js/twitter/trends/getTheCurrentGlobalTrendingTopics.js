const { TwitterService } = require("m3o/twitter");

// Get the current global trending topics
async function getTheCurrentGlobalTrendingTopics() {
  let twitterService = new TwitterService(process.env.MICRO_API_TOKEN);
  let rsp = await twitterService.trends({});
  console.log(rsp);
}

getTheCurrentGlobalTrendingTopics();
