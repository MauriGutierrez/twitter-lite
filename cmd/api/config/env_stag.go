package config

const AppVersionStaging = "1.0.0"

func loadStaging() *Config {
	return &Config{
		PostgresDSN: getEnv("POSTGRES_DSN", "postgres://staguser:stagpass@staging-host:5432/uala"),
		ServerPort:  getEnv("PORT", "8081"),
		Env:         ENVSTG,
		AppName:     getEnv("APP_NAME", AppName),
		Version:     AppVersionStaging,
	}
}
