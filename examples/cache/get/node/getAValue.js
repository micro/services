const { CacheService } = require("m3o/cache");

// Get an item from the cache by key. If key is not found, an empty response is returned.
async function getAvalue() {
  let cacheService = new CacheService(process.env.MICRO_API_TOKEN);
  let rsp = await cacheService.get({
    key: "foo",
  });
  console.log(rsp);
}

getAvalue();
