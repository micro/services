package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/cache"
)

// Delete a value from the cache
func DeleteAvalue() {
	cacheService := cache.NewCacheService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := cacheService.Delete(&cache.DeleteRequest{
		Key: "foo",
	})
	fmt.Println(rsp, err)
}
