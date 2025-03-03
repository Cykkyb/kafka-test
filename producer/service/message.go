package service

import (
	"consumer/kafka"
	"context"
	"log/slog"
)

type MassageService struct {
	log      *slog.Logger
	producer *kafka.Producer
}

func NewMassageService(producer *kafka.Producer, log *slog.Logger) *MassageService {
	return &MassageService{
		log:      log,
		producer: producer,
	}
}

func (s *MassageService) SendMessage(ctx context.Context, message string) error {
	err := s.producer.Send("test", message)
	if err != nil {
		return err
	}
	return nil
}
