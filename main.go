package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/manhtai/golang-nsq-chat/controllers"
	"github.com/manhtai/golang-nsq-chat/models"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("PORT must be set!")
	}
	port = ":" + port

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

	log.Printf("Starting web server %s:%d on port%s", hostName, os.Getpid(), port)
	log.Fatal(http.ListenAndServe(port, router))
}
