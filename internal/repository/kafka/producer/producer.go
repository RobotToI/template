package producer

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"
	"scm.x5.ru/x5m/go-backend/template/internal/repository/kafka"
)

// KafkaProducer encapsulates sarama.SyncProducer, could be changed to Async one.
type KafkaProducer struct {
	producer sarama.SyncProducer
}

// New creates KafkaProducer by calling NewSyncProducerFromClient
func New(client kafka.Client) *KafkaProducer {
	kafkaClient, err := client.GetClient()
	if err != nil {
		errors.Wrap(err, "error getting kafka client")
	}
	kafkaProducer, err := sarama.NewSyncProducerFromClient(kafkaClient)
	if err != nil {
		errors.Wrap(err, "error getting kafka producer")
	}
	return &KafkaProducer{producer: kafkaProducer}
}

// Publish pushes message to particual queue
func (p *KafkaProducer) Publish(queue, message string) error {
	// Создавать publish'e ра при инициализации клиента. При асинхронном.
	// Graceful shudown - вссе сообщения возвращаются в очередь(не обработанные),
	if _, _, err := p.producer.SendMessage(&sarama.ProducerMessage{
		Topic:     queue,
		Value:     sarama.StringEncoder(message),
		Timestamp: time.Now(),
	}); err != nil {
		return err
	}

	return nil
}

// Close used in case of Graceful shutdown.
func (p *KafkaProducer) Close() error {
	err := p.producer.Close()
	if err != nil {
		return fmt.Errorf("error while closing publisher: %w", err)
	}
	return nil
}
