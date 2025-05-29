package health

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthHandler_ServeHTTP(t *testing.T) {
	tests := []struct {
		name    string
		env     string
		appName string
		version string
	}{
		{
			name:    "standard case",
			env:     "local",
			appName: "uala-twitter",
			version: "1.0.0",
		},
		{
			name:    "production",
			env:     "live",
			appName: "uala-twitter",
			version: "2.5.1",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewHealthHandler(tc.env, tc.appName, tc.version)
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)

			assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

			var body ServiceInfo
			err := json.Unmarshal(rr.Body.Bytes(), &body)
			assert.NoError(t, err)
			assert.Equal(t, tc.env, body.Env)
			assert.Equal(t, tc.appName, body.Name)
			assert.Equal(t, tc.version, body.Version)
		})
	}
}
