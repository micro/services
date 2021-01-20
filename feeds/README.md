A single uniform API for crawling and indexing RSS feeds

# Feeds Service

Designed to populate the posts service with RSS feeds from other blogs. Useful for migration or just to get outside content into the posts service.

## Creating a feeed

### cURL

```bash
> curl 'https://api.m3o.com/feeds/New' \
  -H 'micro-namespace: $yourNamespace' \
  -H 'authorization: Bearer $yourToken' \
  -d '{"name":"a16z", "url": "http://a16z.com/feed/"}';
{}
```

### CLI

```shell
micro feeds new --name="a16z" --url=http://a16z.com/feed/
```

## Querying feeded posts

```shell
$ micro posts query
{
	"posts": [
		{
			"id": "39cdfbd6e7534bcd868be9eebbf43f8f",
			"title": "Anthony Albanese: From the NYSE to Crypto",
			"slug": "anthony-albanese-from-the-nyse-to-crypto",
			"created": "1605104742",
			"updated": "1605105364",
			"metadata": {
				"domain": "a16z.com",
				"link": "https://a16z.com/2020/10/28/anthony-albanese-from-the-nyse-to-crypto/"
			}
		},
		{
			"id": "5e9285c01311704e204322ba564cd99e",
			"title": "Journal Club: From Insect Eyes to Nanomaterials",
			"slug": "journal-club-from-insect-eyes-to-nanomaterials",
			"created": "1605104741",
			"updated": "1605105363",
			"metadata": {
				"domain": "a16z.com",
				"link": "https://a16z.com/2020/10/29/journal-club-insect-eyes-nanomaterials/"
			}
		},
	]
}
```

```
make proto
```

Run the service

```
micro run .
```
