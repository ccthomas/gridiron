package team

// Data Transfer Objects

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
