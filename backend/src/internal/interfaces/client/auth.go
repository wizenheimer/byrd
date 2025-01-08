package interfaces

import "net/http"

type AuthMethod interface {
	Apply(req *http.Request)
}
