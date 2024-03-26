package auth

type AccessLevel string

const (
	Owner AccessLevel = "OWNER"
)

type AuthorizerContext struct {
	UserId       string                 `json:"user_id"`
	TenantAccess map[string]AccessLevel `json:"tenant_access"`
}
