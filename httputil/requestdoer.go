package httputil

import (
	"net/http"
)

// RequestDoer is the interface implemented by types that
// can Do an HTTP request.
type RequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}
