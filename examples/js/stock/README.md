# Stock

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Stock/api](https://m3o.com/Stock/api).

Endpoints:

## Price

Get the last price for a given stock ticker


[https://m3o.com/stock/api#Price](https://m3o.com/stock/api#Price)

```js
const { StockService } = require('m3o/stock');

// Get the last price for a given stock ticker
async function getAstockPrice() {
	let stockService = new StockService(process.env.MICRO_API_TOKEN)
	let rsp = await stockService.price({
  "symbol": "AAPL"
})
	console.log(rsp)
}

getAstockPrice()
```
## Quote

Get the last quote for the stock


[https://m3o.com/stock/api#Quote](https://m3o.com/stock/api#Quote)

```js
const { StockService } = require('m3o/stock');

// Get the last quote for the stock
async function getAstockQuote() {
	let stockService = new StockService(process.env.MICRO_API_TOKEN)
	let rsp = await stockService.quote({
  "symbol": "AAPL"
})
	console.log(rsp)
}

getAstockQuote()
```
## History

Get the historic open-close for a given day


[https://m3o.com/stock/api#History](https://m3o.com/stock/api#History)

```js
const { StockService } = require('m3o/stock');

// Get the historic open-close for a given day
async function getHistoricData() {
	let stockService = new StockService(process.env.MICRO_API_TOKEN)
	let rsp = await stockService.history({
  "date": "2020-10-01",
  "stock": "AAPL"
})
	console.log(rsp)
}

getHistoricData()
```
## OrderBook

Get the historic order book and each trade by timestamp


[https://m3o.com/stock/api#OrderBook](https://m3o.com/stock/api#OrderBook)

```js
const { StockService } = require('m3o/stock');

// Get the historic order book and each trade by timestamp
async function orderBookHistory() {
	let stockService = new StockService(process.env.MICRO_API_TOKEN)
	let rsp = await stockService.orderBook({
  "date": "2020-10-01",
  "end": "2020-10-01T11:00:00Z",
  "limit": 3,
  "start": "2020-10-01T10:00:00Z",
  "stock": "AAPL"
})
	console.log(rsp)
}

orderBookHistory()
```
