package test

import (
	"bytes"
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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewUser(t *testing.T) {
	// Given
	db := database.ConnectPostgres()

	unique := uuid.New().String()

	userPass := &useracc.UserPassDTO{
		Username: fmt.Sprintf("TestCreateNewUser%s", unique),
		Password: fmt.Sprintf("TestCreateNewUser%s", unique),
	}

	jsonData, err := json.Marshal(userPass)
	if err != nil {
		fmt.Printf("could not marshal userPass: %s\n", err)
		os.Exit(1)
	}

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/user", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	// When
	res, err := http.DefaultClient.Do(req)

	// Then
	if err != nil {
		log.Fatalln(err)
	}

	assert.Equal(t, http.StatusOK, res.StatusCode, "Status code is not a 200")

	var createdUser useracc.CreatedUserDTO
	err = json.NewDecoder(res.Body).Decode(&createdUser)
	if err != nil {
		t.Fatal(err)
	}

	cleanUpUser(t, createdUser.Id)

	rows, err := db.Query("SELECT * FROM user_account.user_account WHERE id = $1", createdUser.Id)
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

	assert.Equal(t, userAccount.Id, createdUser.Id, "User account id is incorrect")
	assert.Equal(t, userAccount.Username, createdUser.Username, "User account username is incorrect")

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

	jsonData, err := json.Marshal(userPass)
	if err != nil {
		fmt.Printf("could not marshal userPass: %s\n", err)
		os.Exit(1)
	}

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/user", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	// When
	startTime := time.Now().UTC()
	res, err := http.DefaultClient.Do(req)
	endTime := time.Now().UTC()

	// Then
	if err != nil {
		log.Fatalln(err)
	}

	assert.Equal(t, http.StatusBadRequest, res.StatusCode, "Status code is not a 400")
	assertApiError(t, res.Body, "Username is taken.", startTime, endTime)
}

func TestLogin_WithAuthorizerContext(t *testing.T) {
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

	assert.NotNil(t, loginResponse.AccessToken, "Access token is nil.")

	// Given - Authorizer Token
	reqAuth, err := http.NewRequest(http.MethodGet, "http://localhost:8080/user/authorizer-context", nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	reqAuth.Header.Set("Authorization", loginResponse.AccessToken)

	// When - Authorizer Token
	resAuth, err := http.DefaultClient.Do(reqAuth)

	// Then - Authorizer Token
	if err != nil {
		log.Fatalln(err)
	}

	var authCtx auth.AuthorizerContext
	err = json.NewDecoder(resAuth.Body).Decode(&authCtx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, existing.Id, authCtx.UserId, "Authorizer context does not contain user id")
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

	assert.Equal(t, http.StatusBadRequest, res.StatusCode, "Status code is not a 400")
	assertApiError(t, res.Body, "Invalid username or password.", startTime, endTime)
}

func TestLogin_UserDoesNotExist(t *testing.T) {
	// Given - login
	reqLogin, err := http.NewRequest(http.MethodPost, "http://localhost:8080/user/login", nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	reqLogin.SetBasicAuth("Non Existent User", "Wrong Password")

	// When
	startTime := time.Now().UTC()
	res, err := http.DefaultClient.Do(reqLogin)
	endTime := time.Now().UTC()

	// Then
	if err != nil {
		log.Fatalln(err.Error())
	}

	assert.Equal(t, http.StatusBadRequest, res.StatusCode, "Status code is not a 400")
	assertApiError(t, res.Body, "Invalid username or password.", startTime, endTime)
}
