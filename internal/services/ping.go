package services

import (
	"context"

	"scm.x5.ru/x5m/go-backend/template/internal/config"
)

// PingService is a service for healthcheck
type PingService struct {
	cfg *config.Common
}

// NewPingService creates new PingService
func NewPingService(cfg *config.Common) *PingService {
	return &PingService{
		cfg: cfg,
	}
}

// Ping test
func (s *PingService) Ping(_ context.Context) error {
	return nil
}
