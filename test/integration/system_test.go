package main_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ccthomas/gridiron/internal/system"
	"github.com/stretchr/testify/assert"
)

func TestApiSystemHealth(t *testing.T) {
	// Given
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/system/health", nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	// When
	startTime := time.Now().UTC()
	res, err := http.DefaultClient.Do(req)
	endTime := time.Now().UTC()

	// Then
	if err != nil {
		log.Fatalln(err)
	}

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var healthResponse system.HealthMessage
	err = json.NewDecoder(res.Body).Decode(&healthResponse)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Gridiron Service is Healthy", healthResponse.Message)

	parsedTime, err := time.Parse(time.RFC3339, healthResponse.Timestamp)
	startTime = startTime.Truncate(time.Second)
	endTime = endTime.Truncate(time.Second)

	assert.NoError(t, err)
	assert.True(t, parsedTime.Equal(startTime) || parsedTime.After(startTime), "timestamp is not after or equal to the start of the test.")
	assert.True(t, parsedTime.Equal(endTime) || parsedTime.Before(endTime), "timestamp is not before or equal to the end of the test.")
}
