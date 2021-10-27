# Gifs

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Gifs/api](https://m3o.com/Gifs/api).

Endpoints:

## Search

Search for a GIF


[https://m3o.com/gifs/api#Search](https://m3o.com/gifs/api#Search)

```js
const { GifsService } = require('m3o/gifs');

// Search for a GIF
async function search() {
	let gifsService = new GifsService(process.env.MICRO_API_TOKEN)
	let rsp = await gifsService.search({
  "limit": 2,
  "query": "dogs"
})
	console.log(rsp)
}

search()
```
