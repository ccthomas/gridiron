package useracc

// Data Transfer Objects

type CreatedUserDTO struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type UserPassDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Entities

type UserAccount struct {
	Id           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}

// Interfaces

type UserAccountHandlers struct {
	UserAccountRepository UserAccountRepository
}

type UserAccountRepository interface {
	InsertUserAccount(userAccount UserAccount) error
	SelectByUsername(username string) (*UserAccount, error)
}
