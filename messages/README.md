# Messages Service

The messages service is a simplified service for sending messages, much like email. You can send a message using the CLI:
```bash
> micro messages send --to=John --from=Barry --subject=HelloWorld --text="Hello John"
```

And then list the messages a user has recieved:
```bash
> micro messages list --user=John
{
	"messages": [
		{
			"id": "78efd836-ca51-4163-af43-65985f7c6587",
			"to": "John",
			"from": "Barry",
			"subject": "HelloWorld",
			"text": "Hello John",
			"sent_at": "1602777240"
		}
	]
}
```

Or lookup an individual email by ID:
```bash
> micro messages read --id=78efd836-ca51-4163-af43-65985f7c6587
{
	"message": {
		"id": "78efd836-ca51-4163-af43-65985f7c6587",
		"to": "John",
		"from": "Barry",
		"subject": "HelloWorld",
		"text": "Hello John",
		"sent_at": "1602777240"
	}
}
```