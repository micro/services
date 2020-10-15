# PubSub

A simple example of using PubSub

## Overview

This example is a simple pubsub service which subscribes to a topic and periodically publishes 
to the same place. The interface for how it works may change but its simple enough.

## Usage

```
# start micro
micro server

# run pubsub
micro run pubsub

# check status
micro status

# get the logs
micro logs pubsub
```
