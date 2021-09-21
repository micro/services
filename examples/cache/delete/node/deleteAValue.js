import * as cache from "m3o/cache";

// Delete a value from the cache
async function DeleteAvalue() {
  let cacheService = new cache.CacheService(process.env.MICRO_API_TOKEN);
  let rsp = await cacheService.delete({
    key: "foo",
  });
  console.log(rsp);
}

await DeleteAvalue();
