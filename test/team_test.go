package test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/ccthomas/gridiron/internal/team"
	"github.com/ccthomas/gridiron/internal/tenant"
	"github.com/ccthomas/gridiron/pkg/database"
	"github.com/ccthomas/gridiron/pkg/logger"
	"github.com/ccthomas/gridiron/pkg/myhttp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewTeam(t *testing.T) {
	// Given

	u, loginRes := login(t)
	tn := createTenant(t, u.Id, "TestTenantName")

	unique := uuid.New().String()
	dto := &team.CreateNewTeamDTO{
		Name: fmt.Sprintf("TestTeamName%s", unique),
	}

	// When

	res, actual := sendApiReq[team.Team](
		t,
		http.MethodPost,
		"http://localhost:8080/team",
		dto,
		loginRes.AccessToken,
		tn.Id,
	)

	// Then

	assert.Equal(t, http.StatusOK, res.StatusCode, "Status code is not a 200")

	cleanUpTeam(t, actual.Id)

	assert.NotNil(t, actual.Id, "Actual id is not nil.")
	assert.Equal(t, actual.Name, dto.Name, "Actual name is incorrect.")
	assert.Equal(t, actual.TenantId, tn.Id, "Actual name is incorrect.")
}

func TestNewTeam_NoTenantId(t *testing.T) {
	// Given

	_, loginRes := login(t)

	unique := uuid.New().String()
	dto := &team.CreateNewTeamDTO{
		Name: fmt.Sprintf("TestTeamName%s", unique),
	}

	// When

	startTime := time.Now().UTC()
	res, actual := sendApiReq[myhttp.ApiError](
		t,
		http.MethodPost,
		"http://localhost:8080/team",
		dto,
		loginRes.AccessToken,
		"",
	)
	endTime := time.Now().UTC()

	// Then

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode, "Status code is not a 401")
	assertApiError(t, actual, "User is unauthorized to access tenant.", startTime, endTime)
}

func TestNewTeam_DoesNotHaveAccessTenant(t *testing.T) {
	// Given

	_, loginRes := login(t)

	unique := uuid.New().String()
	dto := &team.CreateNewTeamDTO{
		Name: fmt.Sprintf("TestTeamName%s", unique),
	}

	// When

	startTime := time.Now().UTC()
	res, actual := sendApiReq[myhttp.ApiError](
		t,
		http.MethodPost,
		"http://localhost:8080/team",
		dto,
		loginRes.AccessToken,
		uuid.New().String(),
	)
	endTime := time.Now().UTC()

	// Then

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode, "Status code is not a 200")
	assertApiError(t, actual, "User is unauthorized to access tenant.", startTime, endTime)
}

func TestGetAllTeams_EmptyResponse(t *testing.T) {
	// Given

	_, loginRes := login(t)

	// When

	res, actual := sendApiReq[*team.TeamGetAllDTO](
		t,
		http.MethodGet,
		"http://localhost:8080/team",
		nil,
		loginRes.AccessToken,
		"",
	)

	// Then

	assert.Equal(t, http.StatusOK, res.StatusCode, "Status code is not a 200")
	assert.Equal(t, 0, actual.Count, "count does not equal 0.")
	assert.Equal(t, 0, len(actual.Data), "data count does not equal 0.")
}

func TestGetAllTeams_MultipleTeams(t *testing.T) {
	// Given
	u, loginRes := login(t)
	tn := createTenant(t, u.Id, fmt.Sprintf("TestTenant%s", u.Id))
	tm1 := createTeam(t, tn.Id, fmt.Sprintf("TestTeamA%s", tn.Id))
	tm2 := createTeam(t, tn.Id, fmt.Sprintf("TestTeamB%s", tn.Id))

	// When

	res, actual := sendApiReq[*team.TeamGetAllDTO](
		t,
		http.MethodGet,
		"http://localhost:8080/team",
		nil,
		loginRes.AccessToken,
		tn.Id,
	)

	// Then

	assert.Equal(t, http.StatusOK, res.StatusCode, "Status code is not a 200")
	assert.Equal(t, 2, actual.Count, "count does not equal 2.")
	assert.Equal(t, 2, len(actual.Data), "data count does not equal 2.")
	assert.Equal(t, tm1, actual.Data[0], "data does not equal first team.")
	assert.Equal(t, tm2, actual.Data[1], "data does not equal second team.")
}

func TestProcessNewTenantMessage(t *testing.T) {
	// Given

	db := database.ConnectPostgres()
	_, loginRes := login(t)
	unique := uuid.New().String()

	// When

	res, actual := sendApiReq[*tenant.Tenant](
		t,
		http.MethodPost,
		fmt.Sprintf("http://localhost:8080/tenant/%s", unique),
		nil,
		loginRes.AccessToken,
		"",
	)

	// Then

	assert.Equal(t, http.StatusOK, res.StatusCode, "Status code is not a 200")

	cleanUpTenant(t, actual.Id)
	time.Sleep(6 * time.Second)

	rows, err := db.Query("SELECT id, tenant_id, name FROM team.team WHERE tenant_id = $1", actual.Id)
	if err != nil {
		t.Fatal("Failed to prepare query.", err.Error())
	}

	var teams []team.Team
	for rows.Next() {
		var team team.Team
		logger.Get().Debug("Scan next row.")
		if err := rows.Scan(&team.Id, &team.TenantId, &team.Name); err != nil {
			t.Fatal("Failed to scan row.")

		}

		teams = append(teams, team)
	}

	// TODO This can be improved by verifying configurations match what's in the db.
	// Recommended solution would be to key each team by name in a map, then
	// loop through the configs (by name) and verify every config is correct.

	assert.Equal(t, 32, len(teams), "Teams is not of length 32.")
}
