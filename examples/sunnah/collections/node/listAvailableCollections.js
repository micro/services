const { SunnahService } = require("m3o/sunnah");

// Get a list of available collections. A collection is
// a compilation of hadiths collected and written by an author.
async function listAvailableCollections() {
  let sunnahService = new SunnahService(process.env.MICRO_API_TOKEN);
  let rsp = await sunnahService.collections({});
  console.log(rsp);
}

listAvailableCollections();
