package team

import "github.com/ccthomas/gridiron/pkg/rabbitmq"

// Data Transfer Objects

type CreateNewTeamDTO struct {
	Name string `json:"name"`
}

type TeamGetAllDTO struct {
	Count int    `json:"count"`
	Data  []Team `json:"data"`
}

// Entities

type Team struct {
	Id       string `json:"id"`
	TenantId string `json:"tenant_id"`
	Name     string `json:"name"`
}

// Interfaces

type TeamHandlers struct {
	RabbitMqRouter *rabbitmq.RabbitMqRouter
	TeamRepository TeamRepository
}

type TeamRepository interface {
	InsertTeam(team Team) error
	SelectAllTeamsByTenant(tenantId string) ([]Team, error)
}
