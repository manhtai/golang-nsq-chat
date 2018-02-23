package models

import (
	"encoding/gob"
	"time"
)

func init() {
	gob.Register(&User{})
}

// UserKey is key to save user to context
type UserKey int

// User hold information about an user
type User struct {
	ID                string    `json:"id" bson:"_id"`
	Provider          string    `json:"provider" bson:"provider"`
	Email             string    `json:"email" bson:"email"`
	Name              string    `json:"name" bson:"name"`
	FirstName         string    `json:"first_name" bson:"first_name"`
	LastName          string    `json:"last_name" bson:"last_name"`
	NickName          string    `json:"nick_name" bson:"nick_name"`
	Description       string    `json:"description" bson:"description"`
	UserID            string    `json:"user_id" bson:"user_id"`
	AvatarURL         string    `json:"avatar_url" bson:"avatar_url"`
	Location          string    `json:"location" bson:"location"`
	AccessToken       string    `json:"access_token" bson:"access_token"`
	AccessTokenSecret string    `json:"access_token_secret" bson:"access_token_secret"`
	RefreshToken      string    `json:"refresh_token" bson:"refresh_token"`
	ExpiresAt         time.Time `json:"expires_at" bson:"expires_at"`
	CreatedAt         time.Time `json:"created_at" bson:"created_at"`
	Active            bool      `json:"active" bson:"active"`
}
