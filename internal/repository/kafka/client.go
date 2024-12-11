package kafka

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"scm.x5.ru/x5m/go-backend/template/internal/config"
)

// Client is container for storing created KafkaConnection & Group
type Client struct {
	Client sarama.Client
	Group  string
}

// GetClient returns KafkaClient.Client field
func (c *Client) GetClient() (sarama.Client, error) {
	return c.Client, nil
}

// GetGroup returns KafkaClient.Group field
func (c *Client) GetGroup() string {
	return c.Group
}

// Close should be used in Graceful Shutdown of the service
func (c *Client) Close() error {
	return c.Client.Close()
}

// NewClient creates and configures KafkaClient
func NewClient(_ context.Context, cfg config.Kafka) (*Client, error) {
	version, err := sarama.ParseKafkaVersion(cfg.Version)
	if err != nil {
		return nil, fmt.Errorf("kafka version: %v", err)
	}

	conf := sarama.NewConfig()
	conf.Version = version

	switch cfg.Assignor {
	case "sticky":
		conf.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategySticky()
	case "roundrobin":
		conf.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	case "range":
		conf.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	default:
		return nil, fmt.Errorf("unknown group partition assignor: %s", cfg.Assignor)
	}

	if cfg.Oldest {
		conf.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	if cfg.SASL.Enable {
		conf.Net.SASL.Enable = cfg.SASL.Enable
		conf.Net.SASL.Handshake = cfg.SASL.Handshake
		conf.Net.SASL.User = cfg.SASL.User
		conf.Net.SASL.Password = cfg.SASL.Password
	}

	conf.Consumer.Return.Errors = true
	conf.Producer.Return.Successes = true
	conf.Producer.RequiredAcks = sarama.WaitForLocal
	conf.Producer.Partitioner = sarama.NewRoundRobinPartitioner

	client, err := sarama.NewClient(cfg.Brokers, conf)
	if err != nil {
		return nil, err
	}

	return &Client{
		Client: client,
		Group:  cfg.GroupID,
	}, nil
}
