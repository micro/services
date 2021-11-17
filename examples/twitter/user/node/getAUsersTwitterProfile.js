const { TwitterService } = require("micro-js-client/twitter");

// Get a user's twitter profile
async function getAusersTwitterProfile() {
  let twitterService = new TwitterService(process.env.MICRO_API_TOKEN);
  let rsp = await twitterService.user({
    username: "crufter",
  });
  console.log(rsp);
}

getAusersTwitterProfile();
