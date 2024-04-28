package config

import "os"

type Config struct {
	Postgres *PostgresCfg
	App      *AppConfig
}

type PostgresCfg struct {
	User     string
	Password string
	Host     string
	Port     string
	Sslmode  string
	Dbname   string
}

type AppConfig struct {
	Port      string
	JWTSecret string
}

func LoadConfig() *Config {

	pgCfg := loadPostgresConfig()
	appCfg := loadAppConfig()

	return &Config{
		Postgres: pgCfg,
		App:      appCfg,
	}
}

func loadPostgresConfig() *PostgresCfg {

	return &PostgresCfg{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Dbname:   os.Getenv("DB_NAME"),
		Sslmode:  os.Getenv("SSL_MODE"),
	}
}

func loadAppConfig() *AppConfig {

	return &AppConfig{
		Port:      os.Getenv("APP_PORT"),
		JWTSecret: os.Getenv("JWTSECRET"),
	}
}
