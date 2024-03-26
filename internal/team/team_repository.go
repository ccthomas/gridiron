package team

import (
	"database/sql"

	"github.com/ccthomas/gridiron/pkg/logger"
)

type TeamRepositoryImpl struct {
	DB *sql.DB
}

func NewTenantRepository(db *sql.DB) *TeamRepositoryImpl {
	logger.Get().Debug("Construct new team repository.")
	return &TeamRepositoryImpl{
		DB: db,
	}
}

func (r *TeamRepositoryImpl) InsertTeam(team Team) error {
	logger.Get().Debug("Insert team.")
	_, err := r.DB.Exec("INSERT INTO team.team (id, tenant_id, name) VALUES ($1, $2, $3)", team.Id, team.TenantId, team.Name)
	if err != nil {
		logger.Get().Warn("Failed to insert team.")
		return err
	}

	logger.Get().Debug("Successfully inserted team.")
	return nil
}

func (r *TeamRepositoryImpl) SelectAllTeamsByTenant(tenantId string) ([]Team, error) {
	logger.Get().Debug("Select team by tenant id.")
	rows, err := r.DB.Query("SELECT id, tenant_id, name FROM team.team WHERE tenant_id = $1 ORDER BY name ASC;", tenantId)

	if err != nil {
		logger.Get().Warn("Failed to select team by tenant id.")
		return nil, err
	}

	defer rows.Close()

	logger.Get().Debug("Start scanning rows.")
	var teams []Team
	for rows.Next() {
		var t Team
		logger.Get().Debug("Scan next row.")
		if err := rows.Scan(&t.Id, &t.TenantId, &t.Name); err != nil {
			logger.Get().Warn("Failed to scan row.")
			return nil, err
		}

		logger.Get().Debug("Add team to teams array.")
		teams = append(teams, t)
	}

	logger.Get().Debug("Return teams.")
	return teams, nil
}
