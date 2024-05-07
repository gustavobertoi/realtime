package channels

import "errors"

var (
	errInvalidMaxLimitOfConnections = errors.New("invalid max limit of connections for this channel")
	errInvalidChannelType           = errors.New("invalid channel type")
)
