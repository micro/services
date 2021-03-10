The mail service is a simplified service for sending mail, much like email.

# Mail Service

## Send a message

### CLI

```bash
> micro mail send --to=John --from=Barry --subject=HelloWorld --text="Hello John"
```

## List the mail a user has received

### CLI

```bash
> micro mail list --user=John
{
	"mail": [
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
> micro mail read --id=78efd836-ca51-4163-af43-65985f7c6587
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