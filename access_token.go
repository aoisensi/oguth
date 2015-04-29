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
	return f
}

type tokenRequestForm struct {
	GrantType   GrantType
	Code        string
	RedirectUri string
	ClientId    string
	Username    string
	Password    string
}

type accessTokenResponse struct {
	AccessToken  string    `json:"access_token"`
	TokenType    TokenType `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
	RefreshToken string    `json:"refresh_token"`
}

func accessTokenRequestAuthCode(a *OAuth, r *http.Request) (interface{}, int) {
	f := newTokenRequestForm(r)
	if f.Code == "" {
		e := NewError(ErrorCodeInvalidRequest)
		return e, http.StatusBadRequest
	}
	//TODO
	return nil, http.StatusAccepted
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
	cli := a.owner.GetClientWithPasswordGrant(username, password)
	if cli == nil {
		e := NewError(ErrorCodeUnauthorizedClient)
		return e, http.StatusBadRequest
	}

	token := a.config.AccessTokenGenerator()
	access := &accessToken{
		client:  cli,
		expires: a.getTokenExpires(),
	}
	a.config.Storage.AddAccessToken(token, access)
	body := accessTokenResponse{
		TokenType:    a.config.TokenType,
		AccessToken:  token,
		ExpiresIn:    a.config.accessTokenExpiresInt,
		RefreshToken: "test",
	}
	return body, http.StatusAccepted
}

func DefaultAccessTokenGenerator() (code string) {
	return SimpleRandomTokenGenerator(40)
}
