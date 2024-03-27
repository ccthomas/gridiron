package tenant

import (
	"github.com/ccthomas/gridiron/pkg/auth"
	"github.com/ccthomas/gridiron/pkg/rabbitmq"
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
	TenantId      string           `json:"tenant_id"`
	UserAccountId string           `json:"user_account_id"`
	AccessLevel   auth.AccessLevel `json:"access_level"`
}

// Interfaces

type TenantHandlers struct {
	RabbitMqRouter   *rabbitmq.RabbitMqRouter
	TenantRepository TenantRepository
}

type TenantRepository interface {
	InsertTenant(tenant Tenant) error
	InsertUserAccess(userAccess TenantUserAccess) error
	SelectTenantByUser(userId string) ([]Tenant, error)
	SelectTenantAccessByUser(userId string) ([]TenantUserAccess, error)
}
