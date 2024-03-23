package api

import (
	"github.com/ccthomas/gridiron/internal/system"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// Handlers struct to hold dependencies for API handlers
type Handlers struct {
	Logger         *zap.Logger
	SystemHandlers *system.SystemHandlers
}

// NewHandlers initializes and returns a new Handlers instance
func NewHandlers(logger *zap.Logger, systemHandlers *system.SystemHandlers) *Handlers {
	logger.Debug("Constructing new handlers")
	return &Handlers{
		Logger:         logger,
		SystemHandlers: systemHandlers,
	}
}

func (h *Handlers) RouteApis(r *mux.Router) {
	h.Logger.Debug("Route apis")
	h.routeSystemApis(r)
}

func (h *Handlers) routeSystemApis(r *mux.Router) {
	h.Logger.Debug("Configuring health handler routes")
	systemRoutes := r.PathPrefix("/system").Subrouter()

	h.Logger.Debug("Mapping api get /health to health handler function")
	systemRoutes.HandleFunc("/service/health", h.SystemHandlers.HealthHandler).Methods("GET")
	systemRoutes.HandleFunc("/database/health", h.SystemHandlers.DatabaseHealthHandler).Methods("GET")
}
