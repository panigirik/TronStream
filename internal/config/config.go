package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Database      Database `env-prefix:"DATABASE_"`
	Tron          TronNode `env-prefix:"TRON_"`
	MigrationsDir string   `env:"MIGRATIONS_DIR" env-default:"./internal/database/migrations"`
}

type Database struct {
	Addr     string `env:"ADDR" env-default:"localhost:5433"`
	Username string `env:"USER"`
	Password string `env:"PASSWORD"`
	Name     string `env:"NAME"`
}

type TronNode struct {
	BaseURL        string        `env:"BASE_URL" env-default:"https://nile.trongrid.io"`
	APIKey         string        `env:"API_KEY"`
	PollerEnabled  bool          `env:"POLLER_ENABLED" env-default:"true"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" env-default:"10s"`
	RequestTimeout time.Duration `env:"REQUEST_TIMEOUT" env-default:"15s"`
	PageLimit      int           `env:"PAGE_LIMIT" env-default:"50"`
	MaxPages       int           `env:"MAX_PAGES" env-default:"20"`
	OnlyConfirmed  bool          `env:"ONLY_CONFIRMED" env-default:"true"`
	StartLookback  time.Duration `env:"START_LOOKBACK" env-default:"1h"`
	TronNodeUrl    string        `env:"TRON_NODE_URL" env-default:"138.226.223.97"`
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("internal/.env")

	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("read config from env: %w", err)
	}
	return &cfg, nil
}

func (c *Config) PostgresURL() string {
	addr := "localhost:5433"
	if len(c.Database.Addr) > 0 {
		addr = c.Database.Addr
	}
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		c.Database.Username,
		c.Database.Password,
		addr,
		c.Database.Name,
	)
}

func (c *Config) PostgresMigrationURL() string {
	return c.PostgresURL()
}
