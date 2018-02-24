A toy chat app
==============

## Overview

The app contains 2 parts:

- Chat server: Responsible for receiving & sending messages between clients.

- Archive daemon: Listen for Archive channel & save messages to Mongodb.

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
       |            _________________
       |           |     Archive     |             .---------.
       |---------->|   NSQ channel   |------------>| Mongodb |
                   |_________________|             |_________|

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

3. Build

```sh
go get github.com/manhtai/golang-nsq-chat
dep ensure
go build ./pkg/cmd/chat
go build ./pkg/cmd/archive
```

4. Run

- Chat server

```sh
./chat
```

- Archive daemon

```sh
./archive
```

## Generate cert.pem & key.pem

```sh
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem
```

## Live reload for Chat server

```
go get https://github.com/Unknwon/bra
Bra run
```
