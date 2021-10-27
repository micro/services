# Geocoding

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Geocoding/api](https://m3o.com/Geocoding/api).

Endpoints:

## Lookup

Lookup returns a geocoded address including normalized address and gps coordinates. All fields are optional, provide more to get more accurate results


[https://m3o.com/geocoding/api#Lookup](https://m3o.com/geocoding/api#Lookup)

```js
const { GeocodingService } = require('m3o/geocoding');

// Lookup returns a geocoded address including normalized address and gps coordinates. All fields are optional, provide more to get more accurate results
async function geocodeAnAddress() {
	let geocodingService = new GeocodingService(process.env.MICRO_API_TOKEN)
	let rsp = await geocodingService.lookup({
  "address": "10 russell st",
  "city": "london",
  "country": "uk",
  "postcode": "wc2b"
})
	console.log(rsp)
}

geocodeAnAddress()
```
## Reverse

Reverse lookup an address from gps coordinates


[https://m3o.com/geocoding/api#Reverse](https://m3o.com/geocoding/api#Reverse)

```js
const { GeocodingService } = require('m3o/geocoding');

// Reverse lookup an address from gps coordinates
async function reverseGeocodeLocation() {
	let geocodingService = new GeocodingService(process.env.MICRO_API_TOKEN)
	let rsp = await geocodingService.reverse({
  "latitude": 51.5123064,
  "longitude": -0.1216235
})
	console.log(rsp)
}

reverseGeocodeLocation()
```
