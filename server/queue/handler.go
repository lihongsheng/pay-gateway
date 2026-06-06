package queue

import (
	"context"
	"time"
)

// kafka 处理消息相关

type Server interface {
	Start(ctx context.Context) error
	Stop() error
}

type Handler interface {
	// Message 消费消息
	Message(ctx context.Context, msg string) error
	// NotifyClose 退出通知
	NotifyClose()
}

// BatchHandler 可选接口：如果 Handler 实现了这个接口，消费者会自动启用批量模式
type BatchHandler interface {
	// BatchMessage 处理一批消息
	// msgs: 消息内容列表
	// return: 如果返回 error，整个批次会重试（或根据逻辑丢弃）
	BatchMessage(ctx context.Context, msgs []string) error
	// NotifyClose 退出通知
	NotifyClose()
	// BatchSize 批量大小
	BatchSize() int
	// BatchInterval 批量间隔
	BatchInterval() time.Duration
}
