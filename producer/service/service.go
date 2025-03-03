package service

import (
	"consumer/kafka"
	"context"
	"log/slog"
)

type Message interface {
	SendMessage(ctx context.Context, message string) error
}

type Service struct {
	Message
}

func NewService(producer *kafka.Producer, log *slog.Logger) *Service {
	return &Service{
		NewMassageService(producer, log),
	}
}
