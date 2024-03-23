package useracc

import (
	"go.uber.org/zap"
)

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserAccount struct {
	Id           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}

type UserAccountRepository interface {
	InsertUserAccount(userAccount UserAccount) error
	SelectByUsername(username string) (*UserAccount, error)
}

type UserAccountHandlers struct {
	Logger                *zap.Logger
	UserAccountRepository UserAccountRepository
}
