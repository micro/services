import * as quran from "@m3o/services/quran";

// Search the Quran for any form of query or questions
async function SearchTheQuran() {
  let quranService = new quran.QuranService(process.env.MICRO_API_TOKEN);
  let rsp = await quranService.search({
    query: "messenger",
  });
  console.log(rsp);
}

await SearchTheQuran();
