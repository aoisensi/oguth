package test

import (
	"fmt"
	"net/url"
	"testing"
)

func TestCodeAuth(t *testing.T) {
	go startTestServer()

	cfg := OConfig

	cfg.ClientID = CLIENT_ID
	cfg.ClientSecret = CLIENT_SECRET

	aurl := cfg.AuthCodeURL("state")
	lq, _ := url.ParseQuery(fastHttpGet(nil, aurl, t))
	cq := fastHttpPost(nil, LOGIN_URL, lq, t)
	if cq != NICKNAME {
		fmt.Println("failed")
	}
}
