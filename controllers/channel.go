package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/manhtai/golang-mongodb-chat/config"
	"github.com/manhtai/golang-mongodb-chat/models"
	"gopkg.in/mgo.v2/bson"
)

// ChannelList lists all the channel available to chat
func ChannelList(w http.ResponseWriter, r *http.Request) {
	var data []models.Channel
	config.Mgo.DB("").C("channels").Find(nil).All(&data)
	config.Templ.ExecuteTemplate(w, "channel-list.html", data)
}

func channelNew(w http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{}

	user, ok := r.Context().Value(models.UserKey(0)).(*models.User)

	if r.Method == http.MethodPost && ok {
		// Stub an user to be populated from the body
		channel := models.Channel{}

		// Populate the user data
		err := r.ParseForm()
		if err != nil {
			log.Print(err)
		}

		channel.CreatedBy = user.ID
		channel.ID = bson.NewObjectId()
		channel.Name = r.FormValue("name")

		if channel.Name == "" {
			channel.Name = "No name"
		}

		// Write the user to mongo
		config.Mgo.DB("").C("channels").Insert(channel)

		data["channel"] = channel
	}

	config.Templ.ExecuteTemplate(w, "channel-new.html", data)
}

func channelView(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Host": r.Host,
	}

	var channel models.Channel
	vars := mux.Vars(r)
	id := bson.ObjectIdHex(vars["id"])

	config.Mgo.DB("").C("channels").FindId(id).One(&channel)
	data["channel"] = channel

	config.Templ.ExecuteTemplate(w, "channel-view.html", data)
}

func channelHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	const limit = 10
	result := make([]models.Message, limit)

	err := config.Mgo.DB("").C("messages").Find(
		bson.M{"channel": vars["id"]},
	).Sort("-timestamp").Limit(limit).All(&result)

	if err != nil {
		log.Print(err)
	}

	rj, _ := json.Marshal(result)

	w.Header().Set("Content-Type", "application/json")
	w.Write(rj)
}

// ChannelNew is used to create new chat channel
var ChannelNew = MustAuth(channelNew)

// ChannelHistory hold chat history in a channel
var ChannelHistory = MustAuth(channelHistory)

// ChannelView is where we chat, it displays history along with
// current chat in the channel
var ChannelView = MustAuth(channelView)
