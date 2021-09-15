import * as cache from "@m3o/services/cache";

// Increment a value (if it's a number)
async function IncrementAvalue() {
  let cacheService = new cache.CacheService(process.env.MICRO_API_TOKEN);
  let rsp = await cacheService.increment({
    key: "counter",
    value: 2,
  });
  console.log(rsp);
}

await IncrementAvalue();
