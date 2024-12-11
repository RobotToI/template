package events

import (
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
	"scm.x5.ru/x5m/go-backend/template/internal/repository/kafka"
)

var _ sarama.ConsumerGroupHandler = (*ConsumerGroupHandler)(nil)

// ConsumerGroupHandler store channel to accept readiness
type ConsumerGroupHandler struct {
	ready chan bool
}

// Event представляет cобытие любого типа
type Event struct {
	EventID   int64     `json:"event_id"`
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
}

// NewHandler создает ConsumerGroupHandler создавая небуфиризированный канал
func NewHandler() *ConsumerGroupHandler {
	return &ConsumerGroupHandler{
		ready: make(chan bool),
	}
}

// Setup Начинаем новую сессию, до ConsumeClaim
func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup завершает сессию, после того, как все ConsumeClaim завершатся
func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim читаем до тех пор пока сессия не завершилась
func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			msg := kafka.ConvertMsg(message)

			var event Event
			err := json.Unmarshal([]byte(msg.Payload), &event)
			if err != nil {
				log.Println("failed to unmarshall event", msg, err)
				continue
			}

			log.Printf("received event. event id: %d, event type: %s, timestamp: %s", event.EventID, event.EventType, event.Timestamp.Format("02.01.2006 15:04:05"))

			// mark message as successfully handled and ready to commit offset
			// autocommit may commit message offset sometime
			session.MarkMessage(message, "")

			// commit offset manually right now
			// works when autocommit disabled
			session.Commit()

			// autocommit not to work if cg not gracefully shut downed
			// but manual commit does
			// panic("emulate not commit")
		case <-session.Context().Done():
			return nil
		}
	}
}
