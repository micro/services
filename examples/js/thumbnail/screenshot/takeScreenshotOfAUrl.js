const { ThumbnailService } = require("m3o/thumbnail");

// Create a thumbnail screenshot by passing in a url, height and width
async function takeScreenshotOfAurl() {
  let thumbnailService = new ThumbnailService(process.env.MICRO_API_TOKEN);
  let rsp = await thumbnailService.screenshot({
    height: 600,
    url: "https://m3o.com",
    width: 600,
  });
  console.log(rsp);
}

takeScreenshotOfAurl();
