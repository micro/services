# Quran

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Quran/api](https://m3o.com/Quran/api).

Endpoints:

## Chapters

List the Chapters (surahs) of the Quran


[https://m3o.com/quran/api#Chapters](https://m3o.com/quran/api#Chapters)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/quran"
)

// List the Chapters (surahs) of the Quran
func ListChapters() {
	quranService := quran.NewQuranService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := quranService.Chapters(&quran.ChaptersRequest{
		Language: "en",

	})
	fmt.Println(rsp, err)
}
```
## Summary

Get a summary for a given chapter (surah)


[https://m3o.com/quran/api#Summary](https://m3o.com/quran/api#Summary)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/quran"
)

// Get a summary for a given chapter (surah)
func GetChapterSummary() {
	quranService := quran.NewQuranService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := quranService.Summary(&quran.SummaryRequest{
		Chapter: 1,

	})
	fmt.Println(rsp, err)
}
```
## Verses

Lookup the verses (ayahs) for a chapter including
translation, interpretation and breakdown by individual
words.


[https://m3o.com/quran/api#Verses](https://m3o.com/quran/api#Verses)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/quran"
)

// Lookup the verses (ayahs) for a chapter including
// translation, interpretation and breakdown by individual
// words.
func GetVersesOfAchapter() {
	quranService := quran.NewQuranService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := quranService.Verses(&quran.VersesRequest{
		Chapter: 1,

	})
	fmt.Println(rsp, err)
}
```
## Search

Search the Quran for any form of query or questions


[https://m3o.com/quran/api#Search](https://m3o.com/quran/api#Search)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/quran"
)

// Search the Quran for any form of query or questions
func SearchTheQuran() {
	quranService := quran.NewQuranService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := quranService.Search(&quran.SearchRequest{
		Query: "messenger",

	})
	fmt.Println(rsp, err)
}
```
