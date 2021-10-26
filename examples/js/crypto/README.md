# Crypto

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Crypto/api](https://m3o.com/Crypto/api).

Endpoints:

## News

Get news related to a currency


[https://m3o.com/crypto/api#News](https://m3o.com/crypto/api#News)

```js
const { CryptoService } = require('m3o/crypto');

// Get news related to a currency
async function getCryptocurrencyNews() {
	let cryptoService = new CryptoService(process.env.MICRO_API_TOKEN)
	let rsp = await cryptoService.news({
  "symbol": "BTCUSD"
})
	console.log(rsp)
}

getCryptocurrencyNews()
```
## Price

Get the last price for a given crypto ticker


[https://m3o.com/crypto/api#Price](https://m3o.com/crypto/api#Price)

```js
const { CryptoService } = require('m3o/crypto');

// Get the last price for a given crypto ticker
async function getCryptocurrencyPrice() {
	let cryptoService = new CryptoService(process.env.MICRO_API_TOKEN)
	let rsp = await cryptoService.price({
  "symbol": "BTCUSD"
})
	console.log(rsp)
}

getCryptocurrencyPrice()
```
## Quote

Get the last quote for a given crypto ticker


[https://m3o.com/crypto/api#Quote](https://m3o.com/crypto/api#Quote)

```js
const { CryptoService } = require('m3o/crypto');

// Get the last quote for a given crypto ticker
async function getAcryptocurrencyQuote() {
	let cryptoService = new CryptoService(process.env.MICRO_API_TOKEN)
	let rsp = await cryptoService.quote({
  "symbol": "BTCUSD"
})
	console.log(rsp)
}

getAcryptocurrencyQuote()
```
## History

Returns the history for the previous close


[https://m3o.com/crypto/api#History](https://m3o.com/crypto/api#History)

```js
const { CryptoService } = require('m3o/crypto');

// Returns the history for the previous close
async function getPreviousClose() {
	let cryptoService = new CryptoService(process.env.MICRO_API_TOKEN)
	let rsp = await cryptoService.history({
  "symbol": "BTCUSD"
})
	console.log(rsp)
}

getPreviousClose()
```
