package test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/ccthomas/gridiron/internal/tenant"
	"github.com/ccthomas/gridiron/pkg/auth"
	"github.com/ccthomas/gridiron/pkg/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewTenant(t *testing.T) {
	// Given
	db := database.ConnectPostgres()
	existing, loginResp := login(t)

	unique := uuid.New().String()

	url := fmt.Sprintf("http://localhost:8080/tenant/%s", unique)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	req.Header.Add("Authorization", loginResp.AccessToken)

	// When
	res, err := http.DefaultClient.Do(req)

	// Then
	if err != nil {
		log.Fatalln(err)
	}

	assert.Equal(t, http.StatusOK, res.StatusCode, "Status code is not a 200")

	var actual tenant.Tenant
	err = json.NewDecoder(res.Body).Decode(&actual)
	if err != nil {
		t.Fatal(err)
	}

	cleanUpTenant(t, actual.Id)

	rows, err := db.Query("SELECT * FROM tenant.tenant WHERE id = $1", actual.Id)
	if err != nil {
		t.Fatal("Failed to prepare query.", err.Error())
	}

	defer rows.Close()
	rows.Next()

	var expected tenant.Tenant
	err = rows.Scan(&expected.Id, &expected.Name)
	if err != nil {
		t.Fatal("User was not created.", err.Error())
	}

	assert.Equal(t, expected.Id, actual.Id, "tenant id id is incorrect")
	assert.Equal(t, expected.Name, actual.Name, "tenant name is incorrect")

	rows, err = db.Query("SELECT * FROM tenant.tenant_user_access WHERE tenant_id = $1", actual.Id)
	if err != nil {
		t.Fatal("Failed to prepare query.", err.Error())
	}

	defer rows.Close()
	rows.Next()

	var userAccess tenant.TenantUserAccess
	err = rows.Scan(&userAccess.TenantId, &userAccess.UserAccountId, &userAccess.AccessLevel)
	if err != nil {
		t.Fatal("User was not created.", err.Error())
	}

	assert.Equal(t, userAccess.TenantId, actual.Id, "User access tenant id is incorrect")
	assert.Equal(t, userAccess.UserAccountId, existing.Id, "User access user account id is incorrect")
	assert.Equal(t, userAccess.AccessLevel, auth.Owner, "User access access level is incorrect")
}

func TestNewTenant_ConfigurationNFL(t *testing.T) {
	// Given
	db := database.ConnectPostgres()
	existing, loginResp := login(t)

	unique := uuid.New().String()

	url := fmt.Sprintf("http://localhost:8080/tenant/%s", unique)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	req.Header.Add("Authorization", loginResp.AccessToken)

	// When
	res, err := http.DefaultClient.Do(req)

	// Then
	if err != nil {
		log.Fatalln(err)
	}

	assert.Equal(t, http.StatusOK, res.StatusCode, "Status code is not a 200")

	var actual tenant.Tenant
	err = json.NewDecoder(res.Body).Decode(&actual)
	if err != nil {
		t.Fatal(err)
	}

	cleanUpTenant(t, actual.Id)

	rows, err := db.Query("SELECT * FROM tenant.tenant WHERE id = $1", actual.Id)
	if err != nil {
		t.Fatal("Failed to prepare query.", err.Error())
	}

	defer rows.Close()
	rows.Next()

	var expected tenant.Tenant
	err = rows.Scan(&expected.Id, &expected.Name)
	if err != nil {
		t.Fatal("User was not created.", err.Error())
	}

	assert.Equal(t, expected.Id, actual.Id, "tenant id id is incorrect")
	assert.Equal(t, expected.Name, actual.Name, "tenant name is incorrect")

	rows, err = db.Query("SELECT * FROM tenant.tenant_user_access WHERE tenant_id = $1", actual.Id)
	if err != nil {
		t.Fatal("Failed to prepare query.", err.Error())
	}

	defer rows.Close()
	rows.Next()

	var userAccess tenant.TenantUserAccess
	err = rows.Scan(&userAccess.TenantId, &userAccess.UserAccountId, &userAccess.AccessLevel)
	if err != nil {
		t.Fatal("User was not created.", err.Error())
	}

	assert.Equal(t, userAccess.TenantId, actual.Id, "User access tenant id is incorrect")
	assert.Equal(t, userAccess.UserAccountId, existing.Id, "User access user account id is incorrect")
	assert.Equal(t, userAccess.AccessLevel, auth.Owner, "User access access level is incorrect")
}

func TestGetAllTenants(t *testing.T) {
	// Given
	existing, loginResp := login(t)
	existingOther, _ := login(t)

	tenant1 := createTenant(t, existing.Id, "TestingA")
	tenant2 := createTenant(t, existing.Id, "TestingB")
	createTenant(t, existingOther.Id, "TestingFilteredOut") // tenant to be filtered out

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/tenant", nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	req.Header.Add("Authorization", loginResp.AccessToken)

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

	assert.Equal(t, actual.Count, 2, "count does not equal 2.")
	assert.Equal(t, len(actual.Data), 2, "data count does not equal 2.")
	assert.Equal(t, actual.Data[0], tenant1, "data does not equal first tenant.")
	assert.Equal(t, actual.Data[1], tenant2, "data does not equal second tenant.")
}
