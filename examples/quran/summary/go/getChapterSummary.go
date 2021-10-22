package example

import (
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
