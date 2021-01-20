Real time messaging API which enables Chat services to be embedded anywhere

# Chat Service

The Chat service is a programmable instant messaging API service which can be used in any application to immediately create conversations. 

## Create a chat

### cURL

```bash
> curl 'https://api.m3o.com/chat/New' \
  -H 'micro-namespace: $yourNamespace' \
  -H 'authorization: Bearer $yourToken' \
  -d '{"user_ids":["JohnBarry"]}';
{
	"chat_id": "3c9ea66c-d516-45d4-abe8-082089e18b27"
}
```

### CLI

```bash
> micro chat new --user_ids=JohnBarry
{
	"chat_id": "3c9ea66c-d516-45d4-abe8-082089e18b27"
}
```

## Send a message to the chat

### cURL

```bash
> curl 'https://api.m3o.com/chat/Send' \
  -H 'micro-namespace: $yourNamespace' \
  -H 'authorization: Bearer $yourToken' \
  -d '{"user_id": "John", "subject": "Hello", "text": "Hey Barry"}';
{}
```

### CLI

```bash
> micro chat send --chat_id=bed4f0f0-da12-46d2-90d2-17ae1714a214 --user_id=John --subject=Hello --text='Hey Barry'
{}
```

## View the chat history

### cURL

```bash
> curl 'https://api.m3o.com/chat/Send' \
  -H 'micro-namespace: $yourNamespace' \
  -H 'authorization: Bearer $yourToken' \
  -d '{"chat_id": "bed4f0f0-da12-46d2-90d2-17ae1714a214"}';
{
	"messages": [
		{
			"id": "a61284a8-f471-4734-9192-640d89762e98",
			"client_id": "6ba0d2a6-96fa-47d8-8f6f-7f75b4cc8b3e",
			"chat_id": "bed4f0f0-da12-46d2-90d2-17ae1714a214",
			"user_id": "John",
			"subject": "Hello",
			"text": "Hey Barry"
		}
	]
}
```

### CLI
```bash
> micro chat history --chat_id=bed4f0f0-da12-46d2-90d2-17ae1714a214
{
	"messages": [
		{
			"id": "a61284a8-f471-4734-9192-640d89762e98",
			"client_id": "6ba0d2a6-96fa-47d8-8f6f-7f75b4cc8b3e",
			"chat_id": "bed4f0f0-da12-46d2-90d2-17ae1714a214",
			"user_id": "John",
			"subject": "Hello",
			"text": "Hey Barry"
		}
	]
}
```
