package service

import (
	"consumer/repository"
	"context"
	"log/slog"
)

type Message interface {
	AddMessage(ctx context.Context, message string) error
}

type Service struct {
	Message
}

func NewService(repo repository.MessageRepository, log *slog.Logger) *Service {
	return &Service{
		NewMassageService(repo, log),
	}
}
