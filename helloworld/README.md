# Helloworld Service

This is the Helloworld service

## Overview

An example of how to write a simple helloworld service. This can also be generated using `micro new helloworld`.

## Usage

```
# run the server
micro server

# run the service
micro run github.com/micro/services/helloworld

## call the service
micro call helloworld Helloworld.Call '{"name": "Alice"}'
```
