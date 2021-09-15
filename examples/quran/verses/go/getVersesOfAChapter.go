package example

import (
	"fmt"
	"github.com/micro/services/clients/go/quran"
	"os"
)

// Lookup the verses (ayahs) for a chapter
func GetVersesOfAchapter() {
	quranService := quran.NewQuranService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := quranService.Verses(&quran.VersesRequest{
		Chapter: 1,
	})
	fmt.Println(rsp, err)
}
