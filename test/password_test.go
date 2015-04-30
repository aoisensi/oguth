package test

import (
	"io/ioutil"
	"testing"

	"golang.org/x/oauth2"
)

func TestPasswordAuth(t *testing.T) {
	go startTestServer()

	cfg := OConfig

	token, err := cfg.PasswordCredentialsToken(oauth2.NoContext, USERNAME, PASSWORD)
	if err != nil {
		t.Fatal(err)
	}
	client := cfg.Client(oauth2.NoContext, token)
	resp, err := client.Get(TEST_URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	result := string(body)
	if NICKNAME != result {
		t.Fatal(result)
	}
}
