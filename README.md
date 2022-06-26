# Micro Services [![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/micro/services?tab=doc) [![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Reusable real world Micro services.

## Overview

Micro services provide the fundamental building blocks for any products, apps or services. They can be used in isolation 
or combined to create a powerful distributed system. The services are intended to be consumed by each other using RPC 
and from the external world through a Micro API.

## Services

Services available thus far:

- address - Address lookup by postcode
- answer - Instant answers to any question
- cache - Fast access key-value storage
- crypto - Cryptocurrency prices, quotes, and news
- currency - Exchange rates and currency conversion
- db - Serverless postgres database
- email - Send emails in a flash
- emoji - All the emojis you need 🎉
- evchargers - Find electric vehicle (EV) chargers wherever you go 
- event - Event stream processing
- file - Store, list, and retrieve text files
- forex - Foreign exchange (FX) rates
- geocoding - Address geocoding and reverse lookup
- gifs - Quick and simple GIF search
- google - Google search service
- helloworld - Just saying hello world
- holidays - Find the holidays observed in a particular country
- id - Generate unique IDs (uuid, snowflake, etc)
- image - Upload, resize, and convert images
- ip - IP to geolocation lookup
- location - Real time GPS location tracking and search
- notes - Store and retrieve notes
- otp - One time password generation
- postcode - Fast UK postcode lookup
- prayer - Islamic prayer times
- qr - QR code generator
- quran - The Holy Quran
- routing - Etas, routes and turn by turn directions
- rss - RSS feed crawler and reader
- sentiment - Real time sentiment analysis
- sms - Send SMS messages
- stock - Live stock quotes and prices
- sunnah - Traditions and practices of the Islamic prophet, Muhammad (pbuh)
- thumbnail - Create website thumbnails
- time - Time, date, and timezone info
- twitter - Realtime twitter timeline & search
- url - URL shortening, sharing, and tracking
- user - Authentication and user management
- vehicle - UK vehicle lookup
- weather - Real time weather forecast
- youtube - Search for YouTube videos
- mq - PubSub messaging
- stream - Ephemeral message streams
- spam - Check if an email is spam
- news - Get the latest news
- app - Serverless app deployment
- nft - Explore NFT Assets
- space - Infinite cloud storage
- movie - Search for movies
- search - Indexing and full text search
- translate - Language translation service
- function - Serverless functions
- avatar - Generate an avatar
- contact - Store your contacts
- carbon - Purchase carbon offsets
- minecraft - Minecraft server ping
- ping - Ping any URL
- place - Search for places
- chat - Real time messaging
- lists - Make a list
- comments - Add comments to any App
- memegen - Generate funny memes
- password - Generate strong passwords
- bitcoin - Realtime Bitcoin price
- analytics - Track and retrieve events
- tunnel - Tunnel HTTP requests
- price - Global commodities index
- github - # Github Service
- joke - Funny Jokes
- dns - DNS over HTTPS (DoH)
- webhook - # Webhook Service
- twilio - Twilio SMS service
- Wordle - Multiplayer wordle

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
