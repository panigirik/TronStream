package main

import (
	"TronStream/internal/api"
	"TronStream/internal/middleware"
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", api.HealthHandler)
	protectedProfileHandler := middleware.AuthMiddleware(http.HandlerFunc(handleProfile))
	mux.Handle("/user/profile", protectedProfileHandler)
	api.NewHandler().Routes(mux)

	go func() {
		if err := http.ListenAndServe(":8080", mux); err != nil {
			fmt.Println("Error starting HTTP server:", err)
		}
	}()

	select {}

}

func handleProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIdKey).(string)
	if !ok {
		http.Error(w, `{"error": "Unauthorized / User ID not found in context"}`, http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"message": "Успешный доступ к профилю!", "user_id": "%s"}`, userID)))
}
