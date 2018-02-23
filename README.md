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
> in a specific Room.

- NSQ structs: NsqReader

> This struct keep track of NSQ consumer corresponding to each Room in one NSQ
> channel.

```
 _____________
|   "Chat"    |                  Pub to "Chat" topic
|    Topic    |<-------------------------------------------------------.
|_____________|                                                        |
       |                                                               |
       |            _________________              .----> Client 1 --->|
       |           |     Server1     | Go channels |                   |
       |---------->|   NSQ channel   |------------>|----> Client 2 --->|
       |           |_________________|             |                   |
       |                                           '----> Client 3 --->|
       |            _________________
       |           |     Server2     | Go channels
       |---------->|   NSQ channel   |------------> ...
       |           |_________________|
       |
       |
       ...

```

## Get started

1. Start Mongodb

```sh
mongod
```

2. Start nsq

```sh
nsqlookupd
nsqd -lookupd-tcp-address=0.0.0.0:4160
```

Export `NSQLOOKUPD_HTTP_ADDRESS` and `NSQD_HTTP_ADDRESS` to corresponding address.

3. Start chat server

```sh
go get github.com/manhtai/golang-nsq-chat
dep ensure
go build ./pkg/cmd/chat
./chat -cert-file=cert.pem -key-file=key.pem
```

## Generate cert.pem & key.pem

```sh
go get github.com/codegangsta/gin
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem
```

## Live-reload

```
go get https://github.com/Unknwon/bra
Bra run
```
