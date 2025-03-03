package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"os"
	"strings"
)

type Config struct {
	App   AppConfig
	Kafka KafkaConfig
	DB    DBConfig
}

type AppConfig struct {
	Env           string
	MigrationPath string
}

type KafkaConfig struct {
	Brokers       []string
	Topic         string
	ConsumerGroup string
}

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBname   string
	SSL      string
}

func MustLoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("failed to load .env file: " + err.Error())
	}

	var cfg Config

	cfg.App.Env = getEnv("APP_ENV")
	cfg.App.MigrationPath = getEnv("APP_MIGRATION_PATH")

	cfg.DB.Host = getEnv("DB_HOST")
	cfg.DB.Port = getEnv("DB_PORT")
	cfg.DB.Username = getEnv("DB_USER")
	cfg.DB.Password = getEnv("DB_PASSWORD")
	cfg.DB.DBname = getEnv("DB_NAME")
	cfg.DB.SSL = getEnv("DB_SSL")

	cfg.Kafka.Brokers = strings.Split(getEnv("KAFKA_BROKERS"), ",")
	cfg.Kafka.Topic = getEnv("KAFKA_TOPIC")
	cfg.Kafka.ConsumerGroup = getEnv("KAFKA_CONSUMER_GROUP")

	if err = cleanenv.ReadEnv(&cfg); err != nil {
		panic("failed to load config from environment: " + err.Error())
	}

	return &cfg
}

func getEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic("environment variable not found: " + key)
	}
	return value
}
