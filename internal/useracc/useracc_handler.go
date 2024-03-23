package useracc

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/ccthomas/gridiron/pkg/auth"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

// NewHandlers initializes and returns a new Handlers instance
func NewHandlers(db *sql.DB, logger *zap.Logger) *UserAccountHandlers {
	logger.Debug("Constructing user account handlers")
	return &UserAccountHandlers{
		DB:     db,
		Logger: logger,
	}
}

func (h *UserAccountHandlers) CreateNewUserHandler(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Create New User Handler hit.")
	h.Logger.Info("Token", zap.String("Key", os.Getenv("PRIVATE_KEY")))
}

func (h *UserAccountHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Login Handler hit.")
}

func (h *UserAccountHandlers) TokenAuthorizerHandler(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("Token Authorizer Handler hit.")

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
		return fmt.Errorf("authorization header is missing")
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	h.Logger.Debug("Parse token.")
	secretKey := []byte(os.Getenv("PRIVATE_KEY"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		h.Logger.Warn("Failed to parse token.")
		return err
	}

	if !token.Valid {
		h.Logger.Warn("Token was invalid.")
		return fmt.Errorf("token is not valid")
	}

	h.Logger.Debug("Get claims from token.")

	claims := token.Claims.(jwt.MapClaims)
	id := claims["id"].(string)

	h.Logger.Debug("JSON encode authorizer context")
	b, err := json.Marshal(auth.AuthorizerContext{
		UserId: id,
	})

	if err != nil {
		h.Logger.Warn("Failed to encode context.")
		return err
	}

	h.Logger.Debug("Set request context on request header.")
	r.Header.Set("request-context", string(b))
	return nil
}
