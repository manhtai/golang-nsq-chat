package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/manhtai/golang-nsq-chat/pkg/config"
)

func main() {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go subscribeToNsq(config.ArchiveChannelName)

	// Wait here for SigInt or SigTerm
	<-sigs
}
