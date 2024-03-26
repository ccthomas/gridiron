package tenant

type AccessLevel string

const (
	Owner AccessLevel = "OWNER"
)

// Data Transfer Objects

type TenantGetAllDTO struct {
	Count int      `json:"count"`
	Data  []Tenant `json:"data"`
}

// Entities

type Tenant struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type TenantUserAccess struct {
	TenantId      string      `json:"tenant_id"`
	UserAccountId string      `json:"user_account_id"`
	AccessLevel   AccessLevel `json:"access_level"`
}

// Interfaces

type TenantHandlers struct {
	TenantRepository TenantRepository
}

type TenantRepository interface {
	InsertTenant(tenant Tenant) error
	InsertUserAccess(userAccess TenantUserAccess) error
	SelectTenantByUser(userId string) ([]Tenant, error)
}
