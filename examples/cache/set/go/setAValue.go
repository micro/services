package example

import (
	"fmt"
	"github.com/micro/services/clients/go/cache"
	"os"
)

// Set an item in the cache. Overwrites any existing value already set.
func SetAvalue() {
	cacheService := cache.NewCacheService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := cacheService.Set(&cache.SetRequest{
		Key:   "foo",
		Value: "bar",
	})
	fmt.Println(rsp, err)
}
