import * as cache from "@m3o/services/cache";

// Decrement a value (if it's a number)
async function DecrementAvalue() {
  let cacheService = new cache.CacheService(process.env.MICRO_API_TOKEN);
  let rsp = await cacheService.decrement({
    key: "counter",
    value: 2,
  });
  console.log(rsp);
}

await DecrementAvalue();
