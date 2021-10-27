# Thumbnail

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Thumbnail/api](https://m3o.com/Thumbnail/api).

Endpoints:

## Screenshot

Create a thumbnail screenshot by passing in a url, height and width


[https://m3o.com/thumbnail/api#Screenshot](https://m3o.com/thumbnail/api#Screenshot)

```js
const { ThumbnailService } = require('m3o/thumbnail');

// Create a thumbnail screenshot by passing in a url, height and width
async function takeScreenshotOfAurl() {
	let thumbnailService = new ThumbnailService(process.env.MICRO_API_TOKEN)
	let rsp = await thumbnailService.screenshot({
  "height": 600,
  "url": "https://m3o.com",
  "width": 600
})
	console.log(rsp)
}

takeScreenshotOfAurl()
```
