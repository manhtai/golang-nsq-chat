package models

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	nsq "github.com/bitly/go-nsq"
	"github.com/manhtai/golang-nsq-chat/config"
)

// We use only 1 topic for our chat server
const topicName = "Chat"

// ClientsMap holds all clients that subscribed to a NSQ channel
type ClientsMap map[*Client]bool

// ChannelsMap holds all NSQ channels that a client subscribed to
type ChannelsMap map[string]bool

// NsqReader represents a NSQ channel below topic Chat
type NsqReader struct {
	channelName     string
	consumer        *nsq.Consumer
	muMapsLock      sync.RWMutex
	channel2clients map[string]ClientsMap
	client2channels map[*Client]ChannelsMap
}

// NewNsqReader create new NsqReader from a channel name
func NewNsqReader(channelName string) (*NsqReader, error) {

	cfg := nsq.NewConfig()
	cfg.UserAgent = fmt.Sprintf("go-nsq/%s", nsq.VERSION)

	nsqConsumer, err := nsq.NewConsumer(topicName, channelName, cfg)

	if err != nil {
		log.Println("nsq.NewNsqReader error: " + err.Error())
		return nil, err
	}

	httpclient := &http.Client{}
	url := fmt.Sprintf(config.AddrNsqlookupd+"/create_topic?topic=%s", topicName)
	nsqReq, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	nsqResp, err := httpclient.Do(nsqReq)
	if err != nil {
		log.Println("NSQ create topic error: " + err.Error())
		return nil, err
	}
	defer nsqResp.Body.Close()

	nsqReader := &NsqReader{
		channelName:     channelName,
		channel2clients: make(map[string]ClientsMap),
		client2channels: make(map[*Client]ChannelsMap),
	}

	nsqConsumer.AddHandler(nsqReader)

	nsqErr := nsqConsumer.ConnectToNSQLookupd(config.AddrNsqlookupd)
	if nsqErr != nil {
		log.Println("NSQ connection error: " + nsqErr.Error())
		return nil, err
	}
	nsqReader.consumer = nsqConsumer

	return nsqReader, nil
}

// AddClient adds a client to a channel
func (nr *NsqReader) AddClient(c *Client, channel string) {
	log.Println("Add client for channel " + channel)

	nr.muMapsLock.Lock()
	if nr.channel2clients[channel] == nil {
		nr.channel2clients[channel] = make(ClientsMap)
	}
	if nr.client2channels[c] == nil {
		nr.client2channels[c] = make(ChannelsMap)
	}
	nr.channel2clients[channel][c] = true
	nr.client2channels[c][channel] = true
	nr.muMapsLock.Unlock()
}

// RemoveClient removes a client out of a channel
func (nr *NsqReader) RemoveClient(c *Client) {
	nr.muMapsLock.Lock()
	for channel := range nr.client2channels[c] {
		delete(nr.channel2clients[channel], c)
	}
	delete(nr.client2channels, c)
	nr.muMapsLock.Unlock()
}

// HandleMessage pushs messages from NSQ to Client, is used by AddHandler() function
func (nr *NsqReader) HandleMessage(msg *nsq.Message) error {
	message := Message{}
	err := json.Unmarshal(msg.Body, &message)
	if err != nil {
		log.Println("NSQ HandleMessage ERROR: invalid JSON subscribe data")
		return err
	}
	nr.muMapsLock.RLock()
	for c := range nr.channel2clients[message.Channel] {
		c.send <- &message
	}
	nr.muMapsLock.RUnlock()
	return nil
}
