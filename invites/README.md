# Invites Service

The invites services allows you to create and manage invites. Example usage:

```bash
> micro invites create --group_id=myawesomegroup --email=john@doe.com
{
	"invite": {
		"id": "fb3a3552-3c7b-4a18-a1f8-08ab56940862",
		"group_id": "myawesomegroup",
		"email": "john@doe.com",
		"code": "86285587"
	}
}

> micro invites list --group_id=fb3a3552-3c7b-4a18-a1f8-08ab56940862
{
	"invites": [
    {
    "id": "fb3a3552-3c7b-4a18-a1f8-08ab56940862",
    "group_id": "myawesomegroup",
    "email": "john@doe.com",
    "code": "86285587"
    }
  ]
}

> micro invites read --code=86285587
{
	"invite": {
		"id": "fb3a3552-3c7b-4a18-a1f8-08ab56940862",
		"group_id": "myawesomegroup",
		"email": "john@doe.com",
		"code": "86285587"
	}
}

> micro invites delete --id=fb3a3552-3c7b-4a18-a1f8-08ab56940862
{}
```
