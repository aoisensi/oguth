package oguth

import "net/http"

type OAuth struct {
	owner  ResourceOwner
	config Config
}

func NewOAuth(config Config, owner ResourceOwner) *OAuth {
	return &OAuth{config: config, owner: owner}
}

func (a *OAuth) Close() {
	a.owner.Close()
}

func (a *OAuth) AuthorizeRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowd.", http.StatusMethodNotAllowed)
		return
	}
	endpoint := a.config.TokenEndpoint
	t := ResponseType(r.FormValue("response_type"))
	authorizer := a.config.AuthHandlers[t]
	if authorizer == nil {
		v := NewError(ErrorCodeUnsupportedResponseType).ToValues()
		http.Redirect(w, r, endpoint+"?"+v.Encode(), http.StatusFound)
		return
	}
	v := authorizer(a, r)
	http.Redirect(w, r, endpoint+"?"+v.Encode(), http.StatusFound)
}

func (a *OAuth) AccessTokenRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowd.", http.StatusMethodNotAllowed)
		return
	}
	t := GrantType(r.FormValue("grant_type"))
	access := a.config.AccessHandlers[t]
	if access == nil {
		http.Error(w, "error. this code is not yet", http.StatusBadRequest)
		return
	}

}
