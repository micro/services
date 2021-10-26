const { StockService } = require("m3o/stock");

// Get the historic order book and each trade by timestamp
async function orderBookHistory() {
  let stockService = new StockService(process.env.MICRO_API_TOKEN);
  let rsp = await stockService.orderBook({
    date: "2020-10-01",
    end: "2020-10-01T11:00:00Z",
    limit: 3,
    start: "2020-10-01T10:00:00Z",
    stock: "AAPL",
  });
  console.log(rsp);
}

orderBookHistory();
