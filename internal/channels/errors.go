package channels

import "errors"

var (
	errMessageDoesNotExist          = errors.New("message does not exists in store")
	errInvalidMaxLimitOfConnections = errors.New("invalid max limit of connections for this channel")
	errInvalidChannelType           = errors.New("invalid channel type")
)
