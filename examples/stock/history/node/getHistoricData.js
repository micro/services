import * as stock from "m3o/stock";

// Get the historic open-close for a given day
async function GetHistoricData() {
  let stockService = new stock.StockService(process.env.MICRO_API_TOKEN);
  let rsp = await stockService.history({
    date: "2020-10-01",
    stock: "AAPL",
  });
  console.log(rsp);
}

await GetHistoricData();
