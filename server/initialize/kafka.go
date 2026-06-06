package initialize

import (
	"crypto/tls"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/lihongsheng/pay-gateway/global"
	"log"
)

func InitKafka() (sarama.SyncProducer, error) {
	return createAliyunKafkaClientV2()
}

func createAliyunKafkaClientV2() (sarama.SyncProducer, error) {
	kafkaConfig := global.GVA_CONFIG.Kafka
	// 配置Sarama
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

	// 可选：开启调试日志
	sarama.Logger = log.Default()
	// 创建客户端
	client, err := sarama.NewSyncProducer(kafkaConfig.Brokers, config)
	if err != nil {
		return nil, fmt.Errorf("创建Kafka客户端失败: %v", err)
	}

	return client, nil
}
