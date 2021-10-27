# Url

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Url/api](https://m3o.com/Url/api).

Endpoints:

## List

List information on all the shortened URLs that you have created


[https://m3o.com/url/api#List](https://m3o.com/url/api#List)

```js
const { UrlService } = require('m3o/url');

// List information on all the shortened URLs that you have created
async function listYourShortenedUrls() {
	let urlService = new UrlService(process.env.MICRO_API_TOKEN)
	let rsp = await urlService.list({})
	console.log(rsp)
}

listYourShortenedUrls()
```
## Shorten

Shortens a destination URL and returns a full short URL.


[https://m3o.com/url/api#Shorten](https://m3o.com/url/api#Shorten)

```js
const { UrlService } = require('m3o/url');

// Shortens a destination URL and returns a full short URL.
async function shortenAlongUrl() {
	let urlService = new UrlService(process.env.MICRO_API_TOKEN)
	let rsp = await urlService.shorten({
  "destinationURL": "https://mysite.com/this-is-a-rather-long-web-address"
})
	console.log(rsp)
}

shortenAlongUrl()
```
## Proxy

Proxy returns the destination URL of a short URL.


[https://m3o.com/url/api#Proxy](https://m3o.com/url/api#Proxy)

```js
const { UrlService } = require('m3o/url');

// Proxy returns the destination URL of a short URL.
async function resolveAshortUrlToAlongDestinationUrl() {
	let urlService = new UrlService(process.env.MICRO_API_TOKEN)
	let rsp = await urlService.proxy({
  "shortURL": "https://m3o.one/u/ck6SGVkYp"
})
	console.log(rsp)
}

resolveAshortUrlToAlongDestinationUrl()
```
