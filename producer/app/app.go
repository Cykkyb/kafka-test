package app

import (
	"consumer/kafka"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type App struct {
	log        *slog.Logger
	port       int
	HttpServer *http.Server
	Producer   *kafka.Producer
}

func NewApp(log *slog.Logger, port int, routes *http.ServeMux, producer *kafka.Producer) *App {
	app := App{
		log:      log,
		port:     port,
		Producer: producer,
	}

	if routes != nil {
		app.HttpServer = &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: routes,
		}
	}

	return &app
}

func (a *App) Run() {
	go func() {
		err := a.startHttpServer()
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		a.log.Info("Starting producer")
		for i := range 100 {
			err := a.Producer.Send("test", fmt.Sprintf("message %d", i))
			if err != nil {
				a.log.Error("failed to send message", slog.String("err", err.Error()))
			} else {
				a.log.Info("message sent", slog.String("message", fmt.Sprintf("message %d", i)))
			}
		}
	}()
}

func (a *App) startHttpServer() error {
	if a.HttpServer == nil {
		a.log.Info("HttpServer is not configured")
		return nil
	}
	a.log.Info("Starting http server port ", a.port)

	if err := a.HttpServer.ListenAndServe(); err != nil {
		return fmt.Errorf("%s: %w", "", err)
	}

	return nil
}

func (a *App) Stop() {
	a.stopHttpServer()
	a.stopKafka()
}

func (a *App) stopKafka() {
	if a.Producer == nil {
		a.log.Info("Consumer is not configured")
		return
	}

	a.log.Info("Stopping kafka consumer")
	a.Producer.Close()
}

func (a *App) stopHttpServer() {
	if a.HttpServer == nil {
		a.log.Info("HttpServer is not configured")
		return
	}
	a.log.Info("Stopping http server")

	ctx, canc := context.WithTimeout(context.Background(), 5*time.Second)
	defer canc()

	if err := a.HttpServer.Shutdown(ctx); err != nil {
		a.log.Error("failed to stop http server", slog.String("err", err.Error()))
	}
}
