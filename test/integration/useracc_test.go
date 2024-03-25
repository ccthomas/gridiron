package integration

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/ccthomas/gridiron/internal/useracc"
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
		db := database.ConnectPostgres()
		defer db.Close()
		_, err := db.Exec("DELETE FROM user_account.user_account WHERE id = $1", id)
		if err != nil {
			logger.Get().Error("Failed to clean up user.")
			t.Fatal(err.Error())
		}
	})
}

func createUser(t *testing.T, db *sql.DB) UserAccountWithPass {
	id := uuid.New().String()
	password := fmt.Sprintf("password%s", id)
	passwordHash, _ := useracc.HashPassword(password)
	user := UserAccountWithPass{
		Id:           id,
		Username:     fmt.Sprintf("TestCreateNewUser%s", id),
		PasswordHash: passwordHash,
		Password:     password,
	}
	db.Exec("INSERT INTO user_account.user_account VALUES (id, username, password_hash) VALUES ($1, $2, $3)", user.Id, user.Username, user.PasswordHash)

	cleanUp(t, user.Id)
	return user
}

func TestMain(m *testing.M) {
	// Load .env file
	if err := godotenv.Load("../../.env.offline"); err != nil {
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

	logger.Logger.Debug("Created User", zap.Any("Id", createdUser.Id))

	// SELECT * FROM user_account.user_account WHERE id = '9debbfd9-73d0-4e92-ae19-3ded8eb86ff4';
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
	db := database.ConnectPostgres()
	existing := createUser(t, db)

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
	res, err := http.DefaultClient.Do(req)

	// Then
	if err != nil {
		log.Fatalln(err)
	}

	assert.Equal(t, http.StatusBadRequest, res.StatusCode, "Status code is not a 400")

	var createdUser useracc.CreatedUserDTO
	err = json.NewDecoder(res.Body).Decode(&createdUser)
	if err != nil {
		t.Fatal(err)
	}

	logger.Logger.Debug("Created User", zap.Any("Id", createdUser.Id))

	// SELECT * FROM user_account.user_account WHERE id = '9debbfd9-73d0-4e92-ae19-3ded8eb86ff4';
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
