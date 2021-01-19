The messages service is a simplified service for sending messages, much like email.

# Messages Service

## Send a message

### CLI

```bash
> micro messages send --to=John --from=Barry --subject=HelloWorld --text="Hello John"
```

## List the messages a user has received

### CLI

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

## Lookup an individual email by ID

### CLI

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