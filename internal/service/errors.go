package service

import "errors"

var (
	errTagName = errors.New("tag name is required")
	errTitle   = errors.New("title is required")
)
