package oguth

import (
	"net/http"
	"net/url"
	"time"
)

type ResponseType string

const (
	ResponseCode ResponseType = "code"
)

func newAuthorizeForm(r *http.Request) *authorizeForm {
	q := r.URL.Query()
	f := new(authorizeForm)
	f.ResponseType = ResponseType(q.Get("response_type"))
	f.ClientId = q.Get("client_id")
	f.State = q.Get("state")
	f.Scope = q.Get("scope")
	f.RedirectUri = q.Get("redirect_uri")
	return f
}

type authorizeForm struct {
	ResponseType ResponseType
	ClientId     string
	State        string
	Scope        string
	RedirectUri  string
}

func authorizeRequestCode(a *OAuth, r *http.Request) url.Values {
	f := newAuthorizeForm(r)
	if f.ClientId == "" {
		e := NewError(ErrorCodeInvalidRequest)
		e.State = f.State
		return e.ToValues()
	}
	client := a.owner.GetClient(f.ClientId)
	if client == nil {
		e := NewError(ErrorCodeUnauthorizedClient)
		e.State = f.State
		return e.ToValues()
	}
	code := a.config.AuthCodeGenerator()
	auth := authorize{
		id:      f.ClientId,
		expires: time.Now().Add(a.config.AuthCodeExpires),
		uri:     f.RedirectUri,
	}
	a.config.Storage.SetAuthorize(code, auth)
	v := url.Values{"code": {code}}
	if f.State != "" {
		v.Set("state", f.State)
	}
	return v
}

func DefaultAuthCodeGenerator() (code string) {
	return SimpleRandomTokenGenerator(32)
}
