package rcon_go

import (
	"errors"
)

var (
	ErrNoHostSpecified = errors.New("no host specified in config")

	ErrFailedToConnect = errors.New("failed to connect to host")

	ErrPayloadLimitExceeded = errors.New("payload exceeds max payloadload limit")

	ErrBadAuth = errors.New("bad auth, could not authenticate")
)
