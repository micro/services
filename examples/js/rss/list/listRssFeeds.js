const { RssService } = require("m3o/rss");

// List the saved RSS fields
async function listRssFeeds() {
  let rssService = new RssService(process.env.MICRO_API_TOKEN);
  let rsp = await rssService.list({});
  console.log(rsp);
}

listRssFeeds();
