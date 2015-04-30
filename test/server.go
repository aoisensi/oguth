package test

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/aoisensi/oguth"
	"golang.org/x/oauth2"
)

const (
	LOGIN_URL    = "http://localhost:1919/login"
	REDIRECT_URL = "http://localhost:1919/redirect"
	AUTH_URL     = "http://localhost:1919/authorize"
	TOKEN_URL    = "http://localhost:1919/token"
	TEST_URL     = "http://localhost:1919/nickname"
	CLIENT_URL   = "http://localhost:1919/client"

	USERNAME = "aoisensi"
	PASSWORD = "pAss30rD"
	NICKNAME = "I♥GO"

	CLIENT_ID     = "cLiEnTiD"
	CLIENT_SECRET = "cLiEnTsEcReT"
)

var (
	User = &Client{Nickname: NICKNAME}
)

var (
	OConfig = oauth2.Config{
		Endpoint: oauth2.Endpoint{
			AuthURL:  AUTH_URL,
			TokenURL: TOKEN_URL,
		},
	}
)

var Token *oauth2.Token

func startTestServer() {
	cfg := oguth.NewConfig()
	owner := NewOwner()
	cfg.Owner = owner
	cfg.AuthorizeEndpoint = AUTH_URL
	cfg.AccessTokenEndpoint = TOKEN_URL
	oauth := oguth.NewOAuth(cfg)
	server := http.NewServeMux()
	server.HandleFunc("/authorize", oauth.AuthorizeRequestHandler)
	server.HandleFunc("/token", oauth.AccessTokenRequestHandler)
	server.HandleFunc("/nickname", func(w http.ResponseWriter, r *http.Request) {
		cli, err := oauth.VerifyAccess(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		client, ok := cli.(*Client)
		if !ok {
			http.Error(w, "cast failed", http.StatusInternalServerError)
		}
		w.Write([]byte(client.Nickname))
	})
	server.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			q := url.Values{
				"username": {USERNAME},
				"password": {PASSWORD},
			}
			query := r.URL.Query().Encode() + "&" + q.Encode()
			w.Write([]byte(query))
			return
		} else if r.Method == "POST" {
			r.ParseForm()
			q := r.Form
			username := q.Get("username")
			password := q.Get("password")
			cli := owner.GetClientWithPasswordGrant(username, password)
			if cli == nil {
				http.Error(w, "login failed", http.StatusBadRequest)
				return
			}
			oauth.ConnectClientToCode(q.Get("code"), cli)
			query := url.Values{
				"response_type": {"code"},
				"client_id":     {q.Get("client_id")},
				"scope":         {q.Get("scope")},
				"username":      {username},
			}
			state := q.Get("state")
			if state != "" {
				query.Add("state", state)
			}
			oauth.RedirectAuthorize(w, r, query)
			return
		}
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	})
	server.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
		code := r.FormValue("code")
		token, err := OConfig.Exchange(oauth2.NoContext, code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		Token = token
		http.Redirect(w, r, CLIENT_URL, http.StatusFound)
	})
	server.HandleFunc("/client", func(w http.ResponseWriter, r *http.Request) {
		if Token == nil {
			http.Error(w, "error", http.StatusNotFound)
		}
		cli := OConfig.Client(oauth2.NoContext, Token)
		resp, err := cli.Get(TEST_URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		w.Write(body)
	})

	http.ListenAndServe("localhost:1919", server)
}

type Owner struct {
	oguth.ResourceOwner
}

type Client struct {
	Nickname string
}

func NewOwner() *Owner {
	return &Owner{}
}

func (o *Owner) GetClientWithPasswordGrant(username, password string) oguth.Client {
	if username != USERNAME || password != PASSWORD {
		return nil
	}
	return User
}

func (o *Owner) ExistClientId(id string) bool {
	return id == CLIENT_ID
}

func (o *Owner) AuthCodeMissing(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, LOGIN_URL+"?"+r.URL.Query().Encode(), http.StatusFound)
}

func (o *Owner) AuthCodeDecision(r *http.Request, clientId string) oguth.Client {
	if USERNAME == r.URL.Query().Get("username") { //bad method
		return User
	}
	return nil
}

func (o *Owner) GetRedirectUri(clientId string) string {
	if clientId == CLIENT_ID {
		return REDIRECT_URL
	}
	return ""
}

func fastHttpGet(client *http.Client, url string, t *testing.T) string {
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func fastHttpPost(client *http.Client, url string, v url.Values, t *testing.T) string {
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.PostForm(url, v)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}
