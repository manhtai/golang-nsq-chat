package models

import (
	"log"
	"time"

	"github.com/manhtai/golang-mongodb-chat/config"
)

// Message represents a single message which a client sent to a room
// (same meaning as a user send to a channel)
type Message struct {
	Name      string    `json:"name" bson:"name"`
	Body      string    `json:"body" bson:"body"`
	Channel   string    `json:"channel" bson:"channel"`
	User      string    `json:"user" bson:"user"`
	Timestamp time.Time `json:"timestamp,omitempty" bson:"timestamp"`
}

// SaveMessage uses to save Message to db
type SaveMessage struct {
	message *Message
}

func saveMessages(sm *chan SaveMessage) {
	for {
		sM, ok := <-*sm
		if !ok {
			log.Print("Error when receiving message to save")
			return
		}

		err := config.Mgo.DB("").C("messages").Insert(sM.message)
		if err != nil {
			log.Print(err)
		}
	}
}

// NewSaveMessageChan create a new SaveMessage channel
func NewSaveMessageChan() *chan SaveMessage {
	sm := make(chan SaveMessage, 256)
	go saveMessages(&sm)
	return &sm
}
