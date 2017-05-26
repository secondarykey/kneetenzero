package kneetenzero

import (
	"html/template"
	"net/http"

	"kneetenzero/datastore"
	"kneetenzero/dropbox"
	"kneetenzero/twitter"

	"google.golang.org/appengine"
	ds "google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
)

func init() {
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/dashboard/", dashboardHandler)
	http.HandleFunc("/dashboard/register", registerHandler)

	http.HandleFunc("/dashboard/twitter", twitter.RedirectHandler)
	http.HandleFunc("/dashboard/twitter/callback", twitter.CallbackHandler)

	http.HandleFunc("/dashboard/dropbox", dropbox.RedirectHandler)
	http.HandleFunc("/dashboard/dropbox/callback", dropbox.CallbackHandler)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	url, _ := user.LogoutURL(ctx, "/")
	http.Redirect(w, r, url, http.StatusFound)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	err := datastore.PutApplication(r)
	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, "/dashboard/", http.StatusFound)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {

	_, err := datastore.GetApplication(r)

	if err != nil && err.Error() == ds.ErrNoSuchEntity.Error() {
		err = register(w, r)
	} else {

		tw := false
		db := false

		agent, err := datastore.GetAgent(r)

		if err != nil {
			panic(err)
		}

		//Agentを検索
		if agent.TwitterToken != "" && agent.TwitterSecret != "" {
			tw = true
		}
		if agent.DropboxToken != "" {
			db = true
		}

		if tw && db {
			err = word(w, r)
		} else {
			err = token(w, r, tw, db)
		}
	}

	if err != nil {
		panic(err)
	}
}

func register(w http.ResponseWriter, r *http.Request) error {
	return layout(w, r, nil, "templates/api.tmpl")
}

func token(w http.ResponseWriter, r *http.Request, tw, db bool) error {
	return layout(w, r, nil, "templates/token.tmpl")
}

func word(w http.ResponseWriter, r *http.Request) error {
	return layout(w, r, nil, "templates/word.tmpl")
}

func layout(w http.ResponseWriter, r *http.Request, obj interface{}, files ...string) error {

	args := make([]string, len(files)+1)

	copy(args, files)

	args[len(files)] = "templates/layout.tmpl"
	tpl, err := template.New("layout").ParseFiles(args...)
	if err != nil {
		return err
	}
	// テンプレートを出力
	err = tpl.Execute(w, obj)
	if err != nil {
		return err
	}
	return nil
}
