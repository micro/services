import * as gifs from "m3o/gifs";

// Search for a gif
async function Search() {
  let gifsService = new gifs.GifsService(process.env.MICRO_API_TOKEN);
  let rsp = await gifsService.search({
    limit: 2,
    query: "dogs",
  });
  console.log(rsp);
}

await Search();
