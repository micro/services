package example

import (
	"fmt"
	"github.com/micro/services/clients/go/cache"
	"os"
)

// Get an item from the cache by key
func GetAvalue() {
	cacheService := cache.NewCacheService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := cacheService.Get(&cache.GetRequest{
		Key: "foo",
	})
	fmt.Println(rsp, err)
}
