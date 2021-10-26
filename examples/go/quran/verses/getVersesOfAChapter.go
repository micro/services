package example

import (
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
