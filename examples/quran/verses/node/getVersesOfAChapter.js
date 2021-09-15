import * as quran from "@m3o/services/quran";

// Lookup the verses (ayahs) for a chapter
async function GetVersesOfAchapter() {
  let quranService = new quran.QuranService(process.env.MICRO_API_TOKEN);
  let rsp = await quranService.verses({
    chapter: 1,
  });
  console.log(rsp);
}

await GetVersesOfAchapter();
