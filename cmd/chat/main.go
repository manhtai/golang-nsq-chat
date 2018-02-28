package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/manhtai/golang-nsq-chat/pkg/controllers"
	"github.com/manhtai/golang-nsq-chat/pkg/models"
)

var (
	addr     = flag.String("addr", ":3000", "http service address")
	certFile = flag.String("cert-file", "cert.pem", "cert file")
	keyFile  = flag.String("key-file", "key.pem", "key file")
)

func main() {
	flag.Parse()

	r := models.NewRoomChan()

	router := mux.NewRouter()

	// Homepage
	router.HandleFunc("/", controllers.Index)

	// Auth handlers
	router.HandleFunc("/auth/login", controllers.Login)
	router.HandleFunc("/auth/{action:(?:login|callback)}/{provider}",
		controllers.LoginHandle)
	router.HandleFunc("/auth/logout", controllers.Logout)

	// Chat handlers
	router.HandleFunc("/channel", controllers.ChannelList)
	router.HandleFunc("/channel/new", controllers.ChannelNew)
	router.HandleFunc("/channel/{id}/chat", models.RoomChat(r))
	router.HandleFunc("/channel/{id}/view", controllers.ChannelView)
	router.HandleFunc("/channel/{id}/history", controllers.ChannelHistory)

	// User handlers
	// router.GET("/user/", UserList)
	// router.GET("/user/:id", UserDetail)

	// The rest, just not found
	router.HandleFunc("/*", http.NotFound)

	hostName, _ := os.Hostname()

	log.Printf("Starting web server %s:%d on %s", hostName, os.Getpid(), *addr)

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		// Only use curves which have assembly implementations
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
	}

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		TLSConfig:    tlsConfig,
		Handler:      router,
		Addr:         *addr,
	}

	log.Fatal(srv.ListenAndServeTLS(*certFile, *keyFile))
}
