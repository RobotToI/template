package consumer

import (
	"context"
	"errors"
	"time"

	"github.com/IBM/sarama"
	"scm.x5.ru/x5m/go-backend/packages/zlogger"
	"scm.x5.ru/x5m/go-backend/template/internal/config"
)

// ConsumerGroup encapsulates sarama.ConsumerGroup & handlers for topics
type ConsumerGroup struct {
	sarama.ConsumerGroup
	handler sarama.ConsumerGroupHandler
	topics  []string
}

// New Создает новую группу потребителей сообщений
func New(cfg config.Common, consumerGroupHandler sarama.ConsumerGroupHandler, opts ...Option) (*ConsumerGroup, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.MaxVersion
	/*
		sarama.OffsetNewest - получаем только новые сообщений, те, которые уже были игнорируются
		sarama.OffsetOldest - читаем все с самого начала
	*/
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	// Используется, если ваш offset "уехал" далеко и нужно пропустить невалидные сдвиги
	saramaConfig.Consumer.Group.ResetInvalidOffsets = true
	// Сердцебиение консьюмера
	saramaConfig.Consumer.Group.Heartbeat.Interval = 3 * time.Second
	// Таймаут сессии
	saramaConfig.Consumer.Group.Session.Timeout = 60 * time.Second
	// Таймаут ребалансировки
	saramaConfig.Consumer.Group.Rebalance.Timeout = 60 * time.Second
	//
	saramaConfig.Consumer.Return.Errors = true

	saramaConfig.Consumer.Offsets.AutoCommit.Enable = false
	saramaConfig.Consumer.Offsets.AutoCommit.Interval = 5 * time.Second

	// Применяем свои конфигурации
	for _, opt := range opts {
		opt.Apply(saramaConfig)
	}

	cg, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers, cfg.Kafka.GroupID, saramaConfig)
	if err != nil {
		return nil, err
	}

	return &ConsumerGroup{
		ConsumerGroup: cg,
		handler:       consumerGroupHandler,
		topics:        cfg.Kafka.Topics,
	}, nil
}

// Run Запускает потребителей сообщений в рамках consumer-group
func (c *ConsumerGroup) Run(ctx context.Context) {
	logger := zlogger.LoadOrCreateFromCtx(ctx)
	logger.Debug().Msg("[consumer-group] run")

	for {
		if err := ctx.Err(); err != nil {
			logger.Error().Err(err).Msg("[consumer-group] ctx err")
		}

		if err := c.ConsumerGroup.Consume(ctx, c.topics, c.handler); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				logger.Error().Err(err).Msg("[consumer-group] consumer group closed")

				return
			}

			logger.Error().Err(err).Msg("[consumer-group] failed to consume")
		}
	}
}
