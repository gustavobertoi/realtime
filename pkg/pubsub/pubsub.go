package pubsub

import (
	"context"
	"fmt"

	"github.com/gustavobertoi/realtime/internal/channels"
	"github.com/gustavobertoi/realtime/internal/dtos"
)

type (
	PubSub struct {
		Driver   string
		Consumer channels.ConsumerAdapter
		Producer channels.ProducerAdapter
	}
	PubSubError struct {
		Driver  string
		Message string
	}
)

func NewPubSub(ctx context.Context, dto *dtos.PubSubConfig) (*PubSub, error) {
	switch dto.Driver {
	case MemoryDriver:
		return NewMemoryAdapter()
	case RedisDriver:
		return NewRedisAdapter(ctx, dto.Redis)
	case UpstashKafkaDriver:
		return NewUpstashKafkaAdapter(ctx, dto.UpstashKafka)
	default:
		return nil, nil
	}
}

func NewPubSubError(driver, message string) *PubSubError {
	return &PubSubError{
		Driver:  driver,
		Message: message,
	}
}

func (e *PubSubError) Error() string {
	return fmt.Sprintf("pubsub error: driver=%s, message=%s", e.Driver, e.Message)
}
