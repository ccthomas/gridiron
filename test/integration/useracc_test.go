package integration

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/ccthomas/gridiron/internal/useracc"
	"github.com/ccthomas/gridiron/pkg/database"
	"github.com/ccthomas/gridiron/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewUser(t *testing.T) {
	// Given
	Logger := logger.Get()
	defer Logger.Sync()

	dm := database.DatabaseManager{
		Logger: Logger,
	}

	db := dm.ConnectPostgres()
	defer db.Close()

	// createUserData := &useracc.LoginData{
	// 	Username: "My User Name",
	// 	Password: "My Password",
	// }

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/user", nil)
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

	var userAccount useracc.UserAccount
	err = db.QueryRow("SELECT id, name FROM user_account.user_account WHERE id = $1", createdUser.Id).
		Scan(&userAccount.Id, &userAccount.PasswordHash, &userAccount.Username)

	if err != nil {
		t.Fatal("User was not created.")
	}

	assert.Equal(t, userAccount.Id, createdUser.Id)
	assert.Equal(t, userAccount.Username, createdUser.Username)

	// hashedPassword := userAccount.PasswordHash()
	assert.Equal(t, userAccount.Id, createdUser.Id)
}
