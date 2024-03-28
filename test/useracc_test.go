package test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ccthomas/gridiron/internal/useracc"
	"github.com/ccthomas/gridiron/pkg/auth"
	"github.com/ccthomas/gridiron/pkg/database"
	"github.com/ccthomas/gridiron/pkg/logger"
	"github.com/ccthomas/gridiron/pkg/myhttp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCreateNewUser(t *testing.T) {
	// Given
	db := database.ConnectPostgres()

	unique := uuid.New().String()

	userPass := &useracc.UserPassDTO{
		Username: fmt.Sprintf("TestCreateNewUser%s", unique),
		Password: fmt.Sprintf("TestCreateNewUser%s", unique),
	}

	// When

	res, actual := sendApiReq[useracc.UserAccount](
		t,
		http.MethodPost,
		"http://localhost:8080/user",
		userPass,
		"",
		"",
	)

	// Then

	assert.Equal(t, http.StatusOK, res.StatusCode, "Status code is not a 200")

	cleanUpUser(t, actual.Id)

	rows, err := db.Query("SELECT * FROM user_account.user_account WHERE id = $1", actual.Id)
	if err != nil {
		t.Fatal("Failed to prepare query.", err.Error())
	}

	defer rows.Close()
	rows.Next()

	var userAccount useracc.UserAccount
	err = rows.Scan(&userAccount.Id, &userAccount.Username, &userAccount.PasswordHash)
	if err != nil {
		t.Fatal("User was not created.", err.Error())
	}

	assert.Equal(t, userAccount.Id, actual.Id, "User account id is incorrect")
	assert.Equal(t, userAccount.Username, actual.Username, "User account username is incorrect")

	match := useracc.CheckPasswordHash(userPass.Password, userAccount.PasswordHash)
	assert.True(t, match, "User password does not match")
}

func TestCreateNewUser_UsernameTaken(t *testing.T) {
	// Given
	existing := createUser(t)

	unique := uuid.New().String()
	userPass := &useracc.UserPassDTO{
		Username: existing.Username,
		Password: fmt.Sprintf("TestCreateNewUser%s", unique),
	}

	// When

	startTime := time.Now().UTC()
	res, actual := sendApiReq[myhttp.ApiError](
		t,
		http.MethodPost,
		"http://localhost:8080/user",
		userPass,
		"",
		"",
	)
	endTime := time.Now().UTC()

	// Then

	assert.Equal(t, http.StatusBadRequest, res.StatusCode, "Status code is not a 400")
	assertApiError(t, actual, "Username is taken.", startTime, endTime)
}

func TestLogin_WithAuthorizerContext_TenantAccessEmpty(t *testing.T) {
	// Given - login

	existing, loginRes := login(t)

	// When - Authorizer Token

	res, actual := sendApiReq[auth.AuthorizerContext](
		t,
		http.MethodGet,
		"http://localhost:8080/user/authorizer-context",
		nil,
		loginRes.AccessToken,
		"",
	)

	// Then

	assert.Equal(t, http.StatusOK, res.StatusCode, "Status code is not a 200")
	assert.Equal(t, existing.Id, actual.UserId, "Authorizer context does not contain user id")
	assert.Equal(t, 0, len(actual.TenantAccess), "Authorizer context tenant access is not empty")
}

func TestLogin_WithAuthorizerContext_TenantAccessNonEmpty(t *testing.T) {
	// Given
	existing, loginRes := login(t)
	tenant1 := createTenant(t, existing.Id, "Name")

	// When - Authorizer Token

	res, actual := sendApiReq[auth.AuthorizerContext](
		t,
		http.MethodGet,
		"http://localhost:8080/user/authorizer-context",
		nil,
		loginRes.AccessToken,
		"",
	)

	// Then - Authorizer Token

	assert.Equal(t, http.StatusOK, res.StatusCode, "Status code is not a 200")
	assert.Equal(t, existing.Id, actual.UserId, "Authorizer context does not contain user id")
	assert.Equal(t, map[string]auth.AccessLevel{
		tenant1.Id: auth.Owner,
	}, actual.TenantAccess, "Authorizer context tenant access is not empty")
}

func TestLogin_WrongPassword(t *testing.T) {
	// Given - login
	existing := createUser(t)

	reqLogin, err := http.NewRequest(http.MethodPost, "http://localhost:8080/user/login", nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	reqLogin.SetBasicAuth(existing.Username, "Wrong Password")

	// When
	startTime := time.Now().UTC()
	res, err := http.DefaultClient.Do(reqLogin)
	endTime := time.Now().UTC()

	// Then
	if err != nil {
		log.Fatalln(err)
	}

	var actual myhttp.ApiError
	err = json.NewDecoder(res.Body).Decode(&actual)
	if err != nil {
		logger.Get().Fatal("Failed decoded response body for login with wrong password", zap.Error(err))
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusBadRequest, res.StatusCode, "Status code is not a 400")
	assertApiError(t, actual, "Invalid username or password.", startTime, endTime)
}

func TestLogin_UserDoesNotExist(t *testing.T) {
	reqLogin, err := http.NewRequest(http.MethodPost, "http://localhost:8080/user/login", nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	reqLogin.SetBasicAuth("Non existed user", "Wrong Password")

	// When
	startTime := time.Now().UTC()
	res, err := http.DefaultClient.Do(reqLogin)
	endTime := time.Now().UTC()

	// Then
	if err != nil {
		log.Fatalln(err)
	}

	var actual myhttp.ApiError
	err = json.NewDecoder(res.Body).Decode(&actual)
	if err != nil {
		logger.Get().Fatal("Failed decoded response body for login with wrong password", zap.Error(err))
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusBadRequest, res.StatusCode, "Status code is not a 400")
	assertApiError(t, actual, "Invalid username or password.", startTime, endTime)
}
