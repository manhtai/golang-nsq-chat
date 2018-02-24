package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	nsq "github.com/bitly/go-nsq"
	"github.com/manhtai/golang-nsq-chat/pkg/config"
)

// NsqReader represents a NSQ channel below topic Chat
type NsqReader struct {
	channelName string
	consumer    *nsq.Consumer
	rooms       map[*Room]bool
}

// newNsqReader create new NsqReader from a channel name
func newNsqReader(r *Room, channelName string) error {

	cfg := nsq.NewConfig()
	cfg.Set("LookupdPollInterval", config.LookupdPollInterval*time.Second)
	cfg.Set("MaxInFlight", config.MaxInFlight)
	cfg.UserAgent = fmt.Sprintf("Chat client go-nsq/%s", nsq.VERSION)

	nsqConsumer, err := nsq.NewConsumer(config.TopicName, channelName, cfg)

	if err != nil {
		log.Println("Create newNsqReader error: ", err)
		return err
	}

	nsqReader := &NsqReader{
		channelName: channelName,
		rooms:       map[*Room]bool{r: true},
	}
	r.nsqReaders[channelName] = nsqReader

	nsqConsumer.AddHandler(nsqReader)

	nsqErr := nsqConsumer.ConnectToNSQLookupd(config.AddrNsqlookupd)
	if nsqErr != nil {
		log.Println("NSQ connection error: ", nsqErr)
		return err
	}
	nsqReader.consumer = nsqConsumer
	log.Printf("Subscribe to NSQ success to channel %s", channelName)

	return nil
}

// HandleMessage pushes messages from NSQ to Client, is used by AddHandler() function
func (nr *NsqReader) HandleMessage(msg *nsq.Message) error {
	message := Message{}
	err := json.Unmarshal(msg.Body, &message)
	if err != nil {
		log.Println("NSQ HandleMessage ERROR: invalid JSON subscribe data")
		return err
	}
	for r := range nr.rooms {
		r.forward <- &message
	}
	return nil
}

// getChannelName return hostname of our chat server
func getChannelName() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "undefined"
	}
	hostname = "websocket-server-" + hostname

	maxLength := len(hostname)
	if maxLength > 63 {
		maxLength = 63
	}

	return hostname[0:maxLength]
}

// subscribeToNsq subscribes Room to a NSQ channel
func subscribeToNsq(r *Room) error {
	nsqChannelName := getChannelName()
	_, ok := r.nsqReaders[nsqChannelName]

	if !ok {
		err := newNsqReader(r, nsqChannelName)
		if err != nil {
			log.Printf("Failed to subscribe to channel: '%s'",
				nsqChannelName)
			return err
		}
	}
	return nil
}

// SendMessageToTopic send Message to NSQ topic
func SendMessageToTopic(topicName string, message *Message) error {
	httpclient := &http.Client{}
	url := fmt.Sprintf("http://"+config.AddrNsqd+"/pub?topic=%s", topicName)

	msgJSON, _ := json.Marshal(message)
	nsqReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(
		[]byte(string(msgJSON))))

	nsqResp, err := httpclient.Do(nsqReq)

	if err != nil {
		log.Print("NSQ publish error: ", err)
		return err
	}
	defer nsqResp.Body.Close()

	if nsqResp.StatusCode != 200 {
		log.Print("Fail to publish to NSQ: ", nsqResp.Status)
		return errors.New("Send message to NSQ fail: " + nsqResp.Status)
	}

	return nil
}
