# Time

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Time/api](https://m3o.com/Time/api).

Endpoints:

## Now

Get the current time


[https://m3o.com/time/api#Now](https://m3o.com/time/api#Now)

```js
const { TimeService } = require('m3o/time');

// Get the current time
async function returnsCurrentTimeOptionallyWithLocation() {
	let timeService = new TimeService(process.env.MICRO_API_TOKEN)
	let rsp = await timeService.now({})
	console.log(rsp)
}

returnsCurrentTimeOptionallyWithLocation()
```
## Zone

Get the timezone info for a specific location


[https://m3o.com/time/api#Zone](https://m3o.com/time/api#Zone)

```js
const { TimeService } = require('m3o/time');

// Get the timezone info for a specific location
async function getTheTimezoneInfoForAspecificLocation() {
	let timeService = new TimeService(process.env.MICRO_API_TOKEN)
	let rsp = await timeService.zone({
  "location": "London"
})
	console.log(rsp)
}

getTheTimezoneInfoForAspecificLocation()
```
