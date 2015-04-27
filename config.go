package oguth

import (
	"net/http"
	"net/url"
	"time"
)

var DefaultConfig = NewConfig()

type AuthHandler func(*OAuth, *http.Request) url.Values
type AuthHandlers map[ResponseType]AuthHandler

type AccessHandler func(*OAuth, *http.Request)
type AccessHandlers map[GrantType]AccessHandler

type Config struct {
	Storage Storage

	AuthCodeExpires   time.Duration
	AuthCodeGenerator func() string

	AuthrizeEndpoint string
	TokenEndpoint    string
	AuthHandlers     AuthHandlers
	AccessHandlers   AccessHandlers
}

func NewConfig() Config {
	return Config{
		Storage:           NewMemoryStorage(),
		AuthCodeExpires:   time.Minute,
		AuthCodeGenerator: DefaultAuthCodeGenerator,
		AuthHandlers: AuthHandlers{
			ResponseCode: authorizeRequestCode,
		},
		AccessHandlers: AccessHandlers{
			GrantAuthCode: accessTokenRequestAuthCode,
		},
	}
}
