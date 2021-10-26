const { StockService } = require("m3o/stock");

// Get the last quote for the stock
async function getAstockQuote() {
  let stockService = new StockService(process.env.MICRO_API_TOKEN);
  let rsp = await stockService.quote({
    symbol: "AAPL",
  });
  console.log(rsp);
}

getAstockQuote();
