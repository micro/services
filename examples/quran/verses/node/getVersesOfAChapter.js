import * as quran from "m3o/quran";

// Lookup the verses (ayahs) for a chapter including
// translation, interpretation and breakdown by individual
// words.
async function GetVersesOfAchapter() {
  let quranService = new quran.QuranService(process.env.MICRO_API_TOKEN);
  let rsp = await quranService.verses({
    chapter: 1,
  });
  console.log(rsp);
}

await GetVersesOfAchapter();
