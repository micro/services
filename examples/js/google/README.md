# Google

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Google/api](https://m3o.com/Google/api).

Endpoints:

## Search

Search for videos on Google


[https://m3o.com/google/api#Search](https://m3o.com/google/api#Search)

```js
const { GoogleService } = require('m3o/google');

// Search for videos on Google
async function searchForVideos() {
	let googleService = new GoogleService(process.env.MICRO_API_TOKEN)
	let rsp = await googleService.search({
  "query": "how to make donuts"
})
	console.log(rsp)
}

searchForVideos()
```
