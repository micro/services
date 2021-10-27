# Holidays

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Holidays/api](https://m3o.com/Holidays/api).

Endpoints:

## List

List the holiday dates for a given country and year


[https://m3o.com/holidays/api#List](https://m3o.com/holidays/api#List)

```js
const { HolidaysService } = require('m3o/holidays');

// List the holiday dates for a given country and year
async function getHolidays() {
	let holidaysService = new HolidaysService(process.env.MICRO_API_TOKEN)
	let rsp = await holidaysService.list({
  "country_code": "GB",
  "year": 2022
})
	console.log(rsp)
}

getHolidays()
```
## Countries

Get the list of countries that are supported by this API


[https://m3o.com/holidays/api#Countries](https://m3o.com/holidays/api#Countries)

```js
const { HolidaysService } = require('m3o/holidays');

// Get the list of countries that are supported by this API
async function listCountries() {
	let holidaysService = new HolidaysService(process.env.MICRO_API_TOKEN)
	let rsp = await holidaysService.countries({})
	console.log(rsp)
}

listCountries()
```
