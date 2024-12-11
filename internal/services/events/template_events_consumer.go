package events

import (
	"context"
	"errors"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"
	"github.com/rs/zerolog/log"
	"scm.x5.ru/x5m/go-backend/template/internal/config"
	"scm.x5.ru/x5m/go-backend/template/internal/repository/kafka/consumer"
)

// EventsConsumer encapsulates consumerGroup for particular Events
type EventsConsumer struct {
	cg *consumer.ConsumerGroup
}

// NewConsumer Создает потребителя сообщений о событиях заказов в рамках consumer-group
func NewConsumer(cfg config.Common) *EventsConsumer {
	cg, err := consumer.New(
		cfg,
		otelsarama.WrapConsumerGroupHandler(NewHandler()),
		consumer.WithOffsetsInitial(sarama.OffsetNewest))

	if err != nil {
		panic(fmt.Sprintf("failed to create consumer group: %s", err.Error()))
	}

	return &EventsConsumer{
		cg: cg,
	}
}

// Run Запускает потребителя сообщений
func (c *EventsConsumer) Run(ctx context.Context) {
	c.cg.Run(ctx)
}

// RunErrorHandler Запускает обработчик канала ошибок
func (c *EventsConsumer) RunErrorHandler(ctx context.Context) {
	logger := log.Ctx(ctx)
	for {
		select {
		case chErr, ok := <-c.cg.Errors():
			if !ok {
				logger.Info().Msg("[cg-error] error: chan closed")

				return
			}

			if errors.Is(chErr, sarama.ErrClosedConsumerGroup) {
				logger.Info().Msg("[cg-error] consumer group closed")

				return
			}

			logger.Info().Msgf("[cg-error] error: %s\n", chErr)
		case <-ctx.Done():
			logger.Info().Msgf("[cg-error] ctx closed: %s\n", ctx.Err().Error())

			return
		}
	}
}

// Close Закрывает потребителя, завершая все операции
func (c *EventsConsumer) Close() {
	c.cg.PauseAll()
	c.cg.Close()
}
