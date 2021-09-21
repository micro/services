import * as sunnah from "m3o/sunnah";

// Get a list of available collections. A collection is
// a compilation of hadiths collected and written by an author.
async function ListAvailableCollections() {
  let sunnahService = new sunnah.SunnahService(process.env.MICRO_API_TOKEN);
  let rsp = await sunnahService.collections({});
  console.log(rsp);
}

await ListAvailableCollections();
