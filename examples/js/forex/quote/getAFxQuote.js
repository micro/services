const { ForexService } = require("m3o/forex");

// Get the latest quote for the forex
async function getAfxQuote() {
  let forexService = new ForexService(process.env.MICRO_API_TOKEN);
  let rsp = await forexService.quote({
    symbol: "GBPUSD",
  });
  console.log(rsp);
}

getAfxQuote();
