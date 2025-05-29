package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad_Table(t *testing.T) {
	origEnv := os.Getenv("APP_ENV")
	defer os.Setenv("APP_ENV", origEnv)

	tests := []struct {
		name           string
		setEnv         string
		expectedEnvVal string
		loadFunc       func() *Config
	}{
		{
			name:           "uses loadDev for empty APP_ENV",
			setEnv:         "",
			expectedEnvVal: ENVLOCAL,
			loadFunc:       loadDev,
		},
		{
			name:           "uses loadDev for APP_ENV = local",
			setEnv:         ENVLOCAL,
			expectedEnvVal: ENVLOCAL,
			loadFunc:       loadDev,
		},
		{
			name:           "uses loadStaging for APP_ENV = stg",
			setEnv:         ENVSTG,
			expectedEnvVal: ENVSTG,
			loadFunc:       loadStaging,
		},
		{
			name:           "uses loadProd for APP_ENV = live",
			setEnv:         ENVLIVE,
			expectedEnvVal: ENVLIVE,
			loadFunc:       loadProd,
		},
		{
			name:           "uses loadDev for unknown APP_ENV",
			setEnv:         "notarealenv",
			expectedEnvVal: ENVLOCAL,
			loadFunc:       loadDev,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv == "" {
				os.Unsetenv("APP_ENV")
			} else {
				os.Setenv("APP_ENV", tc.setEnv)
			}
			cfg := Load()
			assert.Equal(t, tc.expectedEnvVal, cfg.Env)
		})
	}

	t.Run("loadDev returns local config", func(t *testing.T) {
		cfg := loadDev()
		assert.Equal(t, ENVLOCAL, cfg.Env)
	})
	t.Run("loadStaging returns stg config", func(t *testing.T) {
		cfg := loadStaging()
		assert.Equal(t, ENVSTG, cfg.Env)
	})
	t.Run("loadProd returns live config", func(t *testing.T) {
		cfg := loadProd()
		assert.Equal(t, ENVLIVE, cfg.Env)
	})
}
