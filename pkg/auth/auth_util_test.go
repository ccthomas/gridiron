package auth

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestGetAuthorizerContext(t *testing.T) {
	// Given
	testUserID := "test_user_id"
	authorizerContextJSON := `{"user_id":"` + testUserID + `"}`

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	req.Header.Set("Request-Context", authorizerContextJSON)

	// When
	authorizerContext, err := GetAuthorizerContext(req)

	// Then
	if err != nil {
		t.Fatalf("GetAuthorizerContext returned an error: %v", err)
	}

	if authorizerContext.UserId != testUserID {
		t.Errorf("expected user ID %s, got %s", testUserID, authorizerContext.UserId)
	}
}

func TestGetAuthorizerContext_Error(t *testing.T) {
	// Given
	invalidAuthorizerContextJSON := `{invalid JSON}`

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	req.Header.Set("Request-Context", invalidAuthorizerContextJSON)

	// When
	_, err = GetAuthorizerContext(req)

	// Then
	if err == nil {
		t.Error("expected an error but got none")
	}

	expectedErrorMessage := "invalid character 'i' looking for beginning of object key string"
	if !strings.Contains(err.Error(), expectedErrorMessage) {
		t.Errorf("expected error message containing %q, got %q", expectedErrorMessage, err.Error())
	}

	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected a JSON syntax error, got %T", err)
	}
}
