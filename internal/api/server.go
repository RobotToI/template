package api

import (
	"context"
	"net/http"

	"scm.x5.ru/x5m/go-backend/template/internal/config"

	R "github.com/go-pkgz/rest"
	"scm.x5.ru/x5m/go-backend/packages/zlogger"
)

//go:generate oapi-codegen --config=cfg.yml ../../swagger/swagger.yml

// PingService is a mockable interface for ping service
//
//go:generate moq -rm -out ping_mock.go . PingService
type PingService interface {
	Ping(ctx context.Context) error
}

// Controller is a handler for Controller API
type Controller struct {
	pingSrv PingService
	cfg     *config.Common
}

var _ ServerInterface = (*Controller)(nil)

// ControllerDeps is a dependency for Controller
type ControllerDeps struct {
	PingSrv PingService
}

// NewController constructs a new Controller
//
// Parameters:
//   - cfg: Configuration for the controller
//   - deps: Dependencies for the controller
//
// Returns:
//   - A new instance of Controller
func NewController(cfg *config.Common, deps ControllerDeps) *Controller {
	return &Controller{
		pingSrv: deps.PingSrv,
		cfg:     cfg,
	}
}

// Ping handles the HTTP request to check the service health
//
// Parameters:
//   - w: HTTP response writer (not used)
//   - r: HTTP request
func (q *Controller) Ping(w http.ResponseWriter, r *http.Request) {
	logger := zlogger.LoadOrCreateFromCtx(r.Context())
	logger.Info().Msg("Ping")

	R.RenderJSON(w, "success")
}
