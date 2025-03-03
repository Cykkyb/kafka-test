package service

import (
	"consumer/repository"
	"context"
	"log/slog"
)

type MassageService struct {
	log        *slog.Logger
	repository repository.MessageRepository
}

func NewMassageService(repo repository.MessageRepository, log *slog.Logger) *MassageService {
	return &MassageService{
		log:        log,
		repository: repo,
	}
}

func (s *MassageService) AddMessage(ctx context.Context, message string) error {
	s.log.Info("Message added", slog.String("message", message))
	err := s.repository.SaveMessage(ctx, message)
	if err != nil {
		return err
	}
	return nil
}
