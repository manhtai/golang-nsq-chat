package models

import (
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
	defer func() {
		c.room.leave <- c
		c.socket.Close()
	}()

	c.socket.SetReadLimit(config.MaxMessageSize)
	c.socket.SetReadDeadline(time.Now().Add(config.PongWait))
	c.socket.SetPongHandler(func(string) error {
		c.socket.SetReadDeadline(time.Now().Add(config.PongWait))
		return nil
	})

	for {
		var msg *Message
		if err := c.socket.ReadJSON(&msg); err != nil {
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
	ticker := time.NewTicker(config.PingPeriod)

	defer func() {
		c.socket.Close()
		ticker.Stop()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Drop messages if it's not the same channel
			if c.channel != msg.Channel {
				continue
			}
			if err := c.socket.WriteJSON(msg); err != nil {
				return
			}

		case <-ticker.C:
			c.socket.SetWriteDeadline(time.Now().Add(config.WriteWait))
			if err := c.socket.WriteMessage(websocket.PingMessage,
				[]byte{}); err != nil {
				return
			}
		}
	}

}
