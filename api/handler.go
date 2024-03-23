package api

import (
	"database/sql"
	"net/http"

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
func NewHandlers(db *sql.DB, logger *zap.Logger, userRepo useracc.UserAccountRepository) *Handlers {
	logger.Debug("Constructing new handlers")

	systemHandlers := system.NewHandlers(db, logger)
	userAccHandlers := useracc.NewHandlers(logger, userRepo)

	return &Handlers{
		Logger:              logger,
		SystemHandlers:      systemHandlers,
		UserAccountHandlers: userAccHandlers,
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
	systemRoutes.HandleFunc("", h.UserAccountHandlers.CreateNewUserHandler).Methods("POST")
	systemRoutes.HandleFunc("/login", h.UserAccountHandlers.LoginHandler).Methods("POST")
	systemRoutes.HandleFunc("/authorizer-context", h.tokenAuthorizer(h.UserAccountHandlers.GetAuthorizerContextHandler)).Methods("GET")
}

func (h *Handlers) tokenAuthorizer(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Logger.Debug("Token Authorizer")

		// Extract the token from the Authorization header
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			h.Logger.Warn("Authorization header not provided.")
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		h.Logger.Debug("Authorize request.")
		err := h.UserAccountHandlers.TokenAuthorizerHandler(w, r)
		if err != nil {
			h.Logger.Warn("Is not authorizer")
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Pass the request to the next handler if the token is valid
		h.Logger.Debug("Is Authorized!")
		next.ServeHTTP(w, r)
	}
}
