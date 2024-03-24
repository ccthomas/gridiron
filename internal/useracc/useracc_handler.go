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
func NewHandlers(userAccountRepository UserAccountRepository) *UserAccountHandlers {
	logger.Get().Debug("Constructing user account handlers")
	return &UserAccountHandlers{
		UserAccountRepository: userAccountRepository,
	}
}

func (h *UserAccountHandlers) CreateNewUserHandler(w http.ResponseWriter, r *http.Request) {
	logger.Get().Info("Create New User Handler hit.")

	var newUserData LoginData

	logger.Get().Debug("Decode new user data.")
	err := json.NewDecoder(r.Body).Decode(&newUserData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Get().Debug("Generate id for user.")
	id := uuid.New().String()

	logger.Get().Debug("Hash password.")
	hashedPassword, err := HashPassword(newUserData.Password)
	if err != nil {
		logger.Get().Error("Failed to hash password.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Get().Debug("Create a new user account struct.")
	userAccount := UserAccount{
		Id:           id,
		Username:     newUserData.Username,
		PasswordHash: hashedPassword,
	}

	logger.Get().Debug("Save user account.")
	err = h.UserAccountRepository.InsertUserAccount(userAccount)
	if err != nil {
		logger.Get().Error("Failed to insert user account into db.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Get().Debug("Construct response body.")
	response := &CreatedUser{
		Id:       id,
		Username: newUserData.Username,
	}

	logger.Get().Debug("Encode response JSON and write to response.")
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		logger.Get().Error("Failed to encode user account.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *UserAccountHandlers) GetAuthorizerContextHandler(w http.ResponseWriter, r *http.Request) {
	logger.Get().Info("Get Authorizer Context Handler hit.")

	ctx, err := auth.GetAuthorizerContext(r)
	if err != nil {
		logger.Logger.Debug("Failed to get authorizer context from request")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Get().Debug("Encode response JSON and write to response.")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(ctx)
	if err != nil {
		logger.Get().Error("Failed to encode authorizer context.")
		http.Error(w, "Failed to get context from token.", http.StatusInternalServerError)
		return
	}
}

func (h *UserAccountHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	logger.Get().Info("Login Handler hit.")
	username, password, ok := r.BasicAuth()
	if !ok {
		logger.Get().Warn("Failed to get basic auth.")
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		return
	}

	logger.Get().Debug("Find user by username.", zap.String("username", username))
	userAccount, err := h.UserAccountRepository.SelectByUsername(username)
	if err != nil {
		logger.Get().Warn("Username unknown.", zap.String("username", username))
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		return
	}

	logger.Get().Debug("Check password hash.")
	match := CheckPasswordHash(password, userAccount.PasswordHash)
	if !match {
		logger.Get().Warn("Password does not match.")
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		return
	}

	logger.Get().Debug("Generate new token.")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": userAccount.Id,
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})

	logger.Get().Debug("Fetch secret key.")
	secretKey := []byte(os.Getenv("SECRET_KEY"))

	logger.Get().Debug("Sign token.")
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		logger.Get().Error("Cannot sign token.", zap.Error(err))
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	logger.Get().Debug("Construct response byte array.")
	response := []byte(fmt.Sprintf(`{"accessToken":"%s"}`, tokenString))

	logger.Get().Debug("Set content-type header and write response.")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (h *UserAccountHandlers) TokenAuthorizerHandler(w http.ResponseWriter, r *http.Request) error {
	logger.Get().Info("Token Authorizer Handler hit.")

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
		return fmt.Errorf("authorization header is missing")
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	logger.Get().Debug("Fetch secret key.")
	secretKey := []byte(os.Getenv("SECRET_KEY"))

	logger.Get().Debug("Parse token.")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		logger.Get().Warn("Failed to parse token.")
		return err
	}

	if !token.Valid {
		logger.Get().Warn("Token was invalid.")
		return fmt.Errorf("token is not valid")
	}

	logger.Get().Debug("Get claims from token.", zap.Any("Token", token))
	claims := token.Claims.(jwt.MapClaims)
	id := claims["sub"].(string)

	logger.Get().Debug("JSON encode authorizer context")
	b, err := json.Marshal(auth.AuthorizerContext{
		UserId: id,
	})

	if err != nil {
		logger.Get().Warn("Failed to encode context.")
		return err
	}

	logger.Get().Debug("Set request context on request header.")
	r.Header.Set("request-context", string(b))
	return nil
}
