import * as quran from "m3o/quran";

// List the Chapters (surahs) of the Quran
async function ListChapters() {
  let quranService = new quran.QuranService(process.env.MICRO_API_TOKEN);
  let rsp = await quranService.chapters({
    language: "en",
  });
  console.log(rsp);
}

await ListChapters();
