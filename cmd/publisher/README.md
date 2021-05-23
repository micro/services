# API Publisher

The public api publisher for Micro services

## Overview

The API publisher assumes a few things:

- Your service name is the directory of the service
- You have a `README.md` file we take as description
- You have a `make api` command to generate a `api-{service}.json` for the OpenAPI spec
- You optionally have a `publicapi.json` file to define extra info such a category, icon, etc
- You optionally have a `examples.json` file to separately define usage examples
- You optionaly have a `pricing.json` file to separately define pricing information

All these are combined to produce a Public API.

## Readme

The readmes are taken verbatim. Everything before the newline is used as a short excerpt. Examples are appended to your API 
tab if they exist so no need to add random examples, curls, etc to the readme. Focus on describing the service in brevity.

## Comments

Some rules on how to write protos so they nicely appear in the output of this script:

- The request types (eg. `LoginRequest`) comments will be taken and used as a description for the endpoint (eg. `Login`) itself. This might change.
- The proto message field comments will be taken and displayed to craft them with care

To provide example values use the following format:

```shell
// rss feed name
// eg. a16z
string name = 1;
```

The part after the `eg. ` until the newline will be used as example value.
