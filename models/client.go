package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/manhtai/golang-nsq-chat/config"
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

	// NsqReader that the client subscribed to
	subscribed map[string]*NsqReader
}

func (c *Client) read() {
	defer c.socket.Close()
	for {
		var msg *Message
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			log.Print(err)
			return
		}

		msg.Name = c.user.Name
		msg.Channel = c.channel
		msg.User = c.user.ID
		msg.Timestamp = time.Now()

		// Send to NSQ
		httpclient := &http.Client{}
		url := fmt.Sprintf(config.AddrNsqlookupd+"/pub?topic=%s", topicName)

		msgJSON, _ := json.Marshal(msg)
		nsqReq, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(msgJSON)))
		nsqResp, err := httpclient.Do(nsqReq)

		if err != nil {
			log.Println("NSQ publish error: " + err.Error())
		}
		defer nsqResp.Body.Close()
	}
}

func (c *Client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			return
		}
	}
}
