import * as twitter from "m3o/twitter";

// Search for tweets with a simple query
async function SearchForTweets() {
  let twitterService = new twitter.TwitterService(process.env.MICRO_API_TOKEN);
  let rsp = await twitterService.search({
    query: "cats",
  });
  console.log(rsp);
}

await SearchForTweets();
