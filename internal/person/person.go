package person

type PersonContractEntityType string
type PersonContractType string

const (
	Team PersonContractEntityType = "TEAM"

	Athlete PersonContractType = "ATHLETE"
	Coach   PersonContractType = "COACH"
	Owner   PersonContractType = "OWNER"
)

// Data Transfer Object

type CreateNewPersonDTO struct {
	Name string `json:"name"`
}

type CreateNewPersonContractDTO struct {
	PersonId   string                   `json:"person_id"`
	EntityId   string                   `json:"entity_id"`
	EntityType PersonContractEntityType `json:"entity_type"`
	Type       PersonContractType       `json:"type"`
}

type PersonWithContractDTO struct {
	Person   Person         `json:"person"`
	Contract PersonContract `json:"contract"`
}

type PersonGetDTO struct {
	Count int                     `json:"count"`
	Data  []PersonWithContractDTO `json:"data"`
}

// Entity

type Person struct {
	Id       string `json:"id"`
	TenantId string `json:"tenant_id"`
	Name     string `json:"name"`
}

type PersonContract struct {
	Id         string                   `json:"id"`
	TenantId   string                   `json:"tenant_id"`
	PersonId   string                   `json:"person_id"`
	EntityId   string                   `json:"entity_id"`
	EntityType PersonContractEntityType `json:"entity_type"`
	Type       PersonContractType       `json:"type"`
}
