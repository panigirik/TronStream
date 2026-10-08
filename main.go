package main

import (
	"TronStream/internal/api"
	"TronStream/internal/config"
	"TronStream/internal/database"
	"TronStream/internal/infrastructure/background_jobs"
	"TronStream/internal/infrastructure/tron_node"
	"TronStream/internal/middleware"
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("загрузка конфигурации: %v", err)
	}

	if err := database.RunMigrations(cfg.PostgresMigrationURL(), cfg.MigrationsDir); err != nil {
		log.Printf("[WARN] миграции не применены: %v", err)
	}

	pool, err := pgxpool.New(ctx, cfg.PostgresURL())
	if err != nil {
		log.Fatalf("подключение к БД: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("БД недоступна: %v", err)
	}

	walletRepository := database.NewWalletRepository(pool)
	depositRepository := database.NewDepositRepository(pool)

	tronClient := tron_node.NewClient(
		cfg.Tron.BaseURL,
		cfg.Tron.RequestTimeout,
		tron_node.WithAPIKey(cfg.Tron.APIKey),
	)

	var jobs []background_jobs.Job
	if cfg.Tron.PollerEnabled {
		jobs = append(jobs, background_jobs.NewDepositPoller(
			walletRepository,
			depositRepository,
			tronClient,
			background_jobs.DepositPollerConfig{
				Interval:      cfg.Tron.PollInterval,
				PageLimit:     cfg.Tron.PageLimit,
				MaxPages:      cfg.Tron.MaxPages,
				OnlyConfirmed: cfg.Tron.OnlyConfirmed,
				StartLookback: cfg.Tron.StartLookback,
			},
		))
	} else {
		log.Println("[Deposit Poller] отключён через TRON_POLLER_ENABLED")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", api.HealthHandler)
	protectedProfileHandler := middleware.AuthMiddleware(http.HandlerFunc(handleProfile))
	mux.Handle("/user/profile", protectedProfileHandler)
	api.NewHandler().Routes(mux)

	scheduler := background_jobs.NewScheduler(jobs...)
	go scheduler.Start(ctx)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("ошибка запуска HTTP-сервера: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("не удалось корректно остановить HTTP-сервер: %v", err)
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
