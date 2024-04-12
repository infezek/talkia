package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DBHost            string `envconfig:"DB_HOST"`
	DBPort            string `envconfig:"DB_PORT"`
	DBUser            string `envconfig:"DB_USER"`
	DBPassword        string `envconfig:"DB_PASSWORD"`
	DBName            string `envconfig:"DB_NAME"`
	OpenIAToken       string `envconfig:"OPENIA_TOKEN"`
	OpenIAURL         string `envconfig:"OPENIA_URL"`
	TelegramBotToken  string `envconfig:"TELEGRAM_BOT_TOKEN"`
	Environment       string `envconfig:"ENVIRONMENT" default:"production"`
	S3Region          string `envconfig:"S3_REGION"`
	S3AccessKeyID     string `envconfig:"S3_ACCESS_KEY_ID"`
	S3SecretAccessKey string `envconfig:"S3_SECRET_ACCESS"`
	S3BucketName      string `envconfig:"S3_BUCKET_NAME"`
	NewRelicLicense   string `envconfig:"NEW_RELIC_LICENSE"`
	NewRelicAppName   string `envconfig:"NEW_RELIC_APP_NAME"`
	NewRelicEnabled   bool   `envconfig:"NEW_RELIC_ENABLED"`
	BucketImagesURL   string `envconfig:"BUCKET_IMAGES_URL"`
}

func New(namespace string) (*Config, error) {
	var cfg Config
	err := envconfig.Process(namespace, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
