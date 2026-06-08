package handler

import (
	"encoding/json"
	"net/http"

	"github.com/oshaw1/go-net-test/config"
)

type ConfigHandler struct {
	cfg      *config.Config
	filepath string
}

func NewConfigHandler(cfg *config.Config, filepath string) *ConfigHandler {
	return &ConfigHandler{cfg: cfg, filepath: filepath}
}

type configUpdateRequest struct {
	Dashboard config.DashboardSettings `json:"dashboard"`
	Tests     config.TestConfigs       `json:"tests"`
}

func (h *ConfigHandler) HandleGetConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"dashboard": h.cfg.Dash,
		"tests":     h.cfg.Tests,
	})
}

func (h *ConfigHandler) HandleUpdateConfig(w http.ResponseWriter, r *http.Request) {
	var req configUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	h.cfg.Dash = req.Dashboard
	h.cfg.Tests = req.Tests

	if err := config.Save(h.filepath, h.cfg); err != nil {
		http.Error(w, "Failed to save config: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *ConfigHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.HandleGetConfig(w, r)
	case http.MethodPost:
		h.HandleUpdateConfig(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
