package test

import (
	"net/http"
	"testing"
	"time"

	"github.com/ccthomas/gridiron/pkg/myhttp"
	"github.com/stretchr/testify/assert"
)

// myhttp.ApiError has the same structure as system.HealthMessage.
// we will re-use the functions for api error to test health messages.

func TestSystemServiceHealth(t *testing.T) {

	// When

	startTime := time.Now().UTC()
	res, actual := sendApiReq[myhttp.ApiError](
		t,
		http.MethodGet,
		"http://localhost:8080/system/service/health",
		nil,
		"",
		"",
	)
	endTime := time.Now().UTC()

	// Then

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assertApiError(t, actual, "Gridiron Service is Healthy", startTime, endTime)
}

func TestSystemDatabaseHealth(t *testing.T) {

	// When

	startTime := time.Now().UTC()
	res, actual := sendApiReq[myhttp.ApiError](
		t,
		http.MethodGet,
		"http://localhost:8080/system/database/health",
		nil,
		"",
		"",
	)
	endTime := time.Now().UTC()

	// Then

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assertApiError(t, actual, "Gridiron has a healthy connection to the database.", startTime, endTime)
}
