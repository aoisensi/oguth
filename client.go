package oguth

type Client interface {
	GetClientID() string
	GetRedirectURI() string
}
