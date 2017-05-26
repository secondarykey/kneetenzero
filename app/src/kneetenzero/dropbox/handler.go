package dropbox

import (
	"kneetenzero/datastore"
	"kneetenzero/session"

	"net/http"

	"google.golang.org/appengine"

	"golang.org/x/oauth2"
)

func RedirectHandler(w http.ResponseWriter, r *http.Request) {

	conf, err := GetConfig(r)
	if err != nil {
		panic(err)
	}

	err = session.SetDropboxConfig(w, r, conf)
	if err != nil {
		panic(err)
	}

	redirect := oauth2.SetAuthURLParam("redirect_uri", "https://kneetenzero.appspot.com/dashboard/dropbox/callback")
	types := oauth2.SetAuthURLParam("response_type", "code")

	http.Redirect(w, r, conf.AuthCodeURL("", redirect, types), http.StatusFound)
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	code := r.FormValue("code")

	config, err := session.GetDropboxConfig(w, r)
	if err != nil {
		panic(err)
	}

	config.RedirectURL = "https://kneetenzero.appspot.com/dashboard/dropbox/callback"

	c := appengine.NewContext(r)
	token, err := config.Exchange(c, code)
	if err != nil {
		panic(err)
	}

	err = datastore.UpdateDropboxAgent(r, token.AccessToken)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/dashboard/", http.StatusFound)
}
