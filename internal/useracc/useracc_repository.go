package useracc

import (
	"database/sql"

	"go.uber.org/zap"
)

type UserAccountRepositoryImpl struct {
	DB     *sql.DB
	Logger *zap.Logger
}

func NewUserAccountRepository(db *sql.DB, logger *zap.Logger) *UserAccountRepositoryImpl {
	logger.Debug("Construct new user account repository.")
	return &UserAccountRepositoryImpl{
		DB:     db,
		Logger: logger,
	}
}

func (r *UserAccountRepositoryImpl) InsertUserAccount(userAccount UserAccount) error {
	r.Logger.Debug("Insert user account.")
	stmt, err := r.DB.Prepare("INSERT INTO user_account.user_account (id, username, password_hash) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	r.Logger.Debug("Execute the SQL statement with the user account data.")
	_, err = stmt.Exec(userAccount.Id, userAccount.Username, userAccount.PasswordHash)
	if err != nil {
		return err
	}

	r.Logger.Debug("Successfully inserted user account data.")
	return nil
}

func (r *UserAccountRepositoryImpl) SelectByUsername(username string) (*UserAccount, error) {
	r.Logger.Debug("Select user account by username.")
	stmt, err := r.DB.Prepare("SELECT id, username, password_hash FROM user_account.user_account WHERE username = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	r.Logger.Debug("Execute the SQL statement to retrieve the user account data.")
	row := stmt.QueryRow(username)

	r.Logger.Debug("Scan the data into a user account struct.")
	var userAccount UserAccount
	err = row.Scan(&userAccount.Id, &userAccount.Username, &userAccount.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			r.Logger.Debug("User account not found.")
			return nil, nil // Return nil if no user account is found
		}
		return nil, err
	}

	r.Logger.Debug("Found user account.")
	return &userAccount, nil
}
