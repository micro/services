package example

import (
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
