package datastore

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"

	//"github.com/dghubble/go-twitter/twitter"
	"github.com/knightso/base/gae/ds"
)

const KIND_AGENT = "Agent"

type Agent struct {
	TwitterToken  string
	TwitterSecret string
	DropboxToken  string
	ds.Meta
}

func CreateAgentKey(r *http.Request) *datastore.Key {
	c := appengine.NewContext(r)
	u := user.Current(c)
	return datastore.NewKey(c, KIND_AGENT, u.ID, 0, nil)
}

func GetAgent(r *http.Request) (*Agent, error) {

	var pkgAgent Agent
	c := appengine.NewContext(r)

	key := CreateAgentKey(r)
	err := ds.Get(c, key, &pkgAgent)
	if err != nil && err.Error() == datastore.ErrNoSuchEntity.Error() {
		pkgAgent.SetKey(key)
		return &pkgAgent, nil
	}

	return &pkgAgent, err
}

func UpdateDropboxAgent(r *http.Request, token string) error {

	agent, err := GetAgent(r)
	if err != nil {
		return err
	}

	agent.DropboxToken = token

	c := appengine.NewContext(r)
	err = ds.Put(c, agent)
	if err != nil {
		return err
	}
	return nil
}
