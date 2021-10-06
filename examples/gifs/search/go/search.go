package example

import (
	"fmt"
	"github.com/micro/services/clients/go/gifs"
	"os"
)

// Search for a gif
func Search() {
	gifsService := gifs.NewGifsService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := gifsService.Search(&gifs.SearchRequest{
		Limit: 2,
		Query: "dogs",
	})
	fmt.Println(rsp, err)
}
