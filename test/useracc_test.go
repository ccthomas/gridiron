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
	"github.com/ccthomas/gridiron/pkg/logger"
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

func cleanUp(t *testing.T, id string) {
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

	cleanUp(t, user.Id)
	return user
}

func TestMain(m *testing.M) {
	// Load .env file
	if err := godotenv.Load("../.env.offline"); err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	Logger := logger.Get()
	defer Logger.Sync()

	db := database.ConnectPostgres()
	defer db.Close()

	// Run tests
	exitVal := m.Run()

	// Exit with the same exit code as the tests
	os.Exit(exitVal)
}

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

	t.Cleanup(func() {
		db := database.ConnectPostgres()
		defer db.Close()
		_, err := db.Exec("DELETE FROM user_account.user_account WHERE id = $1", createdUser.Id)
		if err != nil {
			logger.Get().Error("Failed to clean up user.")
			t.Fatal(err.Error())
		}
	})

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
