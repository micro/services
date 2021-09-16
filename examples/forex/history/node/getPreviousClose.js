import * as forex from "m3o/forex";

// Returns the data for the previous close
async function GetPreviousClose() {
  let forexService = new forex.ForexService(process.env.MICRO_API_TOKEN);
  let rsp = await forexService.history({
    symbol: "GBPUSD",
  });
  console.log(rsp);
}

await GetPreviousClose();
