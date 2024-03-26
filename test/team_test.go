package test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ccthomas/gridiron/internal/team"
	"github.com/ccthomas/gridiron/internal/tenant"
	"github.com/ccthomas/gridiron/pkg/myhttp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewTeam(t *testing.T) {
	// Given

	u, loginRes := login(t)
	tn := createTenant(t, u.Id, "TestTenantName")

	unique := uuid.New().String()
	name := fmt.Sprintf("TestTeamName%s", unique)
	url := fmt.Sprintf("http://localhost:8080/teams/%s", name)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	req.Header.Add("Authorization", loginRes.AccessToken)
	req.Header.Add("x-tenant-id", tn.Id)

	// When

	res, err := http.DefaultClient.Do(req)

	// Then

	if err != nil {
		log.Fatalln(err)
	}

	assert.Equal(t, http.StatusOK, res.StatusCode, "Status code is not a 200")

	var actual team.Team
	err = json.NewDecoder(res.Body).Decode(&actual)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, actual.Id, "Actual id is not nil.")
	assert.Equal(t, actual.Name, name, "Actual name is incorrect.")
	assert.Equal(t, actual.TenantId, tn.Id, "Actual name is incorrect.")
}

func TestNewTeam_NoTenantId(t *testing.T) {
	// Given

	_, loginRes := login(t)

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/teams/name", nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	req.Header.Add("Authorization", loginRes.AccessToken)

	// When

	startTime := time.Now().UTC()
	res, err := http.DefaultClient.Do(req)
	endTime := time.Now().UTC()

	// Then

	if err != nil {
		log.Fatalln(err)
	}

	assert.Equal(t, http.StatusOK, res.StatusCode, "Status code is not a 200")

	var actual myhttp.ApiError
	err = json.NewDecoder(res.Body).Decode(&actual)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Tenant id was not provided.", actual.Message)
	parsedTime, err := time.Parse(time.RFC3339, actual.Timestamp)
	startTime = startTime.Truncate(time.Second)
	endTime = endTime.Truncate(time.Second)

	assert.NoError(t, err)
	assert.True(t, parsedTime.Equal(startTime) || parsedTime.After(startTime), "timestamp is not after or equal to the start of the test.")
	assert.True(t, parsedTime.Equal(endTime) || parsedTime.Before(endTime), "timestamp is not before or equal to the end of the test.")
}

func TestGetAllTeams_EmptyResponse(t *testing.T) {
	// Given

	_, loginRes := login(t)

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/teams", nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	req.Header.Add("Authorization", loginRes.AccessToken)

	// When

	res, err := http.DefaultClient.Do(req)

	// Then

	if err != nil {
		log.Fatalln(err)
	}

	assert.Equal(t, http.StatusOK, res.StatusCode, "Status code is not a 200")

	var actual tenant.TenantGetAllDTO
	err = json.NewDecoder(res.Body).Decode(&actual)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 0, actual.Count, "count does not equal 0.")
	assert.Equal(t, 0, len(actual.Data), "data count does not equal 0.")
}

func TestGetAllTeams_MultipleTeams(t *testing.T) {
	// Given

	u, loginRes := login(t)
	ten := createTenant(t, u.Id, fmt.Sprintf("TestTenant%s", u.Id))
	tm1 := createTeam(t, ten.Id, fmt.Sprintf("TestTeamA%s", ten.Id))
	tm2 := createTeam(t, ten.Id, fmt.Sprintf("TestTeamB%s", ten.Id))

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/teams", nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	req.Header.Add("Authorization", loginRes.AccessToken)

	// When

	res, err := http.DefaultClient.Do(req)

	// Then

	if err != nil {
		log.Fatalln(err)
	}

	assert.Equal(t, http.StatusOK, res.StatusCode, "Status code is not a 200")

	var actual tenant.TenantGetAllDTO
	err = json.NewDecoder(res.Body).Decode(&actual)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 2, actual.Count, "count does not equal 2.")
	assert.Equal(t, 2, len(actual.Data), "data count does not equal 2.")
	assert.Equal(t, tm1, actual.Data[0], "data does not equal first team.")
	assert.Equal(t, tm2, actual.Data[1], "data does not equal second team.")
}
