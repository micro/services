package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/cache"
)

// Decrement a value (if it's a number)
func DecrementAvalue() {
	cacheService := cache.NewCacheService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := cacheService.Decrement(&cache.DecrementRequest{
		Key:   "counter",
		Value: 2,
	})
	fmt.Println(rsp, err)
}
