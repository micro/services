# Location

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Location/api](https://m3o.com/Location/api).

Endpoints:

## Read

Read an entity by its ID


[https://m3o.com/location/api#Read](https://m3o.com/location/api#Read)

```js
const { LocationService } = require('m3o/location');

// Read an entity by its ID
async function getLocationById() {
	let locationService = new LocationService(process.env.MICRO_API_TOKEN)
	let rsp = await locationService.read({
  "id": "1"
})
	console.log(rsp)
}

getLocationById()
```
## Search

Search for entities in a given radius


[https://m3o.com/location/api#Search](https://m3o.com/location/api#Search)

```js
const { LocationService } = require('m3o/location');

// Search for entities in a given radius
async function searchForLocations() {
	let locationService = new LocationService(process.env.MICRO_API_TOKEN)
	let rsp = await locationService.search({
  "center": {
    "latitude": 51.511061,
    "longitude": -0.120022
  },
  "numEntities": 10,
  "radius": 100,
  "type": "bike"
})
	console.log(rsp)
}

searchForLocations()
```
## Save

Save an entity's current position


[https://m3o.com/location/api#Save](https://m3o.com/location/api#Save)

```js
const { LocationService } = require('m3o/location');

// Save an entity's current position
async function saveAnEntity() {
	let locationService = new LocationService(process.env.MICRO_API_TOKEN)
	let rsp = await locationService.save({
  "entity": {
    "id": "1",
    "location": {
      "latitude": 51.511061,
      "longitude": -0.120022,
      "timestamp": "1622802761"
    },
    "type": "bike"
  }
})
	console.log(rsp)
}

saveAnEntity()
```
