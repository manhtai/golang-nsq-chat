package config

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	// AddrNsqlookupd holds address for nsqlookupd
	// TODO: Support multiple lookupd addresses
	AddrNsqlookupd string

	// AddrNsqd uses for publishing messages only, it should be your local nsqd
	AddrNsqd string
)

const (
	// TopicName is NSQ topic name for our chat server
	TopicName = "Chat"

	// MaxInFlight is largest number of messages allowed in flight
	MaxInFlight = 10

	// LookupdPollInterval is interval for polling NSQ for new messages
	LookupdPollInterval = 30

	// ArchiveChannelName is the name of Archive Channel
	ArchiveChannelName = "Archive"

	// BotChannelName is the name of Bot Channel
	BotChannelName = "Bot"
)

func init() {
	AddrNsqlookupd = os.Getenv("NSQLOOKUPD_HTTP_ADDRESS")
	AddrNsqd = os.Getenv("NSQD_HTTP_ADDRESS")

	if AddrNsqlookupd == "" || AddrNsqd == "" {
		log.Fatal("NSQLOOKUPD_HTTP_ADDRESS & NSQD_HTTP_ADDRESS must be set.")
	}

	httpclient := &http.Client{}
	url := fmt.Sprintf("http://"+AddrNsqd+"/topic/create?topic=%s", TopicName)
	nsqReq, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	nsqResp, err := httpclient.Do(nsqReq)
	if err != nil {
		log.Fatal("NSQ create topic error: " + err.Error())
	}

	if nsqResp.StatusCode != 200 {
		log.Fatal("Fail to create NSQ topic: ", nsqResp.Status)
	}

	defer nsqResp.Body.Close()
}
