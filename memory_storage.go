package oguth

import "time"

func NewMemoryStorage() Storage {
	return &memoryStorage{
		authes: make(map[string]Authorize),
		access: make(map[string]AccessToken),
	}
}

type memoryStorage struct {
	authes map[string]Authorize
	access map[string]AccessToken
}

func (s *memoryStorage) SetAuthorize(code string, auth Authorize) {
	s.authes[code] = authorize{
		id:      auth.GetClientId(),
		expires: auth.GetExpires(),
	}
}

func (s *memoryStorage) GetAuthorize(code string) Authorize {
	v, ok := s.authes[code]
	if !ok {
		return nil
	}
	if time.Now().After(v.GetExpires()) {
		s.DisableAuthorize(code)
		return nil
	}
	return v
}

func (s *memoryStorage) DisableAuthorize(code string) {
	delete(s.authes, code)
}

func (s *memoryStorage) SetAccessToken(token string, access AccessToken) {
	s.access[token] = access
}

func (s *memoryStorage) GetAccessToken(token string) AccessToken {
	return s.access[token]
}

type authorize struct {
	id, uri string
	expires time.Time
}

type accessToken struct {
	client Client
}

func (c authorize) GetClientId() string {
	return c.id
}

func (c authorize) GetExpires() time.Time {
	return c.expires
}

func (c authorize) GetRedirectUri() string {
	return c.uri
}

func (c accessToken) GetClient() Client {
	return c.client
}
