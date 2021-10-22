package example

import (
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
