package oguth

import (
	"encoding/json"
	"errors"
	"net/http"
)

type OAuth struct {
	owner  ResourceOwner
	config Config
}

func NewOAuth(config Config, owner ResourceOwner) *OAuth {
	return &OAuth{config: config, owner: owner}
}

func (a *OAuth) Close() {
	a.owner.Close()
}

func (a *OAuth) AuthorizeRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowd.", http.StatusMethodNotAllowed)
		return
	}
	endpoint := a.config.AccessTokenEndpoint
	t := ResponseType(r.FormValue("response_type"))
	authorizer := a.config.AuthHandlers[t]
	if authorizer == nil {
		v := NewError(ErrorCodeUnsupportedResponseType).ToValues()
		http.Redirect(w, r, endpoint+"?"+v.Encode(), http.StatusFound)
		return
	}
	v := authorizer(a, r)
	http.Redirect(w, r, endpoint+"?"+v.Encode(), http.StatusFound)
}

func (a *OAuth) AccessTokenRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowd.", http.StatusMethodNotAllowed)
		return
	}
	t := GrantType(r.FormValue("grant_type"))
	access := a.config.AccessHandlers[t]
	if access == nil {
		e := NewError(ErrorCodeUnsupportedGrantType)
		e.Write(w)
		return
	}
	o, scode := access(a, r)
	body, _ := json.Marshal(o)

	h := w.Header()
	h.Set("Content-Type", "application/json")
	h.Add("Cache-Control", "no-store")
	h.Add("Pragma", "no-cache")
	w.WriteHeader(scode)
	w.Write(body)
}

func (a *OAuth) VerifyAccess(w http.ResponseWriter, r *http.Request) (Client, error) {
	token := httpHeaderAuth(r)
	at := a.config.Storage.GetAccessToken(token)
	client := at.GetClient()
	if client != nil {
		return client, nil
	}
	return nil, errors.New("error")
}
