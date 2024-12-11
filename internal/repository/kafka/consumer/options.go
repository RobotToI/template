package consumer

import (
	"github.com/IBM/sarama"
)

// Option is a configuration callback
type Option interface {
	Apply(*sarama.Config) error
}

type optionFn func(*sarama.Config) error

func (fn optionFn) Apply(c *sarama.Config) error {
	return fn(c)
}

// WithOffsetsInitial определяет с какого смещения начинать вычитывать сообщения
func WithOffsetsInitial(v int64) Option {
	return optionFn(func(c *sarama.Config) error {
		c.Consumer.Offsets.Initial = v
		return nil
	})
}
