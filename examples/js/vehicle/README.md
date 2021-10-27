# Vehicle

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Vehicle/api](https://m3o.com/Vehicle/api).

Endpoints:

## Lookup

Lookup a UK vehicle by it's registration number


[https://m3o.com/vehicle/api#Lookup](https://m3o.com/vehicle/api#Lookup)

```js
const { VehicleService } = require('m3o/vehicle');

// Lookup a UK vehicle by it's registration number
async function lookupVehicle() {
	let vehicleService = new VehicleService(process.env.MICRO_API_TOKEN)
	let rsp = await vehicleService.lookup({
  "registration": "LC60OTA"
})
	console.log(rsp)
}

lookupVehicle()
```
