package test

import (
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/ccthomas/gridiron/pkg/myhttp"
	"github.com/stretchr/testify/assert"
)

func assertApiError(t *testing.T, body io.ReadCloser, message string, startTime time.Time, endTime time.Time) {
	var apiErr myhttp.ApiError
	err := json.NewDecoder(body).Decode(&apiErr)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, message, apiErr.Message)

	parsedTime, err := time.Parse(time.RFC3339, apiErr.Timestamp)
	startTime = startTime.Truncate(time.Second)
	endTime = endTime.Truncate(time.Second)

	assert.NoError(t, err)
	assert.True(t, parsedTime.Equal(startTime) || parsedTime.After(startTime), "timestamp is not after or equal to the start of the test.")
	assert.True(t, parsedTime.Equal(endTime) || parsedTime.Before(endTime), "timestamp is not before or equal to the end of the test.")

}
