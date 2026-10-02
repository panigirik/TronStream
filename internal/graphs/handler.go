package graphs

import "net/http"

type Handler struct {
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"status\":\"ok\"}"))
}

func NewHandler() *Handler {
	return &Handler{}
}
