package config

import (
	"log"
	"os"
)

// AddrNsqlookupd holds address for nsqlookupd
// TODO: Support multiple lookupd addresses
var AddrNsqlookupd string

func init() {
	AddrNsqlookupd = os.Getenv("LOOKUPD_HTTP_ADDRESS")
	if AddrNsqlookupd == "" {
		log.Fatal("LOOKUPD_HTTP_ADDRESS must be set.")
	}
}
