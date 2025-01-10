// ./src/internal/interfaces/client/auth.go
package interfaces

import "net/http"

type AuthMethod interface {
	Apply(req *http.Request)
}
