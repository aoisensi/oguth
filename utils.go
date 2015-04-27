package oguth

import (
	"math/rand"
	"net/http"
	"time"
)

type httpHandler func(http.ResponseWriter, *http.Request)

var random *rand.Rand

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}
