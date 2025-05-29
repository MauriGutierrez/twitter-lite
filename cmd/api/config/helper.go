package config

import "os"

const (
	ENVLOCAL = "local"
	ENVSTG   = "stg"
	ENVLIVE  = "live"

	AppName = "uala-twitter"
)

type Config struct {
	PostgresDSN string
	ServerPort  string
	Env         string
	AppName     string
	Version     string
}

func Load() *Config {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = ENVLOCAL
	}

	switch env {
	case ENVLIVE:
		return loadProd()
	case ENVSTG:
		return loadStaging()
	case ENVLOCAL:
		return loadDev()
	default:
		return loadDev()
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
