package kafka

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/lihongsheng/pay-gateway/global"
	"github.com/lihongsheng/pay-gateway/queue"
	"go.uber.org/zap"
	"sync"
	"time"
)

type ConsumerConfig struct {
	//Target []string
	Group string
	Topic string
	//UserName string
	Handler queue.Handler
}

// impl sarama.ConsumerGroupHandler
type ConsumerHandle struct {
	Handler queue.Handler
	BaseCtx context.Context
}

// Consumer 消费服务
type Consumer struct {
	consumer []ConsumerConfig
	cancel   context.CancelFunc
	baseCtx  context.Context
}

// NewConsumer 基于consumer配置可以启动多个消费者组
func NewConsumer(consumer []ConsumerConfig) *Consumer {
	return &Consumer{
		consumer: consumer,
	}
}

func getConfig() ([]string, *sarama.Config) {
	kafkaConfig := global.GVA_CONFIG.Kafka
	config := sarama.NewConfig()
	// 阿里云Kafka版本通常为2.2.x或更高，根据实际版本调整
	config.Version = sarama.V3_3_1_0
	config.ClientID = kafkaConfig.ClientID
	config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.CompressionLevel = sarama.CompressionLevelDefault
	config.Producer.RequiredAcks = sarama.NoResponse // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true
	//config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	// config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	//config.Consumer.Fetch.Default = 8 * 1024 * 1024
	config.Consumer.Return.Errors = true
	config.Net.DialTimeout = 30 * time.Second
	config.Net.ReadTimeout = 30 * time.Second
	config.Net.WriteTimeout = 30 * time.Second
	// 设置60秒，看看重平衡问题是否可以解决
	config.Consumer.Group.Rebalance.Timeout = 60 * time.Second
	config.Consumer.Group.Session.Timeout = time.Second * 120
	if kafkaConfig.NeedAuth {
		// 配置SASL认证（阿里云Kafka常用PLAIN机制）
		config.Net.SASL.Enable = true
		config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		config.Net.SASL.User = kafkaConfig.Username // 阿里云通常为"实例ID#用户名"
		config.Net.SASL.Password = kafkaConfig.Password
		// 配置TLS（使用系统默认信任的根证书）
		/*caCert, _ := ioutil.ReadFile("only-4096-ca-cert")
		  certPool := x509.NewCertPool()
		  certPool.AppendCertsFromPEM(caCert)*/
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
			/*RootCAs:            certPool, // 使用手动导入的根证书
			  MinVersion:         tls.VersionTLS12,*/
		}
		// 配置SSL
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = tlsConfig
	}
	return global.GVA_CONFIG.Kafka.Brokers, config
}

// sarama 会根据topic的分区数，自动建立对应的go去消费
// see https://github.com/IBM/sarama/blob/master/consumer_group.go#L171
// see https://github.com/IBM/sarama/blob/master/consumer_group.go#L339
// see https://github.com/IBM/sarama/blob/master/consumer_group.go#L591
// see https://github.com/IBM/sarama/blob/master/consumer_group.go#L655 根据topic 下的分区数量建立go去消费

func (k *Consumer) Start(ctx context.Context) error {
	k.baseCtx = ctx
	cli := make([]struct {
		client sarama.ConsumerGroup
		config ConsumerConfig
	}, len(k.consumer))
	brokers, config := getConfig()
	for key, consumer := range k.consumer {
		if consumer.Group == "" || consumer.Handler == nil || consumer.Topic == "" {
			return errors.New("consumer MUST set Group,Topic,Target consumer MUST set Group,Topic,Target")
		}
		client, err := sarama.NewConsumerGroup(brokers, consumer.Group, config)
		if err != nil {
			return errors.New("consumer error" + fmt.Sprintf("error %v", err))
		}
		cli[key] = struct {
			client sarama.ConsumerGroup
			config ConsumerConfig
		}{
			client: client,
			config: consumer,
		}
	}
	// 生成cancel ，用于通知消费者通知
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(len(cli))
	errChannel := make(chan error, len(cli))
	for _, val := range cli {
		go func(val struct {
			client sarama.ConsumerGroup
			config ConsumerConfig
		}) {
			defer wg.Done()
			// 处理kafka消息
			retry := 3
			for {
				err := val.client.Consume(ctx, []string{val.config.Topic}, &ConsumerHandle{
					Handler: val.config.Handler,
					BaseCtx: k.baseCtx,
				})
				if err != nil {
					retry--
					if retry < 1 {
						global.GVA_LOG.Info("kafkaConsumerStartError", zap.String("error", err.Error()),
							zap.String("topic", val.config.Topic), zap.String("Group", val.config.Group))
						errChannel <- err
						return
					}
				}
				if ctx.Err() != nil {
					return
				}
			}
		}(val)
	}
	for err := range errChannel {
		if err != nil {
			// k.Stop(k.baseCtx)
			return err
		}
	}
	k.cancel = cancel
	wg.Wait()
	close(errChannel)
	return nil
}

func (k *Consumer) Stop() error {
	// 告知kafka 组件停止消费
	if k.cancel != nil {
		k.cancel()
	}
	for _, consumer := range k.consumer {
		if consumer.Handler != nil {
			consumer.Handler.NotifyClose()
		}
	}
	return nil
}

func (c *ConsumerHandle) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *ConsumerHandle) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *ConsumerHandle) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) (err error) {
	// 1. 检查是否支持批量处理
	batchHandler, isBatchMode := c.Handler.(queue.BatchHandler)

	// ===========================
	// 模式 A: 批量消费模式
	// ===========================
	if isBatchMode {
		// 缓冲区
		buffer := make([]*sarama.ConsumerMessage, 0, batchHandler.BatchSize()) // 批次大小
		ticker := time.NewTicker(batchHandler.BatchInterval())                 // 时间间隔
		defer ticker.Stop()

		// 定义提交函数
		flush := func() error {
			if len(buffer) == 0 {
				return nil
			}

			// 提取数据
			msgs := make([]string, len(buffer))
			for i, msg := range buffer {
				msgs[i] = string(msg.Value)
			}

			// 创建上下文
			// 注意：批量模式下，Context 是一次性的，不再针对单条消息
			newCancelCtx, cancel := context.WithTimeout(c.BaseCtx, 30*time.Second)
			defer cancel()

			// 调用批量接口
			err := batchHandler.BatchMessage(newCancelCtx, msgs)
			if err != nil {
				global.GVA_LOG.Error("kafkaBatchConsumerError", zap.Error(err))
				// 策略：如果批量失败，返回错误会导致 Rebalance 和重试。
				// 如果你希望跳过错误数据，这里应该记录日志并返回 nil。
				return err
			}

			// 成功后，标记所有消息（或者只标记最后一条）
			for _, msg := range buffer {
				session.MarkMessage(msg, "")
			}

			// 清空缓冲区
			buffer = buffer[:0]
			return nil
		}

		// 批量模式的主循环
		for {
			select {
			case <-c.BaseCtx.Done():
				_ = flush() // 退出前尝试提交剩余数据
				c.Handler.NotifyClose()
				return nil
			case <-session.Context().Done():
				return nil
			case <-ticker.C:
				// 时间触发提交
				if err := flush(); err != nil {
					return err
				}
			case message, ok := <-claim.Messages():
				if !ok {
					return nil
				}
				global.GVA_LOG.Info("kafkaExtractBatch", zap.String("topic", message.Topic), zap.String("value", string(message.Value)))
				buffer = append(buffer, message)
				// 数量触发提交
				if len(buffer) >= batchHandler.BatchSize() {
					if err := flush(); err != nil {
						return err
					}
				}
			}
		}
	} else {
		// ===========================
		// 模式 B: 单条消费模式 (保持你原有的逻辑完全不变)
		// ===========================

		// 防止某一个处理异常
		defer func() {
			e := recover()
			if e != nil {
				global.GVA_LOG.Info("kafkaConsumerPanic", zap.Any("err", e))
				_, ok := e.(error)
				if !ok {
					err = errors.New("ERROR" + fmt.Sprintf("errro %v", e))
				}
				err = e.(error)
			}
		}()
		// NOTE:
		// Do not move the code below to a goroutine.
		// The `ConsumeClaim` itself is called within a goroutine, see:
		// https://github.com/IBM/sarama/blob/master/consumer_group.go#L27-L29
		for {
			select {
			case <-c.BaseCtx.Done():
				// 优雅退出
				c.Handler.NotifyClose()
				global.GVA_LOG.Info("kafkaConsumerQuit", zap.String("quit", "quit"))
				return nil
			case message, ok := <-claim.Messages():
				if !ok {
					return nil
				}
				// 节省elk 储存空间
				//_ = c.Log.Log(log.LevelInfo, log.DefaultMessageKey, "kafka-msg", "topic", message.Topic, "offset", message.Offset, "Partition", message.Partition)
				// 从kafka消息头里导出 trace信息
				newCancelCtx, cancel := context.WithCancel(c.BaseCtx)
				global.GVA_LOG.Info("kafkaExtract", zap.String("topic", message.Topic), zap.String("value", string(message.Value)))
				//tracer := otel.Tracer("kafka")
				//newCtx, span := tracer.Start(newCancelCtx, "kafka", trace.WithAttributes(), trace.WithSpanKind(trace.SpanKindConsumer))
				//// 处理消息
				msg := string(message.Value)
				// 具体执行消息处理的 Handler
				err = c.Handler.Message(newCancelCtx, msg)
				// 如果返回错误，不提交偏移量
				if err != nil {
					cancel()
					// 记录日志
					//session.Context().Done()
					global.GVA_LOG.Error("kafkaConsumerError", zap.Error(err))
					return err
				}
				cancel()
				// 提交偏移量
				session.MarkMessage(message, "")
			}
		}
	}
}
