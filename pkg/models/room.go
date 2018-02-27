package models

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/manhtai/golang-nsq-chat/pkg/config"
)

// Room represents a room to chat
type Room struct {
	// forward is a channel that holds incoming messages that should be forwarded
	// to the other clients in the same host
	forward chan *Message
	// join is a channel for clients wishing to join the room.
	join chan *Client
	// leave is a channel for clients wishing to leave the room.
	leave chan *Client
	// clients holds all current clients in this room, decided by join & leave channels
	clients map[*Client]bool
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
			if _, ok := r.clients[client]; ok {
				// leaving
				close(client.send)
				delete(r.clients, client)
			}

		case msg := <-r.forward:
			// forward message to all clients
			for client := range r.clients {
				client.send <- msg
			}
		}

	}
}

// NewRoomChan creates a new room for clients to join
func NewRoomChan() *Room {
	r := &Room{
		join:       make(chan *Client),
		leave:      make(chan *Client),
		forward:    make(chan *Message),
		clients:    make(map[*Client]bool),
		nsqReaders: make(map[string]*NsqReader),
	}
	go run(r)
	// We subscribe the Room to the NSQ Channel in order to
	// receive messages from NSQ later
	subscribeToNsq(r)
	return r
}

// RoomChat take a room, return a HandlerFunc,
// responsible for send & receive websocket data for all channels
func RoomChat(r *Room) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		vars := mux.Vars(req)
		upgrader := &websocket.Upgrader{ReadBufferSize: config.SocketBufferSize,
			WriteBufferSize: config.SocketBufferSize}

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
			socket:  socket,
			send:    make(chan *Message, config.MessageBufferSize),
			room:    r,
			user:    user,
			channel: vars["id"],
		}
		r.join <- client

		defer func() {
			r.leave <- client
		}()
		go client.write()
		client.read()
	}
}
