import * as stock from "@m3o/services/stock";

// Get the last price for a given stock ticker
async function GetAstockPrice() {
  let stockService = new stock.StockService(process.env.MICRO_API_TOKEN);
  let rsp = await stockService.price({
    symbol: "AAPL",
  });
  console.log(rsp);
}

await GetAstockPrice();
