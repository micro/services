const { GifsService } = require("micro-js-client/gifs");

// Search for a GIF
async function search() {
  let gifsService = new GifsService(process.env.MICRO_API_TOKEN);
  let rsp = await gifsService.search({
    limit: 2,
    query: "dogs",
  });
  console.log(rsp);
}

search();
