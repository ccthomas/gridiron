package api

import (
	"github.com/ccthomas/gridiron/internal/system"
	"github.com/ccthomas/gridiron/internal/useracc"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// Handlers struct to hold dependencies for API handlers
type Handlers struct {
	Logger              *zap.Logger
	SystemHandlers      *system.SystemHandlers
	UserAccountHandlers *useracc.UserAccountHandlers
}

// NewHandlers initializes and returns a new Handlers instance
func NewHandlers(logger *zap.Logger, systemHandlers *system.SystemHandlers, userAccountHandlers *useracc.UserAccountHandlers) *Handlers {
	logger.Debug("Constructing new handlers")
	return &Handlers{
		Logger:              logger,
		SystemHandlers:      systemHandlers,
		UserAccountHandlers: userAccountHandlers,
	}
}

func (h *Handlers) RouteApis(r *mux.Router) {
	h.Logger.Debug("Route apis")
	h.routeSystemApis(r)
	h.routeUserAccountApis(r)
}

func (h *Handlers) routeSystemApis(r *mux.Router) {
	h.Logger.Debug("Configuring health handler routes")
	systemRoutes := r.PathPrefix("/system").Subrouter()

	h.Logger.Debug("Mapping api get /health to health handler function")
	systemRoutes.HandleFunc("/service/health", h.SystemHandlers.HealthHandler).Methods("GET")
	systemRoutes.HandleFunc("/database/health", h.SystemHandlers.DatabaseHealthHandler).Methods("GET")
}

func (h *Handlers) routeUserAccountApis(r *mux.Router) {
	h.Logger.Debug("Configuring health handler routes")
	systemRoutes := r.PathPrefix("/user").Subrouter()

	h.Logger.Debug("Mapping api get /health to health handler function")
	systemRoutes.HandleFunc("/authr", h.UserAccountHandlers.CreateNewUserHandler).Methods("GET")
}
