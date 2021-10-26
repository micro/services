const { RssService } = require("m3o/rss");

// Remove an RSS feed by name
async function removeAfeed() {
  let rssService = new RssService(process.env.MICRO_API_TOKEN);
  let rsp = await rssService.remove({
    name: "bbc",
  });
  console.log(rsp);
}

removeAfeed();
