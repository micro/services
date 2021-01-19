The post service is responsible for storing and querying posts by their slugs or IDs. Posts support tags and metadata, for details see `posts.proto`.

# Post Service

## Create a post

### cURL

```shell
> curl 'https://api.m3o.com/posts/Save' \
  -H 'micro-namespace: $yourNamespace' \
  -H 'authorization: Bearer $yourToken' \
  -d '{"post":{"id":"1","title":"How to Micro","content":"Simply put, Micro is awesome."}}';

> curl 'https://api.m3o.com/chat/Save' \
  -H 'micro-namespace: $yourNamespace' \
  -H 'authorization: Bearer $yourToken' \
  -d '{"post":{"id":"2","title":"Fresh posts are fresh","content":"This post is fresher than the How to Micro one"}}';
```

### CLI

```shell
micro call posts Posts.Save '{"post":{"id":"1","title":"How to Micro","content":"Simply put, Micro is awesome."}}'

micro call posts Posts.Save '{"post":{"id":"2","title":"Fresh posts are fresh","content":"This post is fresher than the How to Micro one"}}'
```

## Create a post with tags


### cURL

```shell
> curl 'https://api.m3o.com/posts/Save' \
  -H 'micro-namespace: $yourNamespace' \
  -H 'authorization: Bearer $yourToken' \
  -d '{"post":{"id":"3","title":"How to do epic things with Micro","content":"Everything is awesome.","tagNames":["a","b"]}}';
```

### CLI

```shell
micro call posts Posts.Save '{"post":{"id":"3","title":"How to do epic things with Micro","content":"Everything is awesome.","tagNames":["a","b"]}}'
```

## Query posts

### cURL

```shell
# Query all
> curl 'https://api.m3o.com/posts/Query' \
  -H 'micro-namespace: $yourNamespace' \
  -H 'authorization: Bearer $yourToken' \
  -d '{}';

# Query by slug
> curl 'https://api.m3o.com/posts/Query' \
  -H 'micro-namespace: $yourNamespace' \
  -H 'authorization: Bearer $yourToken' \
  -d '{"slug":"how-to-micro"}';

# Limit and offset
> curl 'https://api.m3o.com/posts/Query' \
  -H 'micro-namespace: $yourNamespace' \
  -H 'authorization: Bearer $yourToken' \
  -d '{"offset": 10, "limit": 10}';
```

### CLI

```shell
micro call posts Posts.Query '{}'
micro call posts Posts.Query '{"slug":"how-to-micro"}'
micro call posts Posts.Query '{"offset": 10, "limit": 10}'
```

## Delete posts

```shell
> curl 'https://api.m3o.com/posts/Delete' \
  -H 'micro-namespace: $yourNamespace' \
  -H 'authorization: Bearer $yourToken' \
  -d '{"id": "3c9ea66c"}';
```

```shell
micro call posts Posts.Delete '{"id": "3c9ea66c"}'
```
