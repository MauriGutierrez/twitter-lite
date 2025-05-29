package usecase

import (
	"testing"
)

func TestUseCaseErrors(t *testing.T) {
	tests := []struct {
		name           string
		createFunc     func(string, ...error) *UseCaseError
		expectedType   string
		expectedMsg    string
		expectedCauses []error
	}{
		{
			name:         "InternalServerError",
			createFunc:   InternalServerError,
			expectedType: TypeInternalServerError,
			expectedMsg:  "Internal server error",
		},
		{
			name:         "NotFound",
			createFunc:   NotFound,
			expectedType: TypeNotFound,
			expectedMsg:  "Resource not found",
		},
		{
			name:         "Unknown",
			createFunc:   Unknown,
			expectedType: TypeUnknown,
			expectedMsg:  "Unknown error",
		},
		{
			name:         "Forbidden",
			createFunc:   Forbidden,
			expectedType: TypeForbidden,
			expectedMsg:  "Access forbidden",
		},
		{
			name:         "InvalidParam",
			createFunc:   InvalidParam,
			expectedType: TypeInvalidParam,
			expectedMsg:  "Invalid parameter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.createFunc(tt.expectedMsg)
			if err.Type != tt.expectedType || err.Message != tt.expectedMsg {
				t.Errorf("Expected error type: %s, message: %s, got: %s, %s", tt.expectedType, tt.expectedMsg, err.Type, err.Message)
			}
		})
	}
}
