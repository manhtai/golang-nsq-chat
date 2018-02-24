package models

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/manhtai/golang-nsq-chat/pkg/config"
)

// Client represents a user connect to a room, one user may have many devices to chat,
// so it should not be the same as user
type Client struct {
	channel string
	// socket is the web socket for this client.
	socket *websocket.Conn
	// send is a channel on which messages are sent.
	send chan *Message
	// room is the room this client is chatting in.
	room *Room
	// user uses this client to chat
	user *User
}

func (c *Client) read() {
	defer c.socket.Close()
	for {
		var msg *Message
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			log.Print("Error when reading message from Websocket: ", err)
			return
		}

		msg.Name = c.user.Name
		msg.Channel = c.channel
		msg.User = c.user.ID
		msg.Timestamp = time.Now()

		SendMessageToTopic(config.TopicName, msg)
	}
}

func (c *Client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		// Drop messages if it's not the same channel
		if c.channel != msg.Channel {
			continue
		}
		err := c.socket.WriteJSON(msg)
		if err != nil {
			return
		}
	}
}
