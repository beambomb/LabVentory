package config

import "os"

// Config holds runtime configuration loaded from environment variables.
type Config struct {
	TableName    string
	ImagesBucket string
	AWSRegion    string
	Port         string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() Config {
	return Config{
		TableName:    getEnv("TABLE_NAME", "LabComponentInventory"),
		ImagesBucket: getEnv("IMAGES_BUCKET", ""),
		AWSRegion:    getEnv("AWS_REGION", "ap-southeast-1"),
		Port:         getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
