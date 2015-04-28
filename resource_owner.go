package oguth

type ResourceOwner interface {
	GetClient(id string) Client
	GetClientWithPasswordGrant(username, password string) Client
	Close()
}
