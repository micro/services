Tag any resource by savin a tag associated with their ID in the tag service.

# Tag Service

## Query tags

```
micro tags list --type=post-tag
```

Generated with

```
micro new --namespace=go.micro --type=service tag
```

## Getting Started

- [Tag Service](#tag-service)
  - [Query tags](#query-tags)
  - [Getting Started](#getting-started)
  - [Configuration](#configuration)
  - [Dependencies](#dependencies)
  - [Usage](#usage)

## Configuration

- FQDN: go.micro.service.tag
- Type: service
- Alias: tag

## Dependencies

Micro services depend on service discovery. The default is multicast DNS, a zeroconf system.

In the event you need a resilient multi-host setup we recommend etcd.

```
# install etcd
brew install etcd

# run etcd
etcd
```

## Usage

A Makefile is included for convenience

Build the binary

```
make build
```

Run the service
```
./tag-service
```

Build a docker image
```
make docker
```