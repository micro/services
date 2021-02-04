# Groups Service

The group serivce is a basic CRUD service for groups. You can use it to create groups, add members and lookup which groups a user is a member of.

Example usage:

```bash
$ micro groups create --name=Micro
{
	"group": {
		"id": "e35562c9-b6f6-459a-b52d-7e6159465fd6",
		"name": "Micro"
	}
}
$ micro groups addMember --group_id=e35562c9-b6f6-459a-b52d-7e6159465fd6 --member_id=Asim
{}
$ micro groups list --member_id=Asim
{
	"groups": [
		{
			"id": "e35562c9-b6f6-459a-b52d-7e6159465fd6",
			"name": "Micro",
			"member_ids": [
				"Asim"
			]
		}
	]
}
$ micro groups list --member_id=Boris
{}
```
