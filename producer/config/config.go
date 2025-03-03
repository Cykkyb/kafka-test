package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	App   AppConfig
	Kafka KafkaConfig
}

type AppConfig struct {
	Port int
	Env  string
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
}

func MustLoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("failed to load .env file: " + err.Error())
	}

	var cfg Config

	cfg.App.Port, _ = strconv.Atoi(getEnv("APP_PORT_PUBLISH"))
	cfg.App.Env = getEnv("APP_ENV")

	cfg.Kafka.Brokers = strings.Split(getEnv("KAFKA_BROKERS"), ",")
	cfg.Kafka.Topic = getEnv("KAFKA_TOPIC")

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
