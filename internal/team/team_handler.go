package team

import (
	"encoding/json"
	"net/http"

	"github.com/ccthomas/gridiron/pkg/auth"
	"github.com/ccthomas/gridiron/pkg/logger"
	"github.com/ccthomas/gridiron/pkg/myhttp"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func NewHandlers(teamRepository TeamRepository) *TeamHandlers {
	logger.Get().Debug("Constructing tenant handlers")
	return &TeamHandlers{
		TeamRepository: teamRepository,
	}
}

func (h *TeamHandlers) CreateNewTeamHandler(w http.ResponseWriter, r *http.Request) {
	logger.Get().Info("Create New Team Handler hit.")

	ctx, err := auth.GetAuthorizerContext(r)
	if err != nil {
		logger.Logger.Debug("Failed to get authorizer context from request")
		myhttp.WriteError(w, http.StatusInternalServerError, "Internal Server Error.")
		return
	}

	logger.Get().Debug("Get tenant id.")
	tenantId := r.Header.Get("x-tenant-id")
	if ctx.TenantAccess[tenantId] != auth.Owner {
		logger.Logger.Debug("User does not have access to tenant.")
		myhttp.WriteError(w, http.StatusUnauthorized, "User is unauthorized to access tenant.")
		return
	}

	var dto CreateNewTeamDTO
	logger.Get().Debug("Decode new team data.")
	err = json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		logger.Get().Error("Failed to parse body.", zap.Error(err))
		myhttp.WriteError(w, http.StatusBadRequest, "payload provided was invalid.")
		return
	}

	t := Team{
		Id:       uuid.New().String(),
		TenantId: tenantId,
		Name:     dto.Name,
	}

	err = h.TeamRepository.InsertTeam(t)
	if err != nil {
		logger.Get().Error("Failed to insert team.", zap.Error(err))
		myhttp.WriteError(w, http.StatusInternalServerError, "Internal Server Error.")
		return
	}

	logger.Get().Debug("Encode response JSON and write to response.")
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(t)
	if err != nil {
		logger.Get().Error("Failed to encode team.")
		myhttp.WriteError(w, http.StatusInternalServerError, "Internal Server Error.")
		return
	}
}

func (h *TeamHandlers) GetAllTeamsHandler(w http.ResponseWriter, r *http.Request) {
	logger.Get().Info("Get All Teams Handler hit.")

	logger.Get().Debug("Get tenant id.")
	tenantId := r.Header.Get("x-tenant-id")

	teams, err := h.TeamRepository.SelectAllTeamsByTenant(tenantId)
	if err != nil {
		logger.Logger.Error("Failed to select teams by user.", zap.Error(err))
		myhttp.WriteError(w, http.StatusInternalServerError, "Internal Server Error.")
		return
	}

	var jsonResponse []byte = []byte(`{"count":0,"data":[]}`)
	if len(teams) != 0 {
		tenantGetAllDTO := &TeamGetAllDTO{
			Count: len(teams),
			Data:  teams,
		}

		// Convert games slice to JSON
		jsonResponse, err = json.Marshal(tenantGetAllDTO)
		if err != nil {
			logger.Logger.Error("Failed to marshal response.", zap.Error(err))
			myhttp.WriteError(w, http.StatusInternalServerError, "Internal Server Error.")
			return
		}
	}

	logger.Get().Debug("Encode response JSON and write to response.")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
