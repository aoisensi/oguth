package oguth

import "net/http"

type ResourceOwner interface {
	ExistClientId(id string) bool
	GetClient(id, secret string) Client
	GetClientWithPasswordGrant(username, password string) Client
	GetRedirectUri(clientId string) string
	Close()
	AuthCodeDecision(r *http.Request, clientId string) Client
	AuthCodeMissing(w http.ResponseWriter, r *http.Request)
}

func (a *OAuth) existClientId(id string) *Error {
	if a.owner.ExistClientId(id) {
		return nil
	}
	return ErrorClientNotFound
}
