package interfaces

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

type IProcessor interface {
	Cache() ICacheProcessor
	Queue() IQueueProcessor
	APIs() IAPIsProcessor
}

// Cache

type ICacheProcessor interface {
	Set(ctx context.Context, key string, value []byte) error
	SetJSON(ctx context.Context, key string, value interface{}) error
	SetJSONWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	GetJSON(ctx context.Context, key string, v interface{}) error
	Delete(ctx context.Context, key string) error
	FlushAll(ctx context.Context) error
}

type ICacheProvider interface {
	Set(ctx context.Context, key string, value []byte) error
	SetWithTTL(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
	FlushAll(ctx context.Context) error
	Close() error
}

// Queue

type IQueueProcessor interface {
	Producers() IQueueProducersProcessor
}

type IQueueProducersProcessor interface {
}

type IQueueProducerProcessor interface {
}

type IQueueProducerProvider interface {
	Send(ctx context.Context, qName string, message []byte, args ...amqp.Table) (err error)
}

// APIs

type IAPIsProcessor interface {
}
