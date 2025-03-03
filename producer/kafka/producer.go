package kafka

import (
	"github.com/IBM/sarama"
	"log/slog"
)

type Producer struct {
	sarama.SyncProducer
	topic string
	log   *slog.Logger
}

func NewProducer(address []string, topic string, log *slog.Logger) (*Producer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(address, cfg)
	if err != nil {
		return nil, err
	}

	return &Producer{
		producer,
		topic,
		log,
	}, nil
}

func (p *Producer) Close() error {
	return p.SyncProducer.Close()
}

func (p *Producer) Send(key, value string) error {
	p.log.Info("Sending message", slog.String("key", key), slog.String("value", value))

	_, _, err := p.SyncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(value),
	})

	return err
}
