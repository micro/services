const { TwitterService } = require("m3o/twitter");

// Get the timeline for a given user
async function getAtwitterTimeline() {
  let twitterService = new TwitterService(process.env.MICRO_API_TOKEN);
  let rsp = await twitterService.timeline({
    limit: 1,
    username: "m3oservices",
  });
  console.log(rsp);
}

getAtwitterTimeline();
