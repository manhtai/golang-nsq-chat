A toy chat app
==============

## Overview

The app contains 3 parts:

- Persistent structs: User, Channel, Message

> These models holds information about User, Channel & Message from channels
> when user send messages, data is saved to Mongodb.

- Websocket structs: Client, Room

> These are responsible for opening Websocket connection to receive messages
> from user, send to NSQ, get messages from NSQ and broadcast to all clients
> in a specific channel.

- NSQ structs: keep track of NSQ


## Get started

1. Start Mongodb

```sh
mongod
```

2. Start nsq

```sh
nsqd
nsqlookupd
```

3. Start chat server

```sh
go get github.com/manhtai/golang-nsq-chat
dep ensure
go run main.go
```

## SSL & Live reload support

```sh
go get github.com/codegangsta/gin
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem
gin --certFile cert.pem --keyFile key.pem --all main.go
```
