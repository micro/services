import * as sunnah from "m3o/sunnah";

// Hadiths returns a list of hadiths and their corresponding text for a
// given book within a collection.
async function ListTheHadithsInAbook() {
  let sunnahService = new sunnah.SunnahService(process.env.MICRO_API_TOKEN);
  let rsp = await sunnahService.hadiths({
    book: 1,
    collection: "bukhari",
  });
  console.log(rsp);
}

await ListTheHadithsInAbook();
