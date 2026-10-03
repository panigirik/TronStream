package main

import (
	"TronStream/internal/api"
	"TronStream/internal/config"
	"TronStream/internal/database"
	"TronStream/internal/infrastructure/background_jobs"
	"TronStream/internal/middleware"
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfig()
	if err != nil {
		return
	}

	if err := database.RunMigrations(cfg.PostgresMigrationURL(), cfg.MigrationsDir); err != nil {
		return //TODO добавить логирование
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", api.HealthHandler)
	protectedProfileHandler := middleware.AuthMiddleware(http.HandlerFunc(handleProfile))
	mux.Handle("/user/profile", protectedProfileHandler)
	api.NewHandler().Routes(mux)
	sheduler := background_jobs.NewScheduler()
	go sheduler.Start(ctx)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Println("Error starting HTTP server:", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("HTTP server Shutdown Failed:%+v\n", err)
	}

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
