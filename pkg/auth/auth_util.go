package auth

import (
	"encoding/json"
	"net/http"

	"github.com/ccthomas/gridiron/pkg/logger"
)

func GetAuthorizerContext(r *http.Request) (AuthorizerContext, error) {
	logger.Get().Debug("Get Authorizer Context.")
	var authorizerContext AuthorizerContext
	err := json.Unmarshal([]byte(r.Header.Get("Request-Context")), &authorizerContext)
	logger.Get().Debug("Returning Authorizer Context.")
	return authorizerContext, err
}
