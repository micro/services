# Rss

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Rss/api](https://m3o.com/Rss/api).

Endpoints:

## Feed

Get an RSS feed by name. If no name is given, all feeds are returned. Default limit is 25 entries.


[https://m3o.com/rss/api#Feed](https://m3o.com/rss/api#Feed)

```js
const { RssService } = require('m3o/rss');

// Get an RSS feed by name. If no name is given, all feeds are returned. Default limit is 25 entries.
async function readAfeed() {
	let rssService = new RssService(process.env.MICRO_API_TOKEN)
	let rsp = await rssService.feed({
  "name": "bbc"
})
	console.log(rsp)
}

readAfeed()
```
## List

List the saved RSS fields


[https://m3o.com/rss/api#List](https://m3o.com/rss/api#List)

```js
const { RssService } = require('m3o/rss');

// List the saved RSS fields
async function listRssFeeds() {
	let rssService = new RssService(process.env.MICRO_API_TOKEN)
	let rsp = await rssService.list({})
	console.log(rsp)
}

listRssFeeds()
```
## Remove

Remove an RSS feed by name


[https://m3o.com/rss/api#Remove](https://m3o.com/rss/api#Remove)

```js
const { RssService } = require('m3o/rss');

// Remove an RSS feed by name
async function removeAfeed() {
	let rssService = new RssService(process.env.MICRO_API_TOKEN)
	let rsp = await rssService.remove({
  "name": "bbc"
})
	console.log(rsp)
}

removeAfeed()
```
## Add

Add a new RSS feed with a name, url, and category


[https://m3o.com/rss/api#Add](https://m3o.com/rss/api#Add)

```js
const { RssService } = require('m3o/rss');

// Add a new RSS feed with a name, url, and category
async function addAnewFeed() {
	let rssService = new RssService(process.env.MICRO_API_TOKEN)
	let rsp = await rssService.add({
  "category": "news",
  "name": "bbc",
  "url": "http://feeds.bbci.co.uk/news/rss.xml"
})
	console.log(rsp)
}

addAnewFeed()
```
