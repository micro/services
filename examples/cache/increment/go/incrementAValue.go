package example

import (
	"fmt"
	"github.com/micro/services/clients/go/cache"
	"os"
)

// Increment a value (if it's a number)
func IncrementAvalue() {
	cacheService := cache.NewCacheService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := cacheService.Increment(&cache.IncrementRequest{
		Key:   "counter",
		Value: 2,
	})
	fmt.Println(rsp, err)
}
