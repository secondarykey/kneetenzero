package twitter

import (
	"kneetenzero/datastore"
	"kneetenzero/session"

	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"

	"github.com/dghubble/oauth1"
	"github.com/dghubble/oauth1/twitter"

	"github.com/knightso/base/gae/ds"
)

func getConfig(r *http.Request) (*oauth1.Config, error) {
	app, err := datastore.GetApplication(r)
	if err != nil {
		return nil, err
	}

	//Application から取得
	config := oauth1.Config{
		ConsumerKey:    app.TwitterKey,
		ConsumerSecret: app.TwitterSecret,
		CallbackURL:    "https://kneetenzero.appspot.com/dashboard/twitter/callback",
		Endpoint:       twitter.AuthorizeEndpoint,
	}
	return &config, nil
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {

	token, err := session.SetTwitterRequestToken(w, r)
	if err != nil {
		panic(err)
	}

	config, err := getConfig(r)
	if err != nil {
		panic(err)
	}

	authorizationURL, err := config.AuthorizationURL(token)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, authorizationURL.String(), http.StatusFound)
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {

	verifier := r.FormValue("oauth_verifier")

	requestToken, requestSecret, err := session.GetTwitterRequestToken(r)
	if err != nil {
		panic(err)
	}

	config, err := getConfig(r)
	if err != nil {
		panic(err)
	}

	c := appengine.NewContext(r)
	cli := urlfetch.Client(c)
	accessToken, accessSecret, err := config.AccessToken(cli, requestToken, requestSecret, verifier)
	if err != nil {
		panic(err)
	}

	agent, err := datastore.GetAgent(r)
	if err != nil {
		panic(err)
	}

	agent.TwitterToken = accessToken
	agent.TwitterSecret = accessSecret

	err = ds.Put(c, agent)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/dashboard/", http.StatusFound)
}
