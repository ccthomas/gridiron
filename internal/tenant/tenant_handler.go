package tenant

import (
	"encoding/json"
	"net/http"

	"github.com/ccthomas/gridiron/pkg/auth"
	"github.com/ccthomas/gridiron/pkg/logger"
	"github.com/ccthomas/gridiron/pkg/myhttp"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func NewHandlers(tenantRepository TenantRepository) *TenantHandlers {
	logger.Get().Debug("Constructing tenant handlers")
	return &TenantHandlers{
		TenantRepository: tenantRepository,
	}
}

func (h *TenantHandlers) X(w http.ResponseWriter, r *http.Request) {
	logger.Get().Info("X Handler hit.")
}

func (h *TenantHandlers) GetAllTenantsHandler(w http.ResponseWriter, r *http.Request) {
	logger.Get().Info("New Tenant Handler hit.")
	ctx, err := auth.GetAuthorizerContext(r)
	if err != nil {
		logger.Logger.Debug("Failed to get authorizer context from request")
		myhttp.WriteError(w, http.StatusInternalServerError, "Internal Server Error.")
		return
	}

	tenants, err := h.TenantRepository.SelectTenantByUser(ctx.UserId)
	if err != nil {
		logger.Logger.Error("Failed to select tenants by user.", zap.Error(err))
		myhttp.WriteError(w, http.StatusInternalServerError, "Internal Server Error.")
		return
	}

	var jsonResponse []byte = []byte(`{"count":0,"data":[]}`)
	if len(tenants) != 0 {
		tenantGetAllDTO := &TenantGetAllDTO{
			Count: len(tenants),
			Data:  tenants,
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

func (h *TenantHandlers) NewTenantHandler(w http.ResponseWriter, r *http.Request) {
	logger.Get().Info("New Tenant Handler hit.")
	ctx, err := auth.GetAuthorizerContext(r)
	if err != nil {
		logger.Logger.Debug("Failed to get authorizer context from request")
		myhttp.WriteError(w, http.StatusInternalServerError, "Internal Server Error.")
		return
	}

	logger.Get().Debug("Get name parameter from path params.")
	params := mux.Vars(r)
	name := params["name"]

	logger.Get().Debug("Generate id for tenant.")
	id := uuid.New().String()

	tenant := Tenant{
		Id:   id,
		Name: name,
	}

	userAccess := TenantUserAccess{
		UserAccountId: ctx.UserId,
		TenantId:      tenant.Id,
		AccessLevel:   auth.Owner,
	}

	err = h.TenantRepository.InsertTenant(tenant)
	if err != nil {
		logger.Get().Error("Failed to insert tenant.", zap.Error(err))
		myhttp.WriteError(w, http.StatusInternalServerError, "Internal Server Error.")
		return
	}

	err = h.TenantRepository.InsertUserAccess(userAccess)
	if err != nil {
		logger.Get().Error("Failed to insert tenant user access.", zap.Error(err))
		myhttp.WriteError(w, http.StatusInternalServerError, "Internal Server Error.")
		return
	}

	logger.Get().Debug("Encode response JSON and write to response.")
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tenant)
	if err != nil {
		logger.Get().Error("Failed to encode user account.")
		myhttp.WriteError(w, http.StatusInternalServerError, "Internal Server Error.")
		return
	}
}
