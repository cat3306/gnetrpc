package gnetrpc

import "errors"

var (
	NotFoundMethod  = errors.New("not found method")
	NotFoundService = errors.New("not found service")
)
