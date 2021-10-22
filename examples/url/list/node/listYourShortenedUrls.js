const { UrlService } = require("m3o/url");

// List information on all the shortened URLs that you have created
async function listYourShortenedUrls() {
  let urlService = new UrlService(process.env.MICRO_API_TOKEN);
  let rsp = await urlService.list({});
  console.log(rsp);
}

listYourShortenedUrls();
