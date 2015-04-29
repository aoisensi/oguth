package oguth

import "strings"

type Scopes []string

func (scopes Scopes) ContainScope(scope string) bool {
	for _, s := range scopes {
		if s == scope {
			return true
		}
	}
	return false
}

func (root Scopes) ContainScopes(scopes Scopes) bool {
	for _, s := range scopes {
		if !root.ContainScope(s) {
			return false
		}
	}
	return true
}

func ParseScopes(s string) Scopes {
	strs := make([]string, 0)
	for _, c := range strings.Split(s, " ") {
		if c != "" {
			strs = append(strs, c)
		}
	}
	return Scopes(strs)
}

func (s Scopes) Available(oauth *OAuth) *Error {
	if oauth.config.AvailableScopes.ContainScopes(s) {
		return nil
	}

	return ErrorUnavailableScope
}
