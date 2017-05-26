package dropbox

import (
	"kneetenzero/datastore"
	"net/http"

	"golang.org/x/oauth2"
)

func GetConfig(r *http.Request) (*oauth2.Config, error) {

	app, err := datastore.GetApplication(r)
	if err != nil {
		return nil, err
	}

	config := &oauth2.Config{
		ClientID:     app.DropboxKey,
		ClientSecret: app.DropboxSecret,
		Scopes:       nil,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.dropbox.com/oauth2/authorize",
			TokenURL: "https://api.dropboxapi.com/oauth2/token",
		},
	}
	return config, nil
}
