package models

import (
	"os"
	"path/filepath"
)

// Config containts all the (global) configuration needed to make Alexandria run.
// In most cases this struct should be created by using the NewConfig function.
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

// NewConfig reads from environment variables to construct the Config object.
// If an environment variable is not defined a default value will be used instead.
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
