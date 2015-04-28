package oguth

type Scope string
type Scopes []Scope

func (sc Scopes) ContainScope(s Scope) bool {
	for _, c := range sc {
		if c == s {
			return true
		}
	}
	return false
}

func (root Scopes) ContainScopes(ss Scopes) bool {
	for _, s := range ss {
		if !root.ContainScope(s) {
			return false
		}
	}
	return true
}

func ParseScopes(s string) Scopes {
	//TODO
	//return Scopes([]Scope(strings.Split(s, " ")))
	return make(Scopes, 0)
}
