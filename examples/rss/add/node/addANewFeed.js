import * as rss from "m3o/rss";

// Add a new RSS feed with a name, url, and category
async function AddAnewFeed() {
  let rssService = new rss.RssService(process.env.MICRO_API_TOKEN);
  let rsp = await rssService.add({
    category: "news",
    name: "bbc",
    url: "http://feeds.bbci.co.uk/news/rss.xml",
  });
  console.log(rsp);
}

await AddAnewFeed();
