package config

const AppVersionLocal = "1.0.0"

func loadDev() *Config {
	return &Config{
		PostgresDSN: getEnv("POSTGRES_DSN", "postgres://postgres:postgres@localhost:5432/uala"),
		ServerPort:  getEnv("PORT", "8080"),
		Env:         ENVLOCAL,
		AppName:     getEnv("APP_NAME", AppName),
		Version:     AppVersionLocal,
	}
}
