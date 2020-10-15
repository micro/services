# Blog

This is a full end to end example of writing a multi-service blog application

## Usage

Check out the [blog tutorial](https://m3o.dev/tutorials/building-a-blog) on the developer docs.

## How it works

### Present

The blog services are designed so a user can deploy them to their own micro namespace, write content with their Micro account with commands like

```sh
micro posts save --id=7 --tags=News,Finance --title="Breaking News" --content="The stock market has just crashed"
```

and display content on their frontend by consuming the API:

```sh
curl -H "Authorization: Bearer $MICRO_API_TOKEN" "Micro-Namespace: $NAMESPACE" https://api.m3o.com/tags/list


{
	"tags": [
		{
			"type": "post-tag",
			"slug": "news",
			"title": "News",
			"count": "3"
		}
    ]
]
```

There are no comments provided yet, just posts and tags.
Access is governed by auth rules, ie. Posts List, Tags List is open, Posts Save requires a Micro login.

### Future possibilities

#### Enable non Micro users to write posts, comments

If we provide a user/login service (markedly different from auth, it can be a simple session based auth) to enable non Micro users to register, the following can be done:

- A user (let's call the user Alice from this point) launches posts, tags, login service in their namespace.
- Alice opens up said endpoints
- People (let's call them Yoga Pants Co and Drone Inc) hosting JS and HTML on Netlify or Github Pages could create accounts in the services hosted by Alice. In this way, Alice, by having a Micro account becomes a headless CMS provider. Multiple blogs can be created on top of Alice's service instances.

Questions:
- How will Yoga Pants Co or Drone Inc pay Alice or M3O for the costs of their backend hosting?