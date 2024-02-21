package storage

import "errors"

var (
	ErrDoNotExists = errors.New("id do not exists")
	ErrExpired     = errors.New("value expired")
)
