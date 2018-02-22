package config

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// AddrNsqlookupd holds address for nsqlookupd
// TODO: Support multiple lookupd addresses
var AddrNsqlookupd string

// AddrNsqd uses for publishing messages only, it should be your local nsqd
var AddrNsqd string

// TopicName is NSQ topic name for our chat server
const TopicName = "Chat"

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
