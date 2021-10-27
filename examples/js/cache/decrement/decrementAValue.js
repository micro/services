const { CacheService } = require("m3o/cache");

// Decrement a value (if it's a number)
async function decrementAvalue() {
  let cacheService = new CacheService(process.env.MICRO_API_TOKEN);
  let rsp = await cacheService.decrement({
    key: "counter",
    value: 2,
  });
  console.log(rsp);
}

decrementAvalue();
