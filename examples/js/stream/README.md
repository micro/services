# Stream

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Stream/api](https://m3o.com/Stream/api).

Endpoints:

## Publish

Publish a message to the stream. Specify a topic to group messages for a specific topic.


[https://m3o.com/stream/api#Publish](https://m3o.com/stream/api#Publish)

```js
const { StreamService } = require('m3o/stream');

// Publish a message to the stream. Specify a topic to group messages for a specific topic.
async function publishAmessage() {
	let streamService = new StreamService(process.env.MICRO_API_TOKEN)
	let rsp = await streamService.publish({
  "message": {
    "id": "1",
    "type": "signup",
    "user": "john"
  },
  "topic": "events"
})
	console.log(rsp)
}

publishAmessage()
```
## Subscribe

Subscribe to messages for a given topic.


[https://m3o.com/stream/api#Subscribe](https://m3o.com/stream/api#Subscribe)

```js
const { StreamService } = require('m3o/stream');

// Subscribe to messages for a given topic.
async function subscribeToAtopic() {
	let streamService = new StreamService(process.env.MICRO_API_TOKEN)
	let rsp = await streamService.subscribe({
  "topic": "events"
})
	console.log(rsp)
}

subscribeToAtopic()
```
