package app

import (
	"consumer/kafka"
	"log/slog"
)

type App struct {
	log             *slog.Logger
	MessageConsumer kafka.MessageConsumer
}

func NewApp(log *slog.Logger, consumer kafka.MessageConsumer) *App {
	return &App{
		log:             log,
		MessageConsumer: consumer,
	}
}

func (a *App) Run() error {
	if err := a.MessageConsumer.Start(); err != nil {
		return err
	}

	return nil
}

func (a *App) Stop() {
	a.log.Info("Stopping kafka consumer")
	if err := a.MessageConsumer.Consumer.Close(); err != nil {
		a.log.Error("failed to stop kafka consumer", slog.String("err", err.Error()))
	}
}
