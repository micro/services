import * as sunnah from "m3o/sunnah";

// Get a list of books from within a collection. A book can contain many chapters
// each with its own hadiths.
async function GetTheBooksWithinAcollection() {
  let sunnahService = new sunnah.SunnahService(process.env.MICRO_API_TOKEN);
  let rsp = await sunnahService.books({
    collection: "bukhari",
  });
  console.log(rsp);
}

await GetTheBooksWithinAcollection();
