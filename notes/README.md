# Notes Service

Notes service is an RPC service which offers CRUD for notes. It demonstrates usage of the store and errors pacakges. Example usage:

```bash
> micro notes create --title="HelloWorld" --text="MyFirstNote"
{
	"id": "6d3fa5c0-6e79-4418-a72a-c1650efb65d2"
}
> micro notes update --id=6d3fa5c0-6e79-4418-a72a-c1650efb65d2 --title="HelloWorld" --text="MyFirstNote (v2)"
{}
> micro notes list
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
> micro notes delete --id=6d3fa5c0-6e79-4418-a72a-c1650efb65d2
{}
```