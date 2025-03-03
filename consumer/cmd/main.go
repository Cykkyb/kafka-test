package main

import (
	"consumer/app"
	"consumer/config"
	"consumer/kafka"
	"consumer/lib/logger"
	"consumer/repository"
	"consumer/service"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoadConfig()

	log := initLogger()
	log.Info("Config loaded",
		slog.Any("cfg", cfg),
	)

	defer func() {
		if r := recover(); r != nil {
			log.Error("Panic", slog.String("err", fmt.Sprintf("%v", r)))
		}
	}()

	db, err := repository.ConnectDb(repository.Config{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Username: cfg.DB.Username,
		Password: cfg.DB.Password,
		DBname:   cfg.DB.DBname,
		SSL:      cfg.DB.SSL,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to connect to db: %s", err))
	}

	rep := repository.NewRepository(db)
	serv := service.NewService(rep, log)
	consumer, err := kafka.NewConsumer(cfg.Kafka.Brokers, cfg.Kafka.ConsumerGroup, cfg.Kafka.Topic, serv, log)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to kafka: %s", err))
	}

	application := app.NewApp(log, consumer)
	go func() {
		err = application.Run()
		if err != nil {
			log.Error("failed to run application", slog.String("err", err.Error()))
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.Stop()
}

func initLogger() *slog.Logger {
	opts := logger.PrettyHandlerOptions{
		SlogOpts: slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := logger.NewPrettyHandler(os.Stdout, opts)
	log := slog.New(handler)

	return log
}
