# Id

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Id/api](https://m3o.com/Id/api).

Endpoints:

## Generate

Generate a unique ID. Defaults to uuid.


[https://m3o.com/id/api#Generate](https://m3o.com/id/api#Generate)

```js
const { IdService } = require('m3o/id');

// Generate a unique ID. Defaults to uuid.
async function generateAuniqueId() {
	let idService = new IdService(process.env.MICRO_API_TOKEN)
	let rsp = await idService.generate({
  "type": "uuid"
})
	console.log(rsp)
}

generateAuniqueId()
```
## Generate

Generate a unique ID. Defaults to uuid.


[https://m3o.com/id/api#Generate](https://m3o.com/id/api#Generate)

```js
const { IdService } = require('m3o/id');

// Generate a unique ID. Defaults to uuid.
async function generateAshortId() {
	let idService = new IdService(process.env.MICRO_API_TOKEN)
	let rsp = await idService.generate({
  "type": "shortid"
})
	console.log(rsp)
}

generateAshortId()
```
## Generate

Generate a unique ID. Defaults to uuid.


[https://m3o.com/id/api#Generate](https://m3o.com/id/api#Generate)

```js
const { IdService } = require('m3o/id');

// Generate a unique ID. Defaults to uuid.
async function generateAsnowflakeId() {
	let idService = new IdService(process.env.MICRO_API_TOKEN)
	let rsp = await idService.generate({
  "type": "snowflake"
})
	console.log(rsp)
}

generateAsnowflakeId()
```
## Generate

Generate a unique ID. Defaults to uuid.


[https://m3o.com/id/api#Generate](https://m3o.com/id/api#Generate)

```js
const { IdService } = require('m3o/id');

// Generate a unique ID. Defaults to uuid.
async function generateAbigflakeId() {
	let idService = new IdService(process.env.MICRO_API_TOKEN)
	let rsp = await idService.generate({
  "type": "bigflake"
})
	console.log(rsp)
}

generateAbigflakeId()
```
## Types

List the types of IDs available. No query params needed.


[https://m3o.com/id/api#Types](https://m3o.com/id/api#Types)

```js
const { IdService } = require('m3o/id');

// List the types of IDs available. No query params needed.
async function listTheTypesOfIdsAvailable() {
	let idService = new IdService(process.env.MICRO_API_TOKEN)
	let rsp = await idService.types({})
	console.log(rsp)
}

listTheTypesOfIdsAvailable()
```
