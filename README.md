# Micro Services [![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/micro/services?tab=doc) [![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Programmable real world Micro services.

## Overview

Micro services provide the fundamental building blocks for any products, apps or services. They can be used in isolation 
or combined to create a powerful distributed system. The services are intended to be consumed by each other using RPC 
and from the external world through a Micro API.

## Cloud

Find cloud hosted services on [m3o.com](https://m3o.com).

## Usage

Run a service from source

```
micro run github.com/micro/services/helloworld
```

To call a service from another

```
import "github.com/micro/services/helloworld/proto"
```

## Clients

API clients are generated in the [clients](https://github.com/micro/services/tree/master/clients) directory.

To call a service via the api client import as follows

```
import "github.com/micro/services/clients/go/helloworld"
```

## Examples

See the [examples](https://github.com/micro/services/tree/master/examples) directory.

## Contribute

We welcome contributions of additional services which are then hosted on [m3o.com](https://m3o.com).

- Services must be built using the Micro platform
- Any dependency must be configured using the Micro Config
- All services to be published must include a `publicapi.json` file

## License

[Apache 2.0](LICENSE)
