package oguth

import "time"

type AuthCode interface {
	GetClientId() string
	GetCode() string
	GetExpires() time.Time
	GetRedirectUri() string
}

type Storage interface {
	SetAuthCode(auth AuthCode)
	GetAuthCode(code string) AuthCode
	DisableAuthCode(id string)
}

func NewMemoryStorage() Storage {
	return &memoryStorage{
		codes: make(map[string]authCode),
	}
}

type memoryStorage struct {
	codes map[string]authCode
}

func (s *memoryStorage) SetAuthCode(auth AuthCode) {
	s.codes[auth.GetCode()] = authCode{
		id:      auth.GetClientId(),
		code:    auth.GetCode(),
		expires: auth.GetExpires(),
	}
}

func (s *memoryStorage) GetAuthCode(code string) AuthCode {
	v, ok := s.codes[code]
	if !ok {
		return nil
	}
	if time.Now().After(v.expires) {
		s.DisableAuthCode(code)
		return nil
	}
	return v
}

func (s *memoryStorage) DisableAuthCode(id string) {
	delete(s.codes, id)
}

type authCode struct {
	AuthCode
	id, code, uri string
	expires       time.Time
}

func (c authCode) GetClientId() string {
	return c.id
}

func (c authCode) GetCode() string {
	return c.code
}

func (c authCode) GetExpires() time.Time {
	return c.expires
}

func (c authCode) GetRedirectUri() string {
	return c.uri
}
