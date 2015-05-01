package oguth

import "time"

type Authorize interface {
	GetClientId() string
	GetExpires() time.Time
	GetRedirectUri() string
	SetClient(Client)
	GetClient() Client
}

type AccessToken interface {
	GetClient() Client
	GetExpires() time.Time
}

type RefreshToken interface {
	GetClient() Client
}

type Storage interface {
	AddAuthorize(code string, auth Authorize)
	GetAuthorize(code string) Authorize
	DisableAuthorize(id string)

	AddAccessToken(token string, access AccessToken)
	GetAccessToken(token string) AccessToken

	AddRefreshToken(token string, refresh RefreshToken)
	GetRefreshToken(token string) RefreshToken
	DisableRefreshToken(token string)
}
