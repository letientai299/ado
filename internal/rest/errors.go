package rest

import (
	"fmt"
)

var (
	ErrNotFound       = adoErr{code: 404, msg: "not found"}
	ErrAuth           = adoErr{code: 203, msg: `PAT might be expired`}
	ErrForbidden      = adoErr{code: 403, msg: `Entra token is created with wrong tenant ID`}
	ErrInternalServer = adoErr{msg: "ADO server internal error"}
	ErrInvalidRequest = adoErr{msg: "invalid request"}
)

// adoErr represents an error specific to the ADO REST API, the underlying value
type adoErr struct {
	code int
	msg  string
	uri  string
}

func (e adoErr) Error() string {
	return fmt.Sprintf("%s, code=%d, uri=%s", e.msg, e.code, e.uri)
}

func (e adoErr) with(update func(cloned *adoErr)) adoErr {
	c := adoErr{code: e.code, msg: e.msg, uri: e.uri}
	update(&c)
	return c
}

func (e adoErr) WithCode(code int) adoErr  { return e.with(func(c *adoErr) { c.code = code }) }
func (e adoErr) WithURI(uri string) adoErr { return e.with(func(c *adoErr) { c.uri = uri }) }
