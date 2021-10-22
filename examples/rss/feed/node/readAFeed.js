const { RssService } = require("m3o/rss");

// Get an RSS feed by name. If no name is given, all feeds are returned. Default limit is 25 entries.
async function readAfeed() {
  let rssService = new RssService(process.env.MICRO_API_TOKEN);
  let rsp = await rssService.feed({
    name: "bbc",
  });
  console.log(rsp);
}

readAfeed();
