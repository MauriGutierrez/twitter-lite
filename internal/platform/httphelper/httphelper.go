package httphelper

import (
	"encoding/json"
	"errors"
	"net/http"

	"ualaTwitter/internal/platform/errors/usecase"
)

func RenderError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

func StatusFromError(err error) int {
	var ue *usecase.UseCaseError
	if errors.As(err, &ue) {
		switch ue.Type {
		case usecase.TypeInvalidParam:
			return http.StatusBadRequest
		case usecase.TypeNotFound:
			return http.StatusNotFound
		case usecase.TypeForbidden:
			return http.StatusForbidden
		case usecase.TypeConflict:
			return http.StatusConflict
		default:
			return http.StatusInternalServerError
		}
	}
	return http.StatusInternalServerError
}
