import * as forex from "m3o/forex";

// Get the latest quote for the forex
async function GetAfxQuote() {
  let forexService = new forex.ForexService(process.env.MICRO_API_TOKEN);
  let rsp = await forexService.quote({
    symbol: "GBPUSD",
  });
  console.log(rsp);
}

await GetAfxQuote();
