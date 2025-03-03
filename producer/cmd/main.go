package main

import (
	"consumer/app"
	"consumer/config"
	"consumer/handlers"
	"consumer/kafka"
	"consumer/lib/logger"
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

	producer, err := kafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic, log)
	serv := service.NewService(producer, log)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to kafka: %s", err))
	}
	handler := handlers.NewHandler(serv, log)

	application := app.NewApp(log, cfg.App.Port, handler.InitRoutes(), producer)
	go func() {
		application.Run()
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
