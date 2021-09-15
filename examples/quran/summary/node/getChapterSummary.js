import * as quran from "@m3o/services/quran";

// Get a summary for a given chapter (surah)
async function GetChapterSummary() {
  let quranService = new quran.QuranService(process.env.MICRO_API_TOKEN);
  let rsp = await quranService.summary({
    chapter: 1,
  });
  console.log(rsp);
}

await GetChapterSummary();
