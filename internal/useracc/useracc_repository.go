package useracc

import (
	"database/sql"

	"github.com/ccthomas/gridiron/pkg/logger"
)

type UserAccountRepositoryImpl struct {
	DB *sql.DB
}

func NewUserAccountRepository(db *sql.DB) *UserAccountRepositoryImpl {
	logger.Get().Debug("Construct new user account repository.")
	return &UserAccountRepositoryImpl{
		DB: db,
	}
}

func (r *UserAccountRepositoryImpl) InsertUserAccount(userAccount UserAccount) error {
	logger.Get().Debug("Insert user account.")
	stmt, err := r.DB.Prepare("INSERT INTO user_account.user_account (id, username, password_hash) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	logger.Get().Debug("Execute the SQL statement with the user account data.")
	_, err = stmt.Exec(userAccount.Id, userAccount.Username, userAccount.PasswordHash)
	if err != nil {
		return err
	}

	logger.Get().Debug("Successfully inserted user account data.")
	return nil
}

func (r *UserAccountRepositoryImpl) SelectByUsername(username string) (*UserAccount, error) {
	logger.Get().Debug("Select user account by username.")
	stmt, err := r.DB.Prepare("SELECT id, username, password_hash FROM user_account.user_account WHERE username = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	logger.Get().Debug("Execute the SQL statement to retrieve the user account data.")
	row := stmt.QueryRow(username)

	logger.Get().Debug("Scan the data into a user account struct.")
	var userAccount UserAccount
	err = row.Scan(&userAccount.Id, &userAccount.Username, &userAccount.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Get().Debug("User account not found.")
			return nil, nil // Return nil if no user account is found
		}
		return nil, err
	}

	logger.Get().Debug("Found user account.")
	return &userAccount, nil
}
