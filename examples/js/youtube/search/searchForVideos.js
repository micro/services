const { YoutubeService } = require("m3o/youtube");

// Search for videos on YouTube
async function searchForVideos() {
  let youtubeService = new YoutubeService(process.env.MICRO_API_TOKEN);
  let rsp = await youtubeService.search({
    query: "donuts",
  });
  console.log(rsp);
}

searchForVideos();
