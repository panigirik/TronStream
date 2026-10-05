package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Database      Database `env:"DATABASE" env_defeault:"DATABASE"`
	MigrationsDir string   `env:"MIGRATIONS_DIR" env_default:"./database/migrations"`
}

type Database struct {
	Addr     string `env:"ADDRESS"`
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
	Name     string `env:"NAME"`
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("read config from env: %w", err)
	}
	return &cfg, nil
}

func (c *Config) PostgresMigrationURL() string {
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
