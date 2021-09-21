import * as twitter from "m3o/twitter";

// Get a user's twitter profile
async function GetAusersTwitterProfile() {
  let twitterService = new twitter.TwitterService(process.env.MICRO_API_TOKEN);
  let rsp = await twitterService.user({
    username: "crufter",
  });
  console.log(rsp);
}

await GetAusersTwitterProfile();
