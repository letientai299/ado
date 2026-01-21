package rest

import (
	"fmt"
)

// Sentinel errors for common Azure DevOps API error conditions.
// Use errors.Is() to check for these error types.
var (
	// ErrNotFound indicates the requested resource was not found (HTTP 404).
	ErrNotFound = adoErr{code: 404, msg: "not found"}

	// ErrAuth indicates an authentication failure (HTTP 203).
	// This typically means the PAT has expired or is invalid.
	ErrAuth = adoErr{code: 203, msg: "PAT might be expired"}

	// ErrForbidden indicates the user lacks permission (HTTP 403).
	// For Entra tokens, this may indicate the wrong tenant ID was used.
	ErrForbidden = adoErr{code: 403, msg: "Entra token is created with wrong tenant ID"}

	// ErrInternalServer indicates an Azure DevOps server error (HTTP 5xx).
	ErrInternalServer = adoErr{msg: "ADO server internal error"}

	// ErrInvalidRequest indicates a client error (HTTP 4xx other than 403, 404).
	ErrInvalidRequest = adoErr{msg: "invalid request"}
)

// adoErr represents an error from the Azure DevOps REST API.
// It includes the HTTP status code and request URI for debugging.
type adoErr struct {
	code int
	msg  string
	uri  string
}

// Error implements the error interface.
func (e adoErr) Error() string {
	return fmt.Sprintf("%s, code=%d, uri=%s", e.msg, e.code, e.uri)
}

// with creates a copy of the error with the specified modification.
func (e adoErr) with(update func(cloned *adoErr)) adoErr {
	c := adoErr{code: e.code, msg: e.msg, uri: e.uri}
	update(&c)
	return c
}

// WithCode returns a copy of the error with the specified HTTP status code.
func (e adoErr) WithCode(code int) adoErr { return e.with(func(c *adoErr) { c.code = code }) }

// WithURI returns a copy of the error with the specified request URI.
func (e adoErr) WithURI(uri string) adoErr { return e.with(func(c *adoErr) { c.uri = uri }) }
