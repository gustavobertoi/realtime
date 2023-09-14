package channels

import "errors"

var (
	errMessageDoesNotExist          = errors.New("message does not exists in store")
	errMessageAlreadyPublished      = errors.New("message already published")
	errMaxLimitOfConnections        = errors.New("max limit of connections for this client")
	errInvalidMaxLimitOfConnections = errors.New("invalid max limit of connections for this channel")
)
