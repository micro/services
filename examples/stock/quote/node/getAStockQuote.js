import * as stock from "@m3o/services/stock";

// Get the last quote for the stock
async function GetAstockQuote() {
  let stockService = new stock.StockService(process.env.MICRO_API_TOKEN);
  let rsp = await stockService.quote({
    symbol: "AAPL",
  });
  console.log(rsp);
}

await GetAstockQuote();
