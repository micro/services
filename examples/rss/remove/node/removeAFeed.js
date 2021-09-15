import * as rss from "@m3o/services/rss";

// Remove an RSS feed by name
async function RemoveAfeed() {
  let rssService = new rss.RssService(process.env.MICRO_API_TOKEN);
  let rsp = await rssService.remove({
    name: "bbc",
  });
  console.log(rsp);
}

await RemoveAfeed();
