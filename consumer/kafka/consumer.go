package kafka

import (
	"consumer/service"
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"log/slog"
	"strings"
)

type MessageConsumer struct {
	Consumer sarama.ConsumerGroup
	service  *service.Service
	topic    string
	log      *slog.Logger
}

func NewConsumer(address []string, consumerGroup, topic string, service *service.Service, log *slog.Logger) (MessageConsumer, error) {
	cfg := sarama.NewConfig()
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumer, err := sarama.NewConsumerGroup(address, consumerGroup, cfg)
	if err != nil {
		return MessageConsumer{}, err
	}
	log.Info("Kafka consumer created",
		slog.Any("address", address),
		slog.String("consumerGroup", consumerGroup),
		slog.String("topic", topic),
	)

	return MessageConsumer{
		consumer,
		service,
		topic,
		log,
	}, nil
}

func (c *MessageConsumer) Start() error {
	if c.Consumer == nil {
		return fmt.Errorf("consumer is nil")
	}
	err := c.Consumer.Consume(context.Background(), strings.Split(c.topic, ","), c)
	if err != nil {
		return err
	}
	return nil
}

func (c *MessageConsumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *MessageConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *MessageConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				c.log.Info("Message channel was closed")
				return nil
			}
			c.log.Info("Claimed message",
				slog.String("topic", message.Topic),
				slog.Int64("partition", int64(message.Partition)),
				slog.Int64("offset", message.Offset),
				slog.String("value", string(message.Value)),
				slog.String("key", string(message.Key)),
			)
			err := c.service.AddMessage(session.Context(), string(message.Value))
			if err != nil {
				c.log.Error("Failed to add message", err)
			}
			c.log.Info("Message added", slog.String("value", string(message.Value)))
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			c.log.Info("Session was closed")
			return nil
		}
	}
}
