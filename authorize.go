package oguth

import (
	"net/http"
	"net/url"
)

type ResponseType string

const (
	ResponseCode  ResponseType = "code"
	ResponseToken              = "token"
)

func newAuthorizeForm(r *http.Request) (*authorizeForm, *Error) {
	q := r.Form
	f := new(authorizeForm)
	f.ResponseType = ResponseType(q.Get("response_type"))
	f.ClientId = q.Get("client_id")
	f.State = q.Get("state")
	f.Scope = ParseScopes(q.Get("scope"))
	f.RedirectUri = q.Get("redirect_uri")
	if f.ClientId == "" {
		return f, ErrorClientIdMissing
	}
	return f, nil
}

type authorizeForm struct {
	ResponseType ResponseType
	ClientId     string
	State        string
	Scope        Scopes
	RedirectUri  string
}

func authorizeRequestCode(a *OAuth, w http.ResponseWriter, r *http.Request) {
	f, err := newAuthorizeForm(r)
	if err != nil {
		http.Error(w, "missing params", http.StatusBadRequest)
	}
	rurl := a.owner.GetRedirectUri(f.ClientId)

	if err := f.Scope.Available(a); err != nil {
		err.State = f.State
		http.Redirect(w, r, rurl+"?"+err.ToValues().Encode(), http.StatusFound)
		return
	}

	ok := a.owner.ExistClientId(f.ClientId)
	if !ok {
		e := NewError(ErrorCodeUnauthorizedClient)
		e.State = f.State
		http.Redirect(w, r, rurl+"?"+e.ToValues().Encode(), http.StatusFound)
		return
	}

	client := a.owner.AuthCodeDecision(r, f.ClientId)
	if client == nil {
		a.owner.AuthCodeMissing(w, r)
		return
	}
	code := a.config.AuthorizeGenerator()
	auth := &authorize{
		id:      f.ClientId,
		expires: a.getAuthExpires(),
		uri:     f.RedirectUri,
		client:  client,
	}
	a.config.Storage.AddAuthorize(code, auth)
	v := url.Values{"code": {code}}
	if f.State != "" {
		v.Set("state", f.State)
	}
	http.Redirect(w, r, rurl+"?"+v.Encode(), http.StatusFound)
	return

}

func authorizeRequestToken(a *OAuth, w http.ResponseWriter, r *http.Request) {
	//TODO
	/*
		var f authorizeForm
		if f, err := newAuthorizeForm(r); err != nil {
			e := NewError(ErrorCodeUnauthorizedClient)
			e.State = f.State
			return e.ToValues()
		}
		if err := f.Scope.Available(a); err != nil {
			err.State = f.State
			return err.ToValues()
		}
		if err := a.existClientId(f.ClientId); err != nil {
			err.State = f.State
			return err.ToValues()
		}
		token := a.config.AccessTokenGenerator()
		access := accessToken{
			client:  nil,
			expires: a.getTokenExpires(),
		}
		a.config.Storage.AddAccessToken(token, access)
		v := url.Values{
			"access_token": {token},
			"token_type":   {string(a.config.TokenType)},
		}
		if f.State != "" {
			v.Set("state", f.State)
		}
	*/
	return
}

func DefaultAuthCodeGenerator() (code string) {
	return SimpleRandomTokenGenerator(32)
}
