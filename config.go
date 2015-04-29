package oguth

import (
	"net/http"
	"net/url"
	"time"
)

var DefaultConfig = NewConfig()

type AuthHandler func(*OAuth, *http.Request) url.Values
type AuthHandlers map[ResponseType]AuthHandler

type AccessHandler func(*OAuth, *http.Request) (interface{}, int)
type AccessHandlers map[GrantType]AccessHandler

type Config struct {
	Storage Storage
	Owner   ResourceOwner

	AuthorizeExpires      time.Duration
	authorizeExpiresInt   int
	AccessTokenExpires    time.Duration
	accessTokenExpiresInt int
	AuthorizeGenerator    func() string
	AccessTokenGenerator  func() string

	AuthorizeEndpoint   string
	AccessTokenEndpoint string
	RedirectEndpoint    string
	AuthHandlers        AuthHandlers
	AccessHandlers      AccessHandlers
	TokenType           TokenType

	AvailableScopes Scopes
}

func NewConfig() Config {
	return Config{
		AuthorizeExpires:     time.Hour,
		AccessTokenExpires:   time.Hour,
		AuthorizeGenerator:   DefaultAuthCodeGenerator,
		AccessTokenGenerator: DefaultAccessTokenGenerator,
		TokenType:            TokenTypeBearer,
		AuthHandlers: AuthHandlers{
			ResponseCode:  authorizeRequestCode,
			ResponseToken: authorizeRequestToken,
		},
		AccessHandlers: AccessHandlers{
			GrantAuthCode: accessTokenRequestAuthCode,
			GrantPassword: accessTokenRequestPassowrd,
		},
	}
}

func (c *Config) init() {
	c.authorizeExpiresInt = int(c.AuthorizeExpires.Minutes())
	c.accessTokenExpiresInt = int(c.AccessTokenExpires.Minutes())
	if c.Storage == nil {
		c.Storage = NewMemoryStorage()
	}
}
