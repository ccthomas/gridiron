package system

import (
	"database/sql"
)

type HealthMessage struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

type SystemHandlers struct {
	DB *sql.DB
}
