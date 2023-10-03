package config

import "errors"

var (
	errYamlConfigNotDefined         = errors.New("yaml config is not created")
	errDriverNotSupported           = errors.New("pubsub driver is not supported in yaml spec")
	errRedisPubSubAdapterNotDefined = errors.New("redis pubsub adapter is not defined in yaml spec")
)
