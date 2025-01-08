package client

import (
	"encoding/base64"
	"fmt"
	"net/http"
)

type BearerAuth struct {
	Token string
}

func (a BearerAuth) Apply(req *http.Request) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

type BasicAuth struct {
	Username string
	Password string
}

func (a BasicAuth) Apply(req *http.Request) {
	auth := base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf("%s:%s", a.Username, a.Password)),
	)
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", auth))
}
