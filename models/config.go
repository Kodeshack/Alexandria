package models

import (
	"os"
	"path/filepath"
)

type Config struct {
	ContentPath       string
	UserStoragePath   string
	TemplateDirectory string
	AssetPath         string
	Host              string
	Port              string
	BaseURL           string
}

func getEnvVar(n, def string) string {
	val, ok := os.LookupEnv(n)
	if ok {
		return val
	}

	return def
}

/// This will read from environment variables.
func NewConfig() *Config {
	dataPath := getEnvVar("ALEXANDRIA_DATA_PATH", "data")
	return &Config{
		ContentPath:       filepath.Join(dataPath, "content"),
		UserStoragePath:   filepath.Join(dataPath, "users.db"),
		TemplateDirectory: getEnvVar("ALEXANDRIA_TEMPLATE_DIR", "view/templates"),
		AssetPath:         getEnvVar("ALEXANDRIA_ASSET_DIR", "assets/public"),
		Host:              getEnvVar("ALEXANDRIA_HOST", "localhost"),
		Port:              getEnvVar("ALEXANDRIA_PORT", ":8080"),
		BaseURL:           getEnvVar("ALEXANDRIA_BASE_URL", "http://localhost:8080/"),
	}
}
