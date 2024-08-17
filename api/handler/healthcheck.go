package handler

import (
	"net/http"
)

type UtilHandler struct {
}

func (h UtilHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("server is alive"))
}
