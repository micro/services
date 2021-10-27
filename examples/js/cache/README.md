# Cache

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Cache/api](https://m3o.com/Cache/api).

Endpoints:

## Set

Set an item in the cache. Overwrites any existing value already set.


[https://m3o.com/cache/api#Set](https://m3o.com/cache/api#Set)

```js
const { CacheService } = require('m3o/cache');

// Set an item in the cache. Overwrites any existing value already set.
async function setAvalue() {
	let cacheService = new CacheService(process.env.MICRO_API_TOKEN)
	let rsp = await cacheService.set({
  "key": "foo",
  "value": "bar"
})
	console.log(rsp)
}

setAvalue()
```
## Get

Get an item from the cache by key


[https://m3o.com/cache/api#Get](https://m3o.com/cache/api#Get)

```js
const { CacheService } = require('m3o/cache');

// Get an item from the cache by key
async function getAvalue() {
	let cacheService = new CacheService(process.env.MICRO_API_TOKEN)
	let rsp = await cacheService.get({
  "key": "foo"
})
	console.log(rsp)
}

getAvalue()
```
## Delete

Delete a value from the cache


[https://m3o.com/cache/api#Delete](https://m3o.com/cache/api#Delete)

```js
const { CacheService } = require('m3o/cache');

// Delete a value from the cache
async function deleteAvalue() {
	let cacheService = new CacheService(process.env.MICRO_API_TOKEN)
	let rsp = await cacheService.delete({
  "key": "foo"
})
	console.log(rsp)
}

deleteAvalue()
```
## Increment

Increment a value (if it's a number)


[https://m3o.com/cache/api#Increment](https://m3o.com/cache/api#Increment)

```js
const { CacheService } = require('m3o/cache');

// Increment a value (if it's a number)
async function incrementAvalue() {
	let cacheService = new CacheService(process.env.MICRO_API_TOKEN)
	let rsp = await cacheService.increment({
  "key": "counter",
  "value": 2
})
	console.log(rsp)
}

incrementAvalue()
```
## Decrement

Decrement a value (if it's a number)


[https://m3o.com/cache/api#Decrement](https://m3o.com/cache/api#Decrement)

```js
const { CacheService } = require('m3o/cache');

// Decrement a value (if it's a number)
async function decrementAvalue() {
	let cacheService = new CacheService(process.env.MICRO_API_TOKEN)
	let rsp = await cacheService.decrement({
  "key": "counter",
  "value": 2
})
	console.log(rsp)
}

decrementAvalue()
```
