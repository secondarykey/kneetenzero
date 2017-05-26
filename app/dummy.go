package dummy

import (
	_ "golang.org/x/oauth2"

	_ "github.com/gorilla/sessions"

	_ "github.com/knightso/base/gae/ds"

	_ "google.golang.org/appengine"
	_ "google.golang.org/appengine/datastore"
	_ "google.golang.org/appengine/user"
)
