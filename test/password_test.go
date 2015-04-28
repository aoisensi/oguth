package test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"golang.org/x/oauth2"

	"github.com/aoisensi/oguth"
)

const (
	AUTH_URL  = "http://localhost:1919/authorize"
	TOKEN_URL = "http://localhost:1919/token"
	TEST_URL  = "http://localhost:1919/nickname"

	USERNAME = "aoisensi"
	PASSWORD = "pAss30rD"
	NICKNAME = "I♥GO"
)

func TestPasswordAuth(t *testing.T) {
	cfg := oguth.NewConfig()
	cfg.AuthorizeEndpoint = AUTH_URL
	cfg.AccessTokenEndpoint = TOKEN_URL
	oauth := oguth.NewOAuth(cfg, NewPWOwner())
	server := http.NewServeMux()
	server.HandleFunc("/authorize", oauth.AuthorizeRequestHandler)
	server.HandleFunc("/token", oauth.AccessTokenRequestHandler)
	server.HandleFunc("/nickname", func(w http.ResponseWriter, r *http.Request) {
		cli, err := oauth.VerifyAccess(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		client, ok := cli.(*PWClient)
		if !ok {
			http.Error(w, "cast failed", http.StatusInternalServerError)
		}
		w.Write([]byte(client.Nickname))
	})
	go http.ListenAndServe("localhost:1919", server)

	client := &oauth2.Config{
		Endpoint: oauth2.Endpoint{
			AuthURL:  AUTH_URL,
			TokenURL: TOKEN_URL,
		},
	}
	token, err := client.PasswordCredentialsToken(oauth2.NoContext, USERNAME, PASSWORD)
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest("GET", TEST_URL, nil)
	token.SetAuthHeader(req)
	resp, errr := new(http.Client).Do(req)
	if errr != nil {
		t.Fatal(errr)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	result := string(body)
	if NICKNAME != result {
		t.Fatal(result)
	}
}

type PWOwner struct {
	oguth.ResourceOwner
}

type PWClient struct {
	Nickname string
}

func NewPWOwner() *PWOwner {
	return &PWOwner{}
}

func (o *PWOwner) GetClientWithPasswordGrant(username, password string) oguth.Client {
	if username != USERNAME || password != PASSWORD {
		return nil
	}
	return &PWClient{Nickname: NICKNAME}
}
