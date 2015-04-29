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
	REDIRECT_URL = "http://localhost:1919/login"
	LOGIN_URL    = REDIRECT_URL
	AUTH_URL     = "http://localhost:1919/authorize"
	TOKEN_URL    = "http://localhost:1919/token"
	TEST_URL     = "http://localhost:1919/nickname"

	USERNAME = "aoisensi"
	PASSWORD = "pAss30rD"
	NICKNAME = "I♥GO"

	CLIENT_ID     = "cLiEnTiD"
	CLIENT_SECRET = "cLiEnTsEcReT"
)

var (
	Endpoint = oauth2.Endpoint{AuthURL: AUTH_URL, TokenURL: TOKEN_URL}
)

func startTestServer() {
	cfg := oguth.NewConfig()
	owner := NewOwner()
	cfg.Owner = owner
	cfg.RedirectEndpoint = REDIRECT_URL
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
			q.Del("username")
			q.Del("password")
			cli := owner.GetClientWithPasswordGrant(username, password)
			if cli == nil {
				http.Error(w, "login failed", http.StatusBadRequest)
				return
			}
			oauth.ConnectClientToCode(r.Form.Get("code"), cli)
			oauth.RedirectAuthorize(w, r, q)
			return
		}
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	})
	/*
		server = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			pp.Println(r.Method)
			pp.Println(r.URL.String())
			r.ParseForm()
			pp.Println(r.Form)
			pp.Println("==========================")
			server.ServeHTTP(w, r)
		})
	*/
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
	return &Client{Nickname: NICKNAME}
}

func (o *Owner) ExistClientId(id string) bool {
	return id == CLIENT_ID
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
