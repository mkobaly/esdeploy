package elastic

import "errors"

//ErrBadHTTPVerb is when...
var ErrBadHTTPVerb = errors.New("Unknown HTTP Verb")

var ErrBadJSON = errors.New("Invalid JSON document")

var ErrEmptyURL = errors.New("URL is empty")

type ErrSchemaChange struct {
	Message string
}

func (e ErrSchemaChange) Error() string {
	return e.Message
}
