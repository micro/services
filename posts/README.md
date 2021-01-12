# Post Service

The posts service stores posts

## Usage

### Create a post

```
micro call posts Posts.Save '{"post":{"id":"1","title":"How to Micro","content":"Simply put, Micro is awesome."}}'
micro call posts Posts.Save '{"post":{"id":"2","title":"Fresh posts are fresh","content":"This post is fresher than the How to Micro one"}}'
```

### Create a post with tags

```
micro call posts Posts.Save '{"post":{"id":"3","title":"How to do epic things with Micro","content":"Everything is awesome.","tagNames":["a","b"]}}'
```

### Query posts

```
micro call posts Posts.Query '{}'
micro call posts Posts.Query '{"slug":"how-to-micro"}'
micro call posts Posts.Query '{"offset": 10, "limit": 10}'
```

### Delete posts

```
micro call posts Posts.Delete '{"offset": 10, "limit": 10}'
```
