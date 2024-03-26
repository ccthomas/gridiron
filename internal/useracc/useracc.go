package useracc

import "github.com/ccthomas/gridiron/internal/tenant"

// Data Transfer Objects

type CreatedUserDTO struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type UserPassDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponseDTO struct {
	AccessToken string `json:"access_token"`
}

// Entities

type UserAccount struct {
	Id           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}

// Interfaces

type UserAccountHandlers struct {
	TenantRepository      tenant.TenantRepository
	UserAccountRepository UserAccountRepository
}

type UserAccountRepository interface {
	InsertUserAccount(userAccount UserAccount) error
	SelectByUsername(username string) (*UserAccount, error)
}
