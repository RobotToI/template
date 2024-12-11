package config

import (
	"fmt"
	"time"
)

// Common configuration
type Common struct {
	Environment string `envconfig:"APP_ENV" required:"true"` // Окружение приложения
	Name        string `envconfig:"APP_NAME" default:"template" required:"false"`
	Server      HTTPServer
	Debug       Debug
	PostgreSQL  PostgreSQL
	Redis       Redis
	Kafka       Kafka
}

// Redis configuration
type Redis struct {
	Host      string `envconfig:"APP_REDIS_HOST" required:"true"`
	Port      int    `envconfig:"APP_REDIS_PORT" required:"true"`
	Password  string `envconfig:"APP_REDIS_PASSWORD" required:"false"`
	Database  int    `envconfig:"APP_REDIS_DATABASE" required:"false"`
	KeyPrefix string `envconfig:"APP_REDIS_KEY_PREFIX" required:"false"`
}

// GetAddr in format `redis.host:redis.port`
func (r *Redis) GetAddr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

// Debug configuration
type Debug struct {
	Enabled bool `envconfig:"APP_DEBUG" required:"false" default:"false"` // Режим отладки
}

// PostgreSQL configuration settings
type PostgreSQL struct {
	Host           string        `envconfig:"DB_HOST" required:"true" default:"127.0.0.1"`
	Port           int           `envconfig:"DB_PORT" required:"true" default:"5432"`
	Username       string        `envconfig:"DB_USERNAME" required:"true" default:"username"`
	Password       string        `envconfig:"DB_PASSWORD" required:"true" default:"password"`
	Database       string        `envconfig:"DB_DATABASE" required:"true" default:"database"`
	Ssl            string        `envconfig:"DB_SSL_MODE" required:"true" default:"disable"`
	ConnectTimeout time.Duration `envconfig:"DB_CONNECT_TIMEOUT" required:"true" default:"10s"`
	QueryTimeout   time.Duration `envconfig:"DB_QUERY_TIMEOUT" required:"true" default:"10s"`
}

// BuildDSN builds connection url for Postgres
func (p *PostgreSQL) BuildDSN() string {
	// TODO: add param -> statement_cache_mode=describe
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s", p.Username, p.Password, p.Host, p.Port, p.Database, p.Ssl)
}

// Kafka configuration settings
type Kafka struct {
	Brokers  []string `envconfig:"APP_KAFKA_HOSTS" required:"true" default:"localhost:9092"`
	GroupID  string   `envconfig:"APP_KAFKA_GROUPID" required:"true" default:"template"`
	Topics   []string `envconfig:"APP_KAFKA_TOPICS" required:"true" default:"template_example"`
	Assignor string   `envconfig:"APP_KAFKA_ASSIGNOR" required:"true" default:"roundrobin"`
	Version  string   `envconfig:"APP_KAFKA_VERSION" required:"true" default:"2.5.0"`
	Oldest   bool     `envconfig:"APP_KAFKA_OLDEST" required:"false" default:"true"`
	SASL     struct {
		Enable    bool   `envconfig:"APP_KAFKA_SASL_ENABLE"`
		Handshake bool   `envconfig:"APP_KAFKA_SASL_HANDSHAKE"`
		User      string `envconfig:"APP_KAFKA_SASL_USER"`
		Password  string `envconfig:"APP_KAFKA_SASL_PASSWORD"`
	}
}

// HTTPServer configuration settings
type HTTPServer struct {
	Port int `envconfig:"HTTP_LISTEN_PORT" required:"true"`
}

// GetListenPort in format `:PORT`
func (s *HTTPServer) GetListenPort() string {
	return fmt.Sprintf(":%d", s.Port)
}
