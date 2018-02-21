package models

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/manhtai/golang-nsq-chat/config"
)

// Room represents a room to chat
type Room struct {
	// join is a channel for clients wishing to join the room.
	join chan *Client
	// leave is a channel for clients wishing to leave the room.
	leave chan *Client
	// clients holds all current clients in this room, decided by join & leave channels
	clients map[*Client]bool
	// subscribe is a channel for subscribing to a NSQ channel
	subscribe chan *Client
	// nsqReaders holds all NsqTopicReader after subscribing
	nsqReaders map[string]*NsqReader
}

// run start a room and run it forever
func run(r *Room) {
	for {
		select {
		case client := <-r.join:
			// joining
			r.clients[client] = true
		case client := <-r.leave:
			// leaving
			delete(r.clients, client)
			for channel, reader := range client.subscribed {
				log.Println("Delete client on channel: " + channel)
				reader.RemoveClient(client)
			}
			close(client.send)
		case client := <-r.subscribe:
			nsqChannelName := client.channel
			reader, readerExists := r.nsqReaders[nsqChannelName]

			if !readerExists {
				var err error
				reader, err = NewNsqReader(nsqChannelName)
				if err != nil {
					log.Printf("Failed to subscribe to channel: '%s'",
						nsqChannelName)
					break
				}
			}
			client.subscribed[nsqChannelName] = reader
			r.nsqReaders[nsqChannelName] = reader

			reader.AddClient(client, nsqChannelName)
		}
	}
}

// NewRoomChan creates a new room for clients to join
func NewRoomChan() *Room {
	r := &Room{
		join:       make(chan *Client),
		leave:      make(chan *Client),
		subscribe:  make(chan *Client),
		clients:    make(map[*Client]bool),
		nsqReaders: make(map[string]*NsqReader),
	}
	go run(r)
	return r
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize,
	WriteBufferSize: socketBufferSize}

// RoomChat take a room, return a HandlerFunc,
// responsible for send & receive websocket data for all channels
func RoomChat(r *Room) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		vars := mux.Vars(req)

		socket, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			log.Print("ServeHTTP:", err)
			return
		}

		// Get user out of session
		session, _ := config.Store.Get(req, "session")
		val := session.Values["user"]
		var user = &User{}
		var ok bool
		if user, ok = val.(*User); !ok {
			log.Print("Invalid session")
			return
		}

		// Create new Client for this connection & join it to the Room
		client := &Client{
			socket:     socket,
			send:       make(chan *Message, messageBufferSize),
			room:       r,
			user:       user,
			channel:    vars["id"],
			subscribed: make(map[string]*NsqReader),
		}
		r.join <- client

		// We also subscribe to the NSQ Channel for the Room in order to receive
		// messages from NSQ later
		r.subscribe <- client

		defer func() {
			r.leave <- client
		}()
		go client.write()
		client.read()
	}
}
