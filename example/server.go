package main

import (
	"net/http"

	"github.com/aoisensi/oguth"
)

func server() {
	s := http.NewServeMux()
	config := oguth.NewConfig()
	config.TokenEndpoint = "http://localhost:19190"
	oauth := oguth.NewOAuth(config, nil)
	s.HandleFunc("/authorize", oauth.AuthorizeRequestHandler)
	http.ListenAndServe("localhost:19190", s)
}
