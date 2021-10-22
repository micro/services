const { ForexService } = require("m3o/forex");

// Get the latest price for a given forex ticker
async function getAnFxPrice() {
  let forexService = new ForexService(process.env.MICRO_API_TOKEN);
  let rsp = await forexService.price({
    symbol: "GBPUSD",
  });
  console.log(rsp);
}

getAnFxPrice();
