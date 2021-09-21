import * as forex from "m3o/forex";

// Get the latest price for a given forex ticker
async function GetAnFxPrice() {
  let forexService = new forex.ForexService(process.env.MICRO_API_TOKEN);
  let rsp = await forexService.price({
    symbol: "GBPUSD",
  });
  console.log(rsp);
}

await GetAnFxPrice();
