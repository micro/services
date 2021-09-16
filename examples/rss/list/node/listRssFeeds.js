import * as rss from "@m3o/services/rss";

// List the saved RSS fields
async function ListRssFeeds() {
  let rssService = new rss.RssService(process.env.MICRO_API_TOKEN);
  let rsp = await rssService.list({});
  console.log(rsp);
}

await ListRssFeeds();
