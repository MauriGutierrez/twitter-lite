package config

const AppVersionProd = "1.0.0"

func loadProd() *Config {
	return &Config{
		PostgresDSN: getEnv("POSTGRES_DSN", "postgres://produser:secure@prod-host:5432/uala"),
		ServerPort:  getEnv("PORT", "80"),
		Env:         ENVLIVE,
		AppName:     getEnv("APP_NAME", AppName),
		Version:     AppVersionProd,
	}
}
