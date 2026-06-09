package infrastructure

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/lihongsheng/pay-gateway/global"
	"github.com/lihongsheng/pay-gateway/plugin/payment/domain/event"
)

type Event interface {
	Publish(ctx context.Context, event event.Event, topic string) error
}
type eventImpl struct {
}

func NewEvent() Event {
	return &eventImpl{}
}

func (e *eventImpl) Publish(ctx context.Context, event event.Event, topic string) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	// 写入kafka
	msg := &sarama.ProducerMessage{
		Topic: topic,                                     // 主题消息
		Value: sarama.StringEncoder(data),                // 消息内容（支持多种Encoder）
		Key:   sarama.StringEncoder(event.GetEventKey()), // 可选：添加消息键（用于分区路由和去重）
	}
	_, _, err = global.GVA_KAFKA_PRODUCER.SendMessage(msg)
	return err
}
