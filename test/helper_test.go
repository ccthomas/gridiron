package test

// Contains no tests.
// Failed named with _test.go to allow MainTest to execute.

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

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

func assertApiError(t *testing.T, body io.ReadCloser, message string, startTime time.Time, endTime time.Time) {
	var apiErr myhttp.ApiError
	err := json.NewDecoder(body).Decode(&apiErr)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, message, apiErr.Message)

	parsedTime, err := time.Parse(time.RFC3339, apiErr.Timestamp)
	startTime = startTime.Truncate(time.Second)
	endTime = endTime.Truncate(time.Second)

	assert.NoError(t, err)
	assert.True(t, parsedTime.Equal(startTime) || parsedTime.After(startTime), "timestamp is not after or equal to the start of the test.")
	assert.True(t, parsedTime.Equal(endTime) || parsedTime.Before(endTime), "timestamp is not before or equal to the end of the test.")

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

func cleanUpTenant(t *testing.T, id string) {
	t.Cleanup(func() {
		logger.Get().Debug("Clean up tenant.", zap.String("ID", id))
		db := database.ConnectPostgres()
		defer db.Close()

		_, err := db.Exec("DELETE FROM tenant.tenant_user_access WHERE tenant_id = $1", id)
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
