package test

// Contains no tests.
// Failed named with _test.go to allow MainTest to execute.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ccthomas/gridiron/internal/person"
	"github.com/ccthomas/gridiron/internal/team"
	"github.com/ccthomas/gridiron/internal/tenant"
	"github.com/ccthomas/gridiron/internal/useracc"
	"github.com/ccthomas/gridiron/pkg/database"
	"github.com/ccthomas/gridiron/pkg/logger"
	"github.com/ccthomas/gridiron/pkg/myhttp"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type UserAccountWithPass struct {
	Id           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Password     string `json:"password"`
}

func TestMain(m *testing.M) {
	// Load .env file
	if err := godotenv.Load("../.env.offline"); err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	Logger := logger.Get()
	defer Logger.Sync()
	Logger.Info("Starting tests.")

	db := database.ConnectPostgres()
	defer db.Close()

	// Run tests
	exitVal := m.Run()

	// Exit with the same exit code as the tests
	os.Exit(exitVal)
}

func sendApiReq[K any](
	t *testing.T,
	method string,
	url string,
	body any,
	auth string,
	tenantId string,
) (*http.Response, K) {
	var req *http.Request
	var err error

	if body == nil {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			t.Fatal("Failed to construct request", err.Error())
		}
	} else {
		jsonData, err := json.Marshal(body)
		if err != nil {
			fmt.Printf("could not marshal userPass: %s\n", err)
			os.Exit(1)
		}

		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatal("Failed to construct request", err.Error())
		}
	}

	if auth != "" {
		req.Header.Set("Authorization", auth)
	}
	if auth != "" {
		req.Header.Set("x-tenant-id", tenantId)
	}

	// req.Close = true

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal("Apu request failed.", err.Error())
	}

	// defer res.Body.Close()

	var actual K
	err = json.NewDecoder(res.Body).Decode(&actual)
	if err != nil {
		logger.Get().Fatal("Failed decoded body", zap.String("url", url), zap.Error(err))
		t.Fatal(err)
	}

	return res, actual
}

func assertApiError(t *testing.T, body myhttp.ApiError, message string, startTime time.Time, endTime time.Time) {
	assert.Equal(t, message, body.Message)

	parsedTime, err := time.Parse(time.RFC3339, body.Timestamp)
	startTime = startTime.Truncate(time.Second)
	endTime = endTime.Truncate(time.Second)

	assert.NoError(t, err)
	assert.True(t, parsedTime.Equal(startTime) || parsedTime.After(startTime), "timestamp is not after or equal to the start of the test.")
	assert.True(t, parsedTime.Equal(endTime) || parsedTime.Before(endTime), "timestamp is not before or equal to the end of the test.")
}

func cleanUpTenant(t *testing.T, id string) {
	t.Cleanup(func() {
		// provide time for for async
		// rabbitmq messages to be processed.
		time.Sleep(2 * time.Second)

		logger.Get().Debug("Clean up tenant.", zap.String("ID", id))
		db := database.ConnectPostgres()
		defer db.Close()

		_, err := db.Exec("DELETE FROM person.person_contract WHERE tenant_id = $1", id)
		if err != nil {
			logger.Get().Error("Failed to clean up person.")
			t.Fatal(err.Error())
		}

		_, err = db.Exec("DELETE FROM person.person WHERE tenant_id = $1", id)
		if err != nil {
			logger.Get().Error("Failed to clean up person.")
			t.Fatal(err.Error())
		}

		_, err = db.Exec("DELETE FROM team.team WHERE tenant_id = $1", id)
		if err != nil {
			logger.Get().Error("Failed to clean up team.")
			t.Fatal(err.Error())
		}

		_, err = db.Exec("DELETE FROM tenant.tenant_user_access WHERE tenant_id = $1", id)
		if err != nil {
			logger.Get().Error("Failed to clean up tenant.")
			t.Fatal(err.Error())
		}

		_, err = db.Exec("DELETE FROM tenant.tenant WHERE id = $1", id)
		if err != nil {
			logger.Get().Error("Failed to clean up user.")
			t.Fatal(err.Error())
		}

	})
}

func cleanUpUser(t *testing.T, id string) {
	t.Cleanup(func() {
		logger.Get().Debug("Clean up user.", zap.String("ID", id))
		db := database.ConnectPostgres()
		defer db.Close()
		_, err := db.Exec("DELETE FROM user_account.user_account WHERE id = $1", id)
		if err != nil {
			logger.Get().Error("Failed to clean up user.")
			t.Fatal(err.Error())
		}
	})
}

func createPerson(t *testing.T, tenantId string, personName string) person.Person {
	db := database.ConnectPostgres()
	defer db.Close()

	p := person.Person{
		Id:       uuid.New().String(),
		TenantId: tenantId,
		Name:     personName,
	}

	_, err := db.Exec("INSERT INTO person.person (id, tenant_id, name) VALUES ($1, $2, $3)", p.Id, p.TenantId, p.Name)
	if err != nil {
		t.Fatal("Failed to insert team as a part of setup.", err.Error())
	}

	cleanUpTenant(t, tenantId)
	return p
}

func createPersonContract(t *testing.T, tenantId string, dto person.CreateNewPersonContractDTO) person.PersonContract {
	db := database.ConnectPostgres()
	defer db.Close()

	pc := person.PersonContract{
		Id:         uuid.New().String(),
		TenantId:   tenantId,
		PersonId:   dto.PersonId,
		EntityId:   dto.EntityId,
		EntityType: dto.EntityType,
		Type:       dto.Type,
	}

	_, err := db.Exec(
		"INSERT INTO person.person_contract (id, tenant_id, person_id, entity_id, entity_type, type) VALUES ($1, $2, $3, $4, $5, $6)",
		pc.Id, pc.TenantId, pc.PersonId, pc.EntityId, pc.EntityType, pc.Type,
	)

	if err != nil {
		t.Fatal("Failed to insert team as a part of setup.", err.Error())
	}

	cleanUpTenant(t, tenantId)
	return pc
}

func createTeam(t *testing.T, tenantId string, teamName string) team.Team {
	db := database.ConnectPostgres()
	defer db.Close()

	myTeam := team.Team{
		Id:       uuid.New().String(),
		TenantId: tenantId,
		Name:     teamName,
	}

	_, err := db.Exec("INSERT INTO team.team (id, tenant_id, name) VALUES ($1, $2, $3)", myTeam.Id, myTeam.TenantId, myTeam.Name)
	if err != nil {
		t.Fatal("Failed to insert team as a part of setup.", err.Error())
	}

	cleanUpTenant(t, tenantId)
	return myTeam
}

func createTenant(t *testing.T, userId string, tenantName string) tenant.Tenant {
	db := database.ConnectPostgres()

	tenant := tenant.Tenant{
		Id:   uuid.New().String(),
		Name: tenantName,
	}

	_, err := db.Exec(
		"INSERT INTO tenant.tenant (id, name) VALUES ($1, $2)",
		tenant.Id, tenant.Name,
	)
	if err != nil {
		t.Fatal("Failed to insert tenants as a part of setup.", err.Error())
	}

	cleanUpTenant(t, tenant.Id)

	_, err = db.Exec(
		"INSERT INTO tenant.tenant_user_access (tenant_id, user_account_id, access_level) VALUES ($1, $2, $3)",
		tenant.Id, userId, "OWNER",
	)
	if err != nil {
		t.Fatal("Failed to insert tenant user access as a part of setup.", err.Error())
	}

	return tenant
}

func createUser(t *testing.T) UserAccountWithPass {
	db := database.ConnectPostgres()
	defer db.Close()

	id := uuid.New().String()
	password := fmt.Sprintf("password%s", id)
	passwordHash, _ := useracc.HashPassword(password)
	user := UserAccountWithPass{
		Id:           id,
		Username:     fmt.Sprintf("TestCreateNewUser%s", id),
		PasswordHash: passwordHash,
		Password:     password,
	}
	_, err := db.Exec("INSERT INTO user_account.user_account (id, username, password_hash) VALUES ($1, $2, $3)", user.Id, user.Username, user.PasswordHash)
	if err != nil {
		t.Fatal("Failed to insert user as a part of setup.", err.Error())
	}

	cleanUpUser(t, user.Id)
	return user
}

func login(t *testing.T) (UserAccountWithPass, useracc.LoginResponseDTO) {
	// Given - login
	existing := createUser(t)

	reqLogin, err := http.NewRequest(http.MethodPost, "http://localhost:8080/user/login", nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	reqLogin.SetBasicAuth(existing.Username, existing.Password)

	// When - login
	resLogin, err := http.DefaultClient.Do(reqLogin)

	// Then - login
	if err != nil {
		log.Fatalln(err)
	}

	assert.Equal(t, http.StatusOK, resLogin.StatusCode, "Status code is not a 200")

	var loginResponse useracc.LoginResponseDTO
	err = json.NewDecoder(resLogin.Body).Decode(&loginResponse)
	if err != nil {
		t.Fatal(err)
	}

	return existing, loginResponse
}
