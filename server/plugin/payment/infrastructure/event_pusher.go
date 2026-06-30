package infrastructure

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/lihongsheng/pay-gateway/plugin/payment/domain/event"
)

type Event interface {
	Publish(ctx context.Context, event event.Event, topic string) error
}

type eventImpl struct {
	producer sarama.SyncProducer
}

func NewEvent(producer sarama.SyncProducer) Event {
	return &eventImpl{producer: producer}
}

func NewNoopEvent() Event {
	return &noopEvent{}
}

func (e *eventImpl) Publish(ctx context.Context, event event.Event, topic string) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(data),
		Key:   sarama.StringEncoder(event.GetEventKey()),
	}
	_, _, err = e.producer.SendMessage(msg)
	return err
}

// noopEvent is a no-op implementation when Kafka is not configured
type noopEvent struct{}

func (n *noopEvent) Publish(_ context.Context, _ event.Event, _ string) error {
	return nil
}
