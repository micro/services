const { CacheService } = require("micro-js-client/cache");

// Decrement a value (if it's a number). If key not found it is equivalent to set.
async function decrementAvalue() {
  let cacheService = new CacheService(process.env.MICRO_API_TOKEN);
  let rsp = await cacheService.decrement({
    key: "counter",
    value: 2,
  });
  console.log(rsp);
}

decrementAvalue();
