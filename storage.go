package oguth

import "time"

type Authorize interface {
	GetClientId() string
	GetExpires() time.Time
	GetRedirectUri() string
	SetClient(Client)
}

type AccessToken interface {
	GetClient() Client
	GetExpires() time.Time
}

type Storage interface {
	AddAuthorize(code string, auth Authorize)
	GetAuthorize(code string) Authorize
	DisableAuthorize(id string)

	GetAccessToken(token string) AccessToken
	AddAccessToken(token string, access AccessToken)
}
