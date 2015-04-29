package oguth

type ResourceOwner interface {
	ExistClientId(id string) bool
	GetClient(id, secret string) Client
	GetClientWithPasswordGrant(username, password string) Client
	Close()
}

func (a *OAuth) existClientId(id string) *Error {
	if a.owner.ExistClientId(id) {
		return nil
	}
	return ErrorClientNotFound
}
