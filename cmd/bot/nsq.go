package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	nsq "github.com/bitly/go-nsq"
	"github.com/manhtai/golang-nsq-chat/pkg/config"
	"github.com/manhtai/golang-nsq-chat/pkg/models"
)

// subscribeToNsq subscribe to NSQ Archive channel
func subscribeToNsq(channelName string) error {

	cfg := nsq.NewConfig()
	cfg.Set("LookupdPollInterval", config.LookupdPollInterval*time.Second)
	cfg.Set("MaxInFlight", config.MaxInFlight)
	cfg.UserAgent = fmt.Sprintf("Bot client go-nsq/%s", nsq.VERSION)

	nsqConsumer, err := nsq.NewConsumer(config.TopicName, channelName, cfg)

	if err != nil {
		log.Println("subscribeToNsq error: ", err)
		return err
	}

	nsqConsumer.AddHandler(nsq.HandlerFunc(handleMessage))
	nsqErr := nsqConsumer.ConnectToNSQLookupd(config.AddrNsqlookupd)

	if nsqErr != nil {
		log.Println("NSQ connection error: ", nsqErr)
		return err
	}

	log.Printf("Subscribe to NSQ success to channel %s", channelName)
	return nil
}

// handleMessage pushes messages from NSQ to Mongodb
func handleMessage(msg *nsq.Message) error {
	message := models.Message{}
	err := json.Unmarshal(msg.Body, &message)
	if err != nil {
		log.Println("NSQ HandleMessage ERROR: invalid JSON subscribe data")
		return err
	}

	// Simple reply for now
	// TODO: Some NLP here?
	if strings.Index(message.Body, "@bot") == 0 {
		replyMessage := &models.Message{
			Name:      config.BotChannelName,
			Channel:   message.Channel,
			User:      config.BotChannelName,
			Timestamp: time.Now(),
			Body:      "Hi human, improve me!",
		}
		msgJSON, _ := json.Marshal(replyMessage)
		err = models.SendMessageToTopic(config.TopicName, []byte(string(msgJSON)))
	}
	return err
}
