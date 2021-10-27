# Helloworld

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Helloworld/api](https://m3o.com/Helloworld/api).

Endpoints:

## Call

Call returns a personalised "Hello $name" response


[https://m3o.com/helloworld/api#Call](https://m3o.com/helloworld/api#Call)

```js
const { HelloworldService } = require('m3o/helloworld');

// Call returns a personalised "Hello $name" response
async function callTheHelloworldService() {
	let helloworldService = new HelloworldService(process.env.MICRO_API_TOKEN)
	let rsp = await helloworldService.call({
  "name": "John"
})
	console.log(rsp)
}

callTheHelloworldService()
```
## Stream

Stream returns a stream of "Hello $name" responses


[https://m3o.com/helloworld/api#Stream](https://m3o.com/helloworld/api#Stream)

```js
const { HelloworldService } = require('m3o/helloworld');

// Stream returns a stream of "Hello $name" responses
async function streamsAreCurrentlyTemporarilyNotSupportedInClients() {
	let helloworldService = new HelloworldService(process.env.MICRO_API_TOKEN)
	let rsp = await helloworldService.stream({
  "name": "not supported"
})
	console.log(rsp)
}

streamsAreCurrentlyTemporarilyNotSupportedInClients()
```
