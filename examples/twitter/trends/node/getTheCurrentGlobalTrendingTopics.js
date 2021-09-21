import * as twitter from "m3o/twitter";

// Get the current global trending topics
async function GetTheCurrentGlobalTrendingTopics() {
  let twitterService = new twitter.TwitterService(process.env.MICRO_API_TOKEN);
  let rsp = await twitterService.trends({});
  console.log(rsp);
}

await GetTheCurrentGlobalTrendingTopics();
