package oguth

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
)

type OAuth struct {
	owner  ResourceOwner
	config *Config
}

func NewOAuth(config Config) *OAuth {
	a := &OAuth{
		config: &config,
		owner:  config.Owner,
	}
	a.config.init()
	return a
}

func (a *OAuth) Close() {
	a.owner.Close()
}

func (a *OAuth) AuthorizeRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowd.", http.StatusMethodNotAllowed)
		return
	}
	t := r.FormValue("response_type")
	if t == "" {
		v := ErrorResponseTypeMissing.ToValues()
		a.writeAuthRedirect(w, r, v)
		return
	}
	authorizer := a.config.AuthHandlers[ResponseType(t)]
	if authorizer == nil {
		v := ErrorUnsupportedResponseType.ToValues()
		a.writeAuthRedirect(w, r, v)
		return
	}
	v := authorizer(a, r)
	a.writeAuthRedirect(w, r, v)
}

func (a *OAuth) AccessTokenRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowd.", http.StatusMethodNotAllowed)
		return
	}
	t := r.FormValue("grant_type")
	if t == "" {
		ErrorGrantTypeMissing.Write(w)
		return
	}

	access := a.config.AccessHandlers[GrantType(t)]
	if access == nil {
		ErrorUnsupportedGrantType.Write(w)
		return
	}
	o, scode := access(a, r)
	body, _ := json.Marshal(o)

	h := w.Header()
	h.Set("Content-Type", "application/json;charset=UTF-8")
	h.Set("Cache-Control", "no-store")
	h.Set("Pragma", "no-cache")
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

func (a *OAuth) ConnectClientToCode(code string, client Client) error {
	auth := a.config.Storage.GetAuthorize(code)
	if auth == nil {
		return errors.New("not found code")
	}
	auth.SetClient(client)
	return nil
}

func (a *OAuth) RedirectAuthorize(w http.ResponseWriter, r *http.Request, v url.Values) {
	url := a.config.AuthorizeEndpoint + "?" + v.Encode()
	http.Redirect(w, r, url, http.StatusFound)
}

func (a *OAuth) getAuthExpires() time.Time {
	return time.Now().Add(a.config.AuthorizeExpires)
}

func (a *OAuth) getTokenExpires() time.Time {
	return time.Now().Add(a.config.AccessTokenExpires)
}

func (a *OAuth) writeAuthRedirect(w http.ResponseWriter, r *http.Request, v url.Values) {
	url := a.config.RedirectEndpoint + "?" + v.Encode()
	http.Redirect(w, r, url, http.StatusFound)
}
