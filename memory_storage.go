package oguth

import "time"

func NewMemoryStorage() Storage {
	return &memoryStorage{
		authes:  make(map[string]Authorize),
		access:  make(map[string]AccessToken),
		refresh: make(map[string]RefreshToken),
	}
}

type memoryStorage struct {
	authes  map[string]Authorize
	access  map[string]AccessToken
	refresh map[string]RefreshToken
}

func (s *memoryStorage) AddAuthorize(code string, auth Authorize) {
	s.authes[code] = auth
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

func (s *memoryStorage) AddAccessToken(token string, access AccessToken) {
	s.access[token] = access
}

func (s *memoryStorage) GetAccessToken(token string) AccessToken {
	return s.access[token]
}

func (s *memoryStorage) AddRefreshToken(token string, refresh RefreshToken) {
	s.refresh[token] = refresh
}

func (s *memoryStorage) SetRefreshToken(token string) RefreshToken {
	return s.refresh[token]
}

type authorize struct {
	id, uri string
	expires time.Time
	client  Client
}

type accessToken struct {
	client  Client
	expires time.Time
}

func (c *authorize) GetClientId() string {
	return c.id
}

func (c *authorize) GetExpires() time.Time {
	return c.expires
}

func (c *authorize) GetRedirectUri() string {
	return c.uri
}

func (c *authorize) GetClient() Client {
	return c.client
}

func (c *authorize) SetClient(client Client) {
	c.client = client
}

func (c *accessToken) GetClient() Client {
	return c.client
}

func (c *accessToken) GetExpires() time.Time {
	return c.expires
}
