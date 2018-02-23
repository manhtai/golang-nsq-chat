package controllers

import (
	"context"
	"time"

	"github.com/gorilla/mux"
	"github.com/manhtai/golang-nsq-chat/pkg/config"
	"github.com/manhtai/golang-nsq-chat/pkg/models"
	"gopkg.in/mgo.v2/bson"

	"fmt"
	"net/http"

	"github.com/markbates/goth/gothic"
)

// Index is the index page
func Index(w http.ResponseWriter, r *http.Request) {
	config.Templ.ExecuteTemplate(w, "index.html", nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	config.Templ.ExecuteTemplate(w, "login.html", nil)
}

func loginHandle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	action := vars["action"]
	provider := vars["provider"]

	switch action {
	case "login":
		q := r.URL.Query()
		q.Set("provider", provider)

		r.URL.RawQuery = q.Encode()

		// FIXME: Find a way to change callback url in place
		config.CreateProvider("https://" + r.Host + "/auth/callback/" + provider)
		gothic.BeginAuthHandler(w, r)

	case "callback":
		gothUser, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		var user = &models.User{}
		var ID = gothUser.Provider + gothUser.UserID

		config.Mgo.DB("").C("users").Find(
			bson.M{"ID": ID},
		).One(&user)

		if user.ID == "" {
			// FIXME: Can you do better?
			user.ID = ID
			user.UserID = gothUser.UserID
			user.Provider = gothUser.Provider
			user.Active = true
			user.CreatedAt = time.Now()
		}

		user.Email = gothUser.Email
		user.Name = gothUser.Name
		user.FirstName = gothUser.FirstName
		user.LastName = gothUser.LastName
		user.NickName = gothUser.NickName
		user.Description = gothUser.Description
		user.AvatarURL = gothUser.AvatarURL
		user.Location = gothUser.Location
		user.AccessToken = gothUser.AccessToken
		user.AccessTokenSecret = gothUser.AccessTokenSecret
		user.RefreshToken = gothUser.RefreshToken
		user.ExpiresAt = gothUser.ExpiresAt

		// Update or insert new user
		config.Mgo.DB("").C("users").UpsertId(ID, user)

		session, _ := config.Store.Get(r, "session")
		session.Values["user"] = user
		session.Save(r, w)

		w.Header().Set("Location", "/channel")
		w.WriteHeader(http.StatusTemporaryRedirect)

	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := config.Store.Get(r, "session")
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// MustAuth is a login required decorator for HandlerFunc
func MustAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := config.Store.Get(r, "session")
		if err != nil {
			fmt.Println(err)
		}

		val := session.Values["user"]
		user, ok := val.(*models.User)

		if !ok {
			// not authenticated
			w.Header().Set("Location", "/auth/login")
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, models.UserKey(0), user)

		// success - call the original handler
		handlerFunc(w, r.WithContext(ctx))
	}
}

// MustNotAuth is a anonymous required decorator for HandlerFunc
func MustNotAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		session, err := config.Store.Get(r, "session")
		if err != nil {
			fmt.Println(err)
		}

		val := session.Values["user"]
		if _, ok := val.(*models.User); ok {
			// already authenticated
			w.Header().Set("Location", "/")
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}

		// success - call the original handler
		handlerFunc(w, r)
	}
}

// Logout uses to log User out
var Logout = MustAuth(logout)

// LoginHandle redirect user to providers' login page & receive callback from them
var LoginHandle = MustNotAuth(loginHandle)

// Login servers our login page
var Login = MustNotAuth(login)
