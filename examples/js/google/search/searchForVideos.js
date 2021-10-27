const { GoogleService } = require("m3o/google");

// Search for videos on Google
async function searchForVideos() {
  let googleService = new GoogleService(process.env.MICRO_API_TOKEN);
  let rsp = await googleService.search({
    query: "how to make donuts",
  });
  console.log(rsp);
}

searchForVideos();
