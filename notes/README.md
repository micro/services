# Notes Service

Notes service is an RPC service which offers CRUD for notes. It demonstrates usage of the store, errors and logger pacakges. Example usage:

Create a note

```bash
micro notes create --title="HelloWorld" --text="MyFirstNote"
{
	"id": "6d3fa5c0-6e79-4418-a72a-c1650efb65d2"
}
```

Update a note

```bash
micro notes update --id=6d3fa5c0-6e79-4418-a72a-c1650efb65d2 --title="HelloWorld" --text="MyFirstNote (v2)"
{}
```

List notes

```bash
micro notes list
{
	"notes": [
		{
			"id": "6d3fa5c0-6e79-4418-a72a-c1650efb65d2",
			"created": "1602849877",
			"title": "HelloWorld",
			"text": "MyFirstNote (v2)"
		}
	]
}
```

Delete a note

```bash
micro notes delete --id=6d3fa5c0-6e79-4418-a72a-c1650efb65d2
{}
```
