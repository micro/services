package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/cache"
)

// Delete a value from the cache. If key not found a success response is returned.
func DeleteAvalue() {
	cacheService := cache.NewCacheService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := cacheService.Delete(&cache.DeleteRequest{
		Key: "foo",
	})
	fmt.Println(rsp, err)
}
