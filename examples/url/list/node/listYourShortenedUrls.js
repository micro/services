import * as url from "@m3o/services/url";

// List information on all the shortened URLs that you have created
async function ListYourShortenedUrls() {
  let urlService = new url.UrlService(process.env.MICRO_API_TOKEN);
  let rsp = await urlService.list({});
  console.log(rsp);
}

await ListYourShortenedUrls();
