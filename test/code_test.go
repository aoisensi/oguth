package test

import (
	"net/url"
	"testing"

	"github.com/k0kubun/pp"
	"golang.org/x/oauth2"
)

func TestCodeAuth(t *testing.T) {
	go startTestServer()
	conf := oauth2.Config{
		Endpoint:     Endpoint,
		ClientID:     CLIENT_ID,
		ClientSecret: CLIENT_SECRET,
	}
	aurl := conf.AuthCodeURL("state")
	lq, _ := url.ParseQuery(fastHttpGet(nil, aurl, t))
	pp.Println(aurl)
	curl := fastHttpPost(nil, LOGIN_URL, lq, t)
	pp.Println(curl)
}
