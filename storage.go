package oguth

import "time"

type Authorize interface {
	GetClientId() string
	GetExpires() time.Time
	GetRedirectUri() string
}

type AccessToken interface {
	GetClient() Client
}

type Storage interface {
	SetAuthorize(code string, auth Authorize)
	GetAuthorize(code string) Authorize
	DisableAuthorize(id string)

	GetAccessToken(token string) AccessToken
	SetAccessToken(token string, access AccessToken)
}
