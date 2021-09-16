import * as cache from "m3o/cache";

// Set an item in the cache. Overwrites any existing value already set.
async function SetAvalue() {
  let cacheService = new cache.CacheService(process.env.MICRO_API_TOKEN);
  let rsp = await cacheService.set({
    key: "foo",
    value: "bar",
  });
  console.log(rsp);
}

await SetAvalue();
