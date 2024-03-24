package useracc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ccthomas/gridiron/pkg/auth"
	"github.com/ccthomas/gridiron/pkg/logger"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// NewHandlers initializes and returns a new Handlers instance
func NewHandlers(logger *zap.Logger, userAccountRepository UserAccountRepository) *UserAccountHandlers {
	logger.Debug("Constructing user account handlers")
	return &UserAccountHandlers{
		Logger:                logger,
		UserAccountRepository: userAccountRepository,
	}
}

func (h *UserAccountHandlers) CreateNewUserHandler(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Create New User Handler hit.")

	var newUserData LoginData

	h.Logger.Debug("Decode new user data.")
	err := json.NewDecoder(r.Body).Decode(&newUserData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.Logger.Debug("Generate id for user.")
	id := uuid.New().String()

	h.Logger.Debug("Hash password.")
	hashedPassword, err := HashPassword(newUserData.Password)
	if err != nil {
		h.Logger.Error("Failed to hash password.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Debug("Create a new user account struct.")
	userAccount := UserAccount{
		Id:           id,
		Username:     newUserData.Username,
		PasswordHash: hashedPassword,
	}

	h.Logger.Debug("Save user account.")
	err = h.UserAccountRepository.InsertUserAccount(userAccount)
	if err != nil {
		h.Logger.Error("Failed to insert user account into db.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Debug("Construct response body.")
	response := &CreatedUser{
		Id:       id,
		Username: newUserData.Username,
	}

	h.Logger.Debug("Encode response JSON and write to response.")
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.Logger.Error("Failed to encode user account.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *UserAccountHandlers) GetAuthorizerContextHandler(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Get Authorizer Context Handler hit.")

	ctx, err := auth.GetAuthorizerContext(r)
	if err != nil {
		logger.Logger.Debug("Failed to get authorizer context from request")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Debug("Encode response JSON and write to response.")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(ctx)
	if err != nil {
		h.Logger.Error("Failed to encode authorizer context.")
		http.Error(w, "Failed to get context from token.", http.StatusInternalServerError)
		return
	}
}

func (h *UserAccountHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Login Handler hit.")
	username, password, ok := r.BasicAuth()
	if !ok {
		h.Logger.Warn("Failed to get basic auth.")
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		return
	}

	h.Logger.Debug("Find user by username.", zap.String("username", username))
	userAccount, err := h.UserAccountRepository.SelectByUsername(username)
	if err != nil {
		h.Logger.Warn("Username unknown.", zap.String("username", username))
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		return
	}

	h.Logger.Debug("Check password hash.")
	match := CheckPasswordHash(password, userAccount.PasswordHash)
	if !match {
		h.Logger.Warn("Password does not match.")
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		return
	}

	h.Logger.Debug("Generate new token.")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": userAccount.Id,
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})

	h.Logger.Debug("Fetch secret key.")
	secretKey := []byte(os.Getenv("PRIVATE_KEY"))

	h.Logger.Debug("Sign token.")
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		h.Logger.Error("Cannot sign token.", zap.Error(err))
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	h.Logger.Debug("Construct response byte array.")
	response := []byte(fmt.Sprintf(`{"accessToken":"%s"}`, tokenString))

	h.Logger.Debug("Set content-type header and write response.")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (h *UserAccountHandlers) TokenAuthorizerHandler(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("Token Authorizer Handler hit.")

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
		return fmt.Errorf("authorization header is missing")
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	h.Logger.Debug("Fetch secret key.")
	secretKey := []byte(os.Getenv("PRIVATE_KEY"))

	h.Logger.Debug("Parse token.")
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

	h.Logger.Debug("Get claims from token.", zap.Any("Token", token))
	claims := token.Claims.(jwt.MapClaims)
	id := claims["sub"].(string)

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
