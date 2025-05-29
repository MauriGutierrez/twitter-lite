package health

import (
	"encoding/json"
	"net/http"
	"sync"
)

type ServiceInfo struct {
	Env     string `json:"env"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

type HealthHandler struct {
	info     ServiceInfo
	initOnce sync.Once
}

func NewHealthHandler(env, name, version string) *HealthHandler {
	return &HealthHandler{
		info: ServiceInfo{
			Env:     env,
			Name:    name,
			Version: version,
		},
	}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(h.info)
}
