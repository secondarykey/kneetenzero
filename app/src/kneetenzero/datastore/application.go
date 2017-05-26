package datastore

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	//"github.com/dghubble/go-twitter/twitter"
	"github.com/knightso/base/gae/ds"
)

const KIND_APPLICATION = "Application"

type Application struct {
	TwitterKey    string
	TwitterSecret string
	DropboxKey    string
	DropboxSecret string
	ds.Meta
}

func createApplicationKey(r *http.Request) *datastore.Key {
	c := appengine.NewContext(r)
	return datastore.NewKey(c, KIND_APPLICATION, "Fixing", 0, nil)
}

func GetApplication(r *http.Request) (*Application, error) {
	var pkgApp Application
	c := appengine.NewContext(r)
	err := ds.Get(c, createApplicationKey(r), &pkgApp)
	return &pkgApp, err
}

func PutApplication(r *http.Request) error {

	c := appengine.NewContext(r)
	app := Application{
		TwitterKey:    r.FormValue("tw_id"),
		TwitterSecret: r.FormValue("tw_secret"),
		DropboxKey:    r.FormValue("db_id"),
		DropboxSecret: r.FormValue("db_secret"),
	}

	app.SetKey(createApplicationKey(r))

	err := ds.Put(c, &app)
	if err != nil {
		return err
	}

	return nil
}
