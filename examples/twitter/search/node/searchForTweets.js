const { TwitterService } = require("m3o/twitter");

// Search for tweets with a simple query
async function searchForTweets() {
  let twitterService = new TwitterService(process.env.MICRO_API_TOKEN);
  let rsp = await twitterService.search({
    query: "cats",
  });
  console.log(rsp);
}

searchForTweets();
