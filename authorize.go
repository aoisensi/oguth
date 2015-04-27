package oguth

import (
	"encoding/base64"
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
	auth := authCode{
		id:      f.ClientId,
		code:    code,
		expires: time.Now().Add(a.config.AuthCodeExpires),
		uri:     f.RedirectUri,
	}
	a.config.Storage.SetAuthCode(auth)
	v := url.Values{"code": {code}}
	if f.State != "" {
		v.Set("state", f.State)
	}
	return v
}

func DefaultAuthCodeGenerator() (code string) {
	size := 16 + random.Intn(16)
	body := make([]byte, size)
	for i := range body {
		body[i] = byte(random.Intn(255))
	}
	return base64.StdEncoding.EncodeToString(body)
}
