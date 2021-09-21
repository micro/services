import * as url from "m3o/url";

// Shortens a destination URL and returns a full short URL.
async function ShortenAlongUrl() {
  let urlService = new url.UrlService(process.env.MICRO_API_TOKEN);
  let rsp = await urlService.shorten({
    destinationURL: "https://mysite.com/this-is-a-rather-long-web-address",
  });
  console.log(rsp);
}

await ShortenAlongUrl();
