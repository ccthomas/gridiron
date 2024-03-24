package integration

import (
	"bytes"
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
)

// func createUser(t *testing.T, db *sql.DB) (useracc.UserAccount) {
// 	id := uuid.New().String()
// 	password := fmt.Sprintf("password%s", id)
// 	passwordHash, _ := useracc.HashPassword(password)
// 	user := useracc.UserAccount{
// 		Id:           id,
// 		Username:     fmt.Sprintf("My Username %s", id),
// 		PasswordHash: passwordHash,
// 	}
// 	db.Exec("INSERT INTO user_account.user_account VALUES (id, username, password_hash) VALUES ($1, $2, $3)", user.Id, user.Username, user.PasswordHash)

// 	t.Cleanup(func() {
// 		db.Exec("DELETE FROM user_account.user_account WHERE id = $1")
// 	})

// 	return user
// }

func TestMain(m *testing.M) {
	// Load .env file
	if err := godotenv.Load("../../.env"); err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	// Run tests
	exitVal := m.Run()

	// Exit with the same exit code as the tests
	os.Exit(exitVal)
}

func TestCreateNewUser(t *testing.T) {
	// Given
	Logger := logger.Get()
	defer Logger.Sync()

	db := database.ConnectPostgres()
	defer db.Close()

	unique := uuid.New().String()

	createUserData := &useracc.LoginData{
		Username: fmt.Sprintf("TestCreateNewUser%s", unique),
		Password: fmt.Sprintf("TestCreateNewUser%s", unique),
	}

	jsonData, err := json.Marshal(createUserData)
	if err != nil {
		fmt.Printf("could not marshal createUserData: %s\n", err)
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

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var createdUser useracc.CreatedUser
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

	stmt, err := db.Prepare("SELECT id, username, password_hash FROM user_account.user_account WHERE id = $1")
	if err != nil {
		t.Fatal("Failed to prepare query.", err.Error())
	}
	defer stmt.Close()

	// Execute the SQL statement to retrieve the user account data
	row := stmt.QueryRow(createdUser.Id)

	// Scan the data into a UserAccount struct
	var userAccount useracc.UserAccount
	err = row.Scan(&userAccount.Id, &userAccount.Username, &userAccount.PasswordHash)
	if err != nil {
		t.Fatal("User was not created.")
	}

	assert.Equal(t, userAccount.Id, createdUser.Id)
	assert.Equal(t, userAccount.Username, createdUser.Username)

	match := useracc.CheckPasswordHash(createUserData.Password, userAccount.PasswordHash)
	assert.True(t, match)
}
