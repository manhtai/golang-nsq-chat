package config

import (
	"math"
	"os"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/gplus"
)

// Store keep our sessions
var Store *sessions.FilesystemStore

func init() {
	Store = sessions.NewFilesystemStore(os.TempDir(), []byte(os.Getenv("SECRET_KEY")))
	Store.MaxLength(math.MaxInt64)
	gothic.Store = Store
}

var providers = make(map[string]*gplus.Provider)

// CreateProvider return provider with corresponding redirect_uri
// We only init goth with Google+ for now
func CreateProvider(redirectURI string) {

	if providers[redirectURI] == nil {
		providers[redirectURI] = gplus.New(
			os.Getenv("GPLUS_KEY"),
			os.Getenv("GPLUS_SECRET"),
			redirectURI,
		)

		goth.UseProviders(
			providers[redirectURI],
		)
	}
}
