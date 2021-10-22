const { UrlService } = require("m3o/url");

// Proxy returns the destination URL of a short URL.
async function resolveAshortUrlToAlongDestinationUrl() {
  let urlService = new UrlService(process.env.MICRO_API_TOKEN);
  let rsp = await urlService.proxy({
    shortURL: "https://m3o.one/u/ck6SGVkYp",
  });
  console.log(rsp);
}

resolveAshortUrlToAlongDestinationUrl();
