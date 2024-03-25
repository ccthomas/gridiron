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
	"github.com/ccthomas/gridiron/pkg/myhttp"
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

	var userPass UserPassDTO

	logger.Get().Debug("Decode new user data.")
	err := json.NewDecoder(r.Body).Decode(&userPass)
	if err != nil {
		myhttp.WriteError(w, http.StatusBadRequest, "payload provided was invalid.")
		return
	}

	logger.Get().Debug("Generate id for user.")
	id := uuid.New().String()

	logger.Get().Debug("Hash password.")
	hashedPassword, err := HashPassword(userPass.Password)
	if err != nil {
		logger.Get().Error("Failed to hash password.")
		myhttp.WriteError(w, http.StatusInternalServerError, "Internal Server Error.")
		return
	}

	logger.Get().Debug("Create a new user account struct.")
	userAccount := UserAccount{
		Id:           id,
		Username:     userPass.Username,
		PasswordHash: hashedPassword,
	}

	logger.Get().Debug("Save user account.")
	err = h.UserAccountRepository.InsertUserAccount(userAccount)
	if err != nil {
		logger.Get().Error("Failed to insert user account into db.", zap.Error(err))

		if err.Error() == "pq: duplicate key value violates unique constraint \"user_account_username_key\"" {
			myhttp.WriteError(w, http.StatusBadRequest, "Username is taken.")
			return
		}

		myhttp.WriteError(w, http.StatusInternalServerError, "Internal Server Error.")
		return
	}

	logger.Get().Debug("Construct response body.")
	response := &CreatedUserDTO{
		Id:       id,
		Username: userPass.Username,
	}

	logger.Get().Debug("Encode response JSON and write to response.")
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		logger.Get().Error("Failed to encode user account.")
		myhttp.WriteError(w, http.StatusInternalServerError, "Internal Server Error.")
		return
	}
}

func (h *UserAccountHandlers) GetAuthorizerContextHandler(w http.ResponseWriter, r *http.Request) {
	logger.Get().Info("Get Authorizer Context Handler hit.")

	ctx, err := auth.GetAuthorizerContext(r)
	if err != nil {
		logger.Logger.Debug("Failed to get authorizer context from request")
		myhttp.WriteError(w, http.StatusInternalServerError, "Internal Server Error.")
		return
	}

	logger.Get().Debug("Encode response JSON and write to response.")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(ctx)
	if err != nil {
		logger.Get().Error("Failed to encode authorizer context.")
		myhttp.WriteError(w, http.StatusInternalServerError, "Internal Server Error.")
		return
	}
}

func (h *UserAccountHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	logger.Get().Info("Login Handler hit.")
	username, password, ok := r.BasicAuth()
	if !ok {
		logger.Get().Warn("Failed to get basic auth.")
		myhttp.WriteError(w, http.StatusBadRequest, "Invalid username or password.")
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
		myhttp.WriteError(w, http.StatusInternalServerError, "Internal Server Error.")
		return
	}

	logger.Get().Debug("Construct response byte array.")
	response := []byte(fmt.Sprintf(`{"access_token":"%s"}`, tokenString))

	logger.Get().Debug("Set content-type header and write response.")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (h *UserAccountHandlers) TokenAuthorizerHandler(w http.ResponseWriter, r *http.Request) error {
	logger.Get().Info("Token Authorizer Handler hit.")

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		myhttp.WriteError(w, http.StatusUnauthorized, "Authorization header is missing.")
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
		myhttp.WriteError(w, http.StatusInternalServerError, "Internal Server Error.")
		return err
	}

	if !token.Valid {
		logger.Get().Warn("Token was invalid.")
		myhttp.WriteError(w, http.StatusUnauthorized, "Authorization header is missing.")
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
		myhttp.WriteError(w, http.StatusInternalServerError, "Internal Server Error.")
		return err
	}

	logger.Get().Debug("Set request context on request header.")
	r.Header.Set("request-context", string(b))
	return nil
}
