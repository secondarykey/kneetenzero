package session

import (
	"kneetenzero/datastore"

	"encoding/gob"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/dghubble/oauth1"
	"github.com/dghubble/oauth1/twitter"
	"github.com/gorilla/sessions"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

var store = sessions.NewCookieStore([]byte("Kneetenzero-sessions-secret"))

func init() {
	gob.Register(&oauth2.Config{})
}

func GetTwitterRequestToken(r *http.Request) (string, string, error) {

	session, err := store.Get(r, "request")
	if err != nil {
		return "", "", err
	}

	requestToken := session.Values["token"].(string)
	requestSecret := session.Values["secret"].(string)

	return requestToken, requestSecret, nil
}

func SetTwitterRequestToken(w http.ResponseWriter, r *http.Request) (string, error) {

	app, err := datastore.GetApplication(r)
	if err != nil {
		return "", err
	}

	//Application から取得
	config := oauth1.Config{
		ConsumerKey:    app.TwitterKey,
		ConsumerSecret: app.TwitterSecret,
		CallbackURL:    "https://kneetenzero.appspot.com/dashboard/twitter/callback",
		Endpoint:       twitter.AuthorizeEndpoint,
	}

	c := appengine.NewContext(r)

	requestToken, requestSecret, err := config.RequestToken(urlfetch.Client(c))
	if err != nil {
		return "", err
	}

	session, err := store.Get(r, "request")
	if err != nil {
		return "", err
	}

	session.Values["token"] = requestToken
	session.Values["secret"] = requestSecret
	return requestToken, session.Save(r, w)
}

func GetDropboxConfig(w http.ResponseWriter, r *http.Request) (*oauth2.Config, error) {
	session, err := store.Get(r, "dropbox")
	if err != nil {
		return nil, err
	}

	config := session.Values["config"].(*oauth2.Config)
	return config, nil
}

func SetDropboxConfig(w http.ResponseWriter, r *http.Request, config *oauth2.Config) error {

	session, err := store.Get(r, "dropbox")
	if err != nil {
		return err
	}

	session.Values["config"] = config
	return session.Save(r, w)
}
