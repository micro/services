# URL Proxy

The URL proxy provides the single entrypoint for m3o.one

## Overview

The [url](https://github.com/micro/services) service provides link shortening and sharing. The URL Proxy fronts those urls 
as a single entrypoint at https://m3o.one. We don't serve directly because of time wasted on ssl certs, etc.

## Usage

- Assumes url is of format `https://m3o.one/u/AArfeZE`
- Will call `https://api.m3o.com/url/proxy?shortURL=AArfeZE`
- URL service should return `destinationURL=https://foobar.com/example`
- Proxy will issue a 301 redirect
