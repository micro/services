# Twitter Service

This is the Twitter service. It let's you write tweets via the service using the twitter api.

## Usage

Create a new app in the twitter [developer portal](https://developer.twitter.com/en/apps).

Then set the config for the api token, secret, consumer key and consumer secret

```
micro config set twitter.api_token xxxx
micro config set twitter.api_token_secret xxxx
micro config set twitter.consumer_key xxxx
micro config set twitter.consumer_secret xxxx
```

Now tweet stuff

```
micro twitter api tweet --status "Tweeting via a micro service"
```
