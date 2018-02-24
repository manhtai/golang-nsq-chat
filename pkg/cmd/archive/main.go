package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
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
	cfg.UserAgent = fmt.Sprintf("Archive client go-nsq/%s", nsq.VERSION)

	nsqConsumer, err := nsq.NewConsumer(config.TopicName, channelName, cfg)

	if err != nil {
		log.Println("nsq.NewNsqReader error: ", err)
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
	err = config.Mgo.DB("").C("messages").Insert(message)
	return err
}

func main() {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go subscribeToNsq(config.ArchiveChannelName)

	// Wait here for SigInt or SigTerm
	<-sigs
}
