# Micro Services [![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/micro/services?tab=doc) [![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Reusable real world Micro services.

## Overview

Micro services provide the fundamental building blocks for any products, apps or services. They can be used in isolation 
or combined to create a powerful distributed system. The services are intended to be consumed by each other using RPC 
and from the external world through a Micro API.

## Usage

Run a service from source

```
micro run github.com/micro/services/helloworld
```

To call a service from another

```
import "github.com/micro/services/helloworld/proto"
```

Call it from the api

```
curl http://localhost:8080/helloworld
```

## Contribute

We welcome contributions of additional services:

- Services must be built using the Micro platform
- Any dependency must be configured using the Micro Config
- All services must include a README.md and be well commented
