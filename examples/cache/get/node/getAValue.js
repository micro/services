import * as cache from "@m3o/services/cache";

// Get an item from the cache by key
async function GetAvalue() {
  let cacheService = new cache.CacheService(process.env.MICRO_API_TOKEN);
  let rsp = await cacheService.get({
    key: "foo",
  });
  console.log(rsp);
}

await GetAvalue();
