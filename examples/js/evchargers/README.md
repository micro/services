# Evchargers

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Evchargers/api](https://m3o.com/Evchargers/api).

Endpoints:

## Search

Search by giving a coordinate and a max distance, or bounding box and optional filters


[https://m3o.com/evchargers/api#Search](https://m3o.com/evchargers/api#Search)

```js
const { EvchargersService } = require('m3o/evchargers');

// Search by giving a coordinate and a max distance, or bounding box and optional filters
async function searchByLocation() {
	let evchargersService = new EvchargersService(process.env.MICRO_API_TOKEN)
	let rsp = await evchargersService.search({
  "distance": 2000,
  "location": {
    "latitude": 51.53336351319885,
    "longitude": -0.0252
  },
  "max_results": 2
})
	console.log(rsp)
}

searchByLocation()
```
## Search

Search by giving a coordinate and a max distance, or bounding box and optional filters


[https://m3o.com/evchargers/api#Search](https://m3o.com/evchargers/api#Search)

```js
const { EvchargersService } = require('m3o/evchargers');

// Search by giving a coordinate and a max distance, or bounding box and optional filters
async function searchByBoundingBox() {
	let evchargersService = new EvchargersService(process.env.MICRO_API_TOKEN)
	let rsp = await evchargersService.search({
  "box": {
    "bottom_left": {
      "latitude": 51.52627543859447,
      "longitude": -0.03635349400295168
    },
    "top_right": {
      "latitude": 51.56717121807993,
      "longitude": -0.002293530559768285
    }
  },
  "max_results": 2
})
	console.log(rsp)
}

searchByBoundingBox()
```
## Search

Search by giving a coordinate and a max distance, or bounding box and optional filters


[https://m3o.com/evchargers/api#Search](https://m3o.com/evchargers/api#Search)

```js
const { EvchargersService } = require('m3o/evchargers');

// Search by giving a coordinate and a max distance, or bounding box and optional filters
async function searchWithFiltersFastChargersOnly() {
	let evchargersService = new EvchargersService(process.env.MICRO_API_TOKEN)
	let rsp = await evchargersService.search({
  "distance": 2000,
  "levels": [
    "3"
  ],
  "location": {
    "latitude": 51.53336351319885,
    "longitude": -0.0252
  },
  "max_results": 2
})
	console.log(rsp)
}

searchWithFiltersFastChargersOnly()
```
## ReferenceData

Retrieve reference data as used by this API and in conjunction with the Search endpoint


[https://m3o.com/evchargers/api#ReferenceData](https://m3o.com/evchargers/api#ReferenceData)

```js
const { EvchargersService } = require('m3o/evchargers');

// Retrieve reference data as used by this API and in conjunction with the Search endpoint
async function getReferenceData() {
	let evchargersService = new EvchargersService(process.env.MICRO_API_TOKEN)
	let rsp = await evchargersService.referenceData({})
	console.log(rsp)
}

getReferenceData()
```
