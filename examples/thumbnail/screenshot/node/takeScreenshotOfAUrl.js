import * as thumbnail from "@m3o/services/thumbnail";

// Create a thumbnail screenshot by passing in a url, height and width
async function TakeScreenshotOfAurl() {
  let thumbnailService = new thumbnail.ThumbnailService(
    process.env.MICRO_API_TOKEN
  );
  let rsp = await thumbnailService.screenshot({
    height: 600,
    url: "https://m3o.com",
    width: 600,
  });
  console.log(rsp);
}

await TakeScreenshotOfAurl();
