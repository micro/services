# Quran

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Quran/api](https://m3o.com/Quran/api).

Endpoints:

## Chapters

List the Chapters (surahs) of the Quran


[https://m3o.com/quran/api#Chapters](https://m3o.com/quran/api#Chapters)

```js
const { QuranService } = require('m3o/quran');

// List the Chapters (surahs) of the Quran
async function listChapters() {
	let quranService = new QuranService(process.env.MICRO_API_TOKEN)
	let rsp = await quranService.chapters({
  "language": "en"
})
	console.log(rsp)
}

listChapters()
```
## Summary

Get a summary for a given chapter (surah)


[https://m3o.com/quran/api#Summary](https://m3o.com/quran/api#Summary)

```js
const { QuranService } = require('m3o/quran');

// Get a summary for a given chapter (surah)
async function getChapterSummary() {
	let quranService = new QuranService(process.env.MICRO_API_TOKEN)
	let rsp = await quranService.summary({
  "chapter": 1
})
	console.log(rsp)
}

getChapterSummary()
```
## Verses

Lookup the verses (ayahs) for a chapter including
translation, interpretation and breakdown by individual
words.


[https://m3o.com/quran/api#Verses](https://m3o.com/quran/api#Verses)

```js
const { QuranService } = require('m3o/quran');

// Lookup the verses (ayahs) for a chapter including
// translation, interpretation and breakdown by individual
// words.
async function getVersesOfAchapter() {
	let quranService = new QuranService(process.env.MICRO_API_TOKEN)
	let rsp = await quranService.verses({
  "chapter": 1
})
	console.log(rsp)
}

getVersesOfAchapter()
```
## Search

Search the Quran for any form of query or questions


[https://m3o.com/quran/api#Search](https://m3o.com/quran/api#Search)

```js
const { QuranService } = require('m3o/quran');

// Search the Quran for any form of query or questions
async function searchTheQuran() {
	let quranService = new QuranService(process.env.MICRO_API_TOKEN)
	let rsp = await quranService.search({
  "query": "messenger"
})
	console.log(rsp)
}

searchTheQuran()
```
