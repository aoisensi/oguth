package oguth

import "net/http"

type GrantType string
type TokenType string

const (
	GrantRefreshToken      GrantType = "refresh_token"
	GrantAuthCode                    = "authorization_code"
	GrantPassword                    = "password"
	GrantClientCredentials           = "client_credentials"
)

const (
	TokenTypeBearer TokenType = "Bearer"
	TokenTypeMAC              = "MAC"
)

func newTokenRequestForm(r *http.Request) *tokenRequestForm {
	q := r.PostForm
	f := new(tokenRequestForm)
	f.GrantType = GrantType(q.Get("grant_type"))
	f.Code = q.Get("code")
	f.RedirectUri = q.Get("redirect_uri")
	f.ClientId = q.Get("client_id")
	f.Username = q.Get("username")
	f.Password = q.Get("password")
	f.RefreshToken = q.Get("refresh_token")
	return f
}

type tokenRequestForm struct {
	GrantType    GrantType
	Code         string
	RedirectUri  string
	ClientId     string
	Username     string
	Password     string
	RefreshToken string
}

type accessTokenResponse struct {
	AccessToken  string    `json:"access_token"`
	TokenType    TokenType `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
	RefreshToken string    `json:"refresh_token,omitempty"`
}

func accessTokenRequestAuthCode(a *OAuth, r *http.Request) (interface{}, int) {
	f := newTokenRequestForm(r)
	if f.Code == "" {
		e := NewError(ErrorCodeInvalidRequest)
		return e, http.StatusBadRequest
	}
	auth := a.config.Storage.GetAuthorize(f.Code)
	if auth == nil {
		e := NewError(ErrorCodeInvalidRequest)
		return e, http.StatusBadRequest
	}

	return a.newAccessTokenResponse(auth.GetClient(), true)
}

func accessTokenRequestPassowrd(a *OAuth, r *http.Request) (interface{}, int) {
	f := newTokenRequestForm(r)
	username, password, ok := r.BasicAuth()
	if !ok || username == "" {
		username = f.Username
		password = f.Password
	}
	if username == "" || password == "" {
		e := NewError(ErrorCodeInvalidRequest)
		return e, http.StatusBadRequest
	}
	if !ok || (username != f.Username) || (password != f.Password) {
		e := NewError(ErrorCodeInvalidRequest)
		return e, http.StatusBadRequest
	}
	client := a.owner.GetClientWithPasswordGrant(username, password)

	return a.newAccessTokenResponse(client, true)
}

func accessTokenRequestClient(a *OAuth, r *http.Request) (interface{}, int) {
	clientId, clientSecret, ok := r.BasicAuth()
	if !ok {
		return nil, 500
	}
	client := a.owner.GetClient(clientId, clientSecret)

	return a.newAccessTokenResponse(client, false)
}

func accessTokenRefresh(a *OAuth, r *http.Request) (interface{}, int) {
	f := newTokenRequestForm(r)
	token := a.config.Storage.GetRefreshToken(f.RefreshToken)
	if token == nil {
		return ErrorRefreshTokenInvalid, http.StatusBadRequest
	}
	a.config.Storage.DisableRefreshToken(f.RefreshToken)
	return a.newAccessTokenResponse(token.GetClient(), true)
}

func (a *OAuth) newAccessTokenResponse(client Client, genRefresh bool) (interface{}, int) {
	if client == nil {
		return ErrorClientNotFound, http.StatusBadRequest
	}
	token := a.config.AccessTokenGenerator()
	access := &accessToken{
		client:  client,
		expires: a.getTokenExpires(),
	}
	a.config.Storage.AddAccessToken(token, access)
	var refresh string
	if genRefresh {
		refresh := a.config.RefreshTokenGenerator()

		a.config.Storage.AddRefreshToken(refresh, &refreshToken{
			client: client,
		})
	}
	body := accessTokenResponse{
		TokenType:    a.config.TokenType,
		AccessToken:  token,
		ExpiresIn:    a.config.accessTokenExpiresInt,
		RefreshToken: refresh,
	}
	return body, http.StatusAccepted
}

func DefaultAccessTokenGenerator() (code string) {
	return SimpleRandomTokenGenerator(40)
}
