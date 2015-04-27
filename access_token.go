package oguth

import "net/http"

type GrantType string

const (
	GrantRefreshToken      GrantType = "refresh_token"
	GrantAuthCode                    = "authorization_code"
	GrantPassword                    = "password"
	GrantClientCredentials           = "client_credentials"
)

func accessTokenRequestAuthCode(a *OAuth, r *http.Request) {

}
