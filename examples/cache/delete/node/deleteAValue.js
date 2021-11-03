const { CacheService } = require("m3o/cache");

// Delete a value from the cache. If key not found a success response is returned.
async function deleteAvalue() {
  let cacheService = new CacheService(process.env.MICRO_API_TOKEN);
  let rsp = await cacheService.delete({
    key: "foo",
  });
  console.log(rsp);
}

deleteAvalue();
