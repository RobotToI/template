package kafka

import (
	"context"
	"fmt"
	"time"

	"scm.x5.ru/x5m/go-backend/template/internal/config"
)

// Consumer Queue consumer
type Consumer interface {
	// GetName returns consumer name
	GetName() string
	// Consume process queue message
	Consume(ctx context.Context, queue string, since time.Duration) error
	// IsFailProof if false, consumer will be stopped on error
	IsFailProof() bool
	// GetTopic returns topic stored in consumer configuration
	GetTopic() string
}

// Client Queue client
// type Client interface {
// 	GetClient() (sarama.Client, error)
// 	GetGroup() string
// 	Close() error
// }

// NewQueueClient creates client depending on provider type
func NewQueueClient(ctx context.Context, cfg config.Common, provider string) (*Client, error) {
	switch provider {
	case "kafka":
		return NewClient(ctx, cfg.Kafka)
	default:
		return nil, fmt.Errorf("unrecognize message provider: %s", provider)
	}
}
