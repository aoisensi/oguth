package oguth

import (
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const (
	vschars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type httpHandler func(http.ResponseWriter, *http.Request)

var random *rand.Rand

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func SimpleRandomTokenGenerator(size int) string {
	body := make([]rune, size)
	for i := range body {
		r := vschars[random.Intn(61)]
		body[i] = rune(r)
	}
	return string(body)
}

func httpHeaderAuth(r *http.Request) string {
	v := r.Header.Get("Authorization")
	if v == "" {
		return ""
	}
	i := strings.Index(v, " ")
	if i < 0 {
		return ""
	}
	t := v[i+1:]
	switch TokenType(v[:i]) {
	case TokenTypeBearer:
		return httpHeaderAuthBearer(t)
	}
	return ""
}

func httpHeaderAuthBearer(token string) string {
	return token
}
