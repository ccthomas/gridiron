package api

import (
	"database/sql"
	"net/http"

	"github.com/ccthomas/gridiron/internal/system"
	"github.com/ccthomas/gridiron/internal/team"
	"github.com/ccthomas/gridiron/internal/tenant"
	"github.com/ccthomas/gridiron/internal/useracc"
	"github.com/ccthomas/gridiron/pkg/logger"
	"github.com/ccthomas/gridiron/pkg/myhttp"
	"github.com/ccthomas/gridiron/pkg/rabbitmq"
	"github.com/gorilla/mux"
)

type Handlers struct {
	SystemHandlers      *system.SystemHandlers
	TeamHandlers        *team.TeamHandlers
	TenantHandlers      *tenant.TenantHandlers
	UserAccountHandlers *useracc.UserAccountHandlers
}

func NewHandlers(
	db *sql.DB,
	rmq *rabbitmq.RabbitMqRouter,
	teamRepo team.TeamRepository,
	tenantRepo tenant.TenantRepository,
	userRepo useracc.UserAccountRepository,
) *Handlers {
	logger.Get().Debug("Constructing new handlers")

	systemHandlers := system.NewHandlers(db)
	teamHandlers := team.NewHandlers(rmq, teamRepo)
	tenantHandlers := tenant.NewHandlers(rmq, tenantRepo)
	userAccHandlers := useracc.NewHandlers(tenantRepo, userRepo)

	return &Handlers{
		SystemHandlers:      systemHandlers,
		TeamHandlers:        teamHandlers,
		TenantHandlers:      tenantHandlers,
		UserAccountHandlers: userAccHandlers,
	}
}

func (h *Handlers) RouteApis(r *mux.Router, rmq *rabbitmq.RabbitMqRouter) {
	logger.Get().Debug("Route apis")
	h.routeSystemApis(r)
	h.routeTeamApis(r, rmq)
	h.routeTenantApis(r)
	h.routeUserAccountApis(r)
}

func (h *Handlers) routeSystemApis(r *mux.Router) {
	logger.Get().Debug("Configuring system handler routes")
	systemRoutes := r.PathPrefix("/system").Subrouter()

	systemRoutes.HandleFunc("/service/health", h.SystemHandlers.HealthHandler).Methods("GET")
	systemRoutes.HandleFunc("/database/health", h.SystemHandlers.DatabaseHealthHandler).Methods("GET")
}

func (h *Handlers) routeTeamApis(r *mux.Router, rmq *rabbitmq.RabbitMqRouter) {
	logger.Get().Debug("Configuring team handler routes")
	tenantRoutes := r.PathPrefix("/team").Subrouter()

	tenantRoutes.HandleFunc("", h.tokenAuthorizer(h.TeamHandlers.CreateNewTeamHandler)).Methods("POST")
	tenantRoutes.HandleFunc("", h.tokenAuthorizer(h.TeamHandlers.GetAllTeamsHandler)).Methods("GET")

	rmq.HandleFunc("tenant-exchange", "New Tenant", h.TeamHandlers.ProcessNewTenantMessageHandler)
}

func (h *Handlers) routeTenantApis(r *mux.Router) {
	logger.Get().Debug("Configuring tenant handler routes")
	tenantRoutes := r.PathPrefix("/tenant").Subrouter()

	tenantRoutes.HandleFunc("", h.tokenAuthorizer(h.TenantHandlers.GetAllTenantsHandler)).Methods("GET")
	tenantRoutes.HandleFunc("/{name}", h.tokenAuthorizer(h.TenantHandlers.NewTenantHandler)).Methods("POST")
}

func (h *Handlers) routeUserAccountApis(r *mux.Router) {
	logger.Get().Debug("Configuring user account handler routes")
	systemRoutes := r.PathPrefix("/user").Subrouter()

	systemRoutes.HandleFunc("", h.UserAccountHandlers.CreateNewUserHandler).Methods("POST")
	systemRoutes.HandleFunc("/login", h.UserAccountHandlers.LoginHandler).Methods("POST")
	systemRoutes.HandleFunc("/authorizer-context", h.tokenAuthorizer(h.UserAccountHandlers.GetAuthorizerContextHandler)).Methods("GET")
}

func (h *Handlers) tokenAuthorizer(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Get().Debug("Token Authorizer")

		// Extract the token from the Authorization header
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			logger.Get().Warn("Authorization header not provided.")
			myhttp.WriteError(w, http.StatusUnauthorized, "Authorization header is missing.")
			return
		}

		logger.Get().Debug("Authorize request.")
		err := h.UserAccountHandlers.TokenAuthorizerHandler(w, r)
		if err != nil {
			logger.Get().Warn("Is not authorized.")
			return
		}

		// Pass the request to the next handler if the token is valid
		logger.Get().Debug("Is Authorized!")
		next.ServeHTTP(w, r)
	}
}
