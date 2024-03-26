package tenant

import (
	"database/sql"

	"github.com/ccthomas/gridiron/pkg/logger"
)

type TenantRepositoryImpl struct {
	DB *sql.DB
}

func NewTenantRepository(db *sql.DB) *TenantRepositoryImpl {
	logger.Get().Debug("Construct new tenant repository.")
	return &TenantRepositoryImpl{
		DB: db,
	}
}

// InsertTenant inserts a new tenant into the database.
func (r *TenantRepositoryImpl) InsertTenant(tenant Tenant) error {
	logger.Get().Debug("Insert tenant.")
	_, err := r.DB.Exec("INSERT INTO tenant.tenant (id, name) VALUES ($1, $2)", tenant.Id, tenant.Name)
	if err != nil {
		logger.Get().Warn("Failed to insert tenant.")
		return err
	}

	logger.Get().Debug("Successfully inserted tenant.")
	return nil
}

func (r *TenantRepositoryImpl) InsertUserAccess(userAccess TenantUserAccess) error {
	logger.Get().Debug("Insert user access.")
	_, err := r.DB.Exec("INSERT INTO tenant.tenant_user_access (user_account_id, tenant_id, access_level) VALUES ($1, $2, $3)", userAccess.UserAccountId, userAccess.TenantId, userAccess.AccessLevel)

	if err != nil {
		logger.Get().Warn("Failed to insert user access.")
		return err
	}

	logger.Get().Debug("Successfully inserted user access.")
	return nil
}

func (r *TenantRepositoryImpl) SelectTenantByUser(userId string) ([]Tenant, error) {
	logger.Get().Debug("Select tenants by user id.")
	rows, err := r.DB.Query("SELECT t.id, t.name FROM tenant.tenant t JOIN tenant.tenant_user_access ua ON t.id = ua.tenant_id WHERE ua.user_account_id = $1 ORDER BY t.name ASC", userId)

	if err != nil {
		logger.Get().Warn("Failed to select tenant by user.")
		return nil, err
	}

	defer rows.Close()

	logger.Get().Debug("Start scanning rows.")
	var tenants []Tenant
	for rows.Next() {
		var tenant Tenant
		logger.Get().Debug("Scan next row.")
		if err := rows.Scan(&tenant.Id, &tenant.Name); err != nil {
			logger.Get().Warn("Failed to scan row.")
			return nil, err
		}

		logger.Get().Debug("Add tenant to tenant array.")
		tenants = append(tenants, tenant)
	}

	logger.Get().Debug("Return tenants.")
	return tenants, nil
}
