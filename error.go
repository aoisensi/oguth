package oguth

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type ErrorCode string

const (
	ErrorCodeInvalidRequest          ErrorCode = "invalid_request"
	ErrorCodeUnauthorizedClient                = "unauthorized_client"
	ErrorCodeUnsupportedGrantType              = "unsupported_grant_type"
	ErrorCodeUnsupportedResponseType           = "unsupported_response_type"
	ErrorCodeAccessDenied                      = "access_denied"
	ErrorCodeInvalidScope                      = "invalid_scope"
	ErrorCodeServerError                       = "server_error"
	ErrorCodeTemporarilyUnavailable            = "temporarily_unavailable"
)

type Error struct {
	Error            ErrorCode `json:"error"`
	ErrorDescription string    `json:"error_description"`
	ErrorUri         string    `json:"error_uri"`
	State            string    `json:"state"`
}

func NewError(code ErrorCode) Error {
	return Error{
		Error: code,
	}
}

func (e Error) ToValues() url.Values {
	v := url.Values{
		"error":             {string(e.Error)},
		"error_description": {e.ErrorDescription},
		"error_uri":         {e.ErrorUri},
	}
	if e.State != "" {
		v.Set("state", e.State)
	}
	return v
}

func (e Error) Write(w http.ResponseWriter) {
	body, _ := json.Marshal(&e)
	w.Write(body)
}
