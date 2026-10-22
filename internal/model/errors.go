package model

import "errors"

var (
	ErrNotFound = errors.New("url not found")

	ErrIDExists = errors.New("url id already exists")
)
