import * as sunnah from "m3o/sunnah";

// Get all the chapters of a given book within a collection.
async function ListTheChaptersInAbook() {
  let sunnahService = new sunnah.SunnahService(process.env.MICRO_API_TOKEN);
  let rsp = await sunnahService.chapters({
    book: 1,
    collection: "bukhari",
  });
  console.log(rsp);
}

await ListTheChaptersInAbook();
