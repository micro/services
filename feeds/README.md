# Feeds Service

This is the Feeds service

Generated with

```
micro new feeds
```

## Usage


```
micro feeds new --name="a16z" --url=http://a16z.com/feed/
```

```
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
