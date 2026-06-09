package config

type Kafka struct {
	Brokers  []string `mapstructure:"brokers" json:"brokers" yaml:"brokers"`
	Username string   `mapstructure:"username" json:"username" yaml:"username"`
	Password string   `mapstructure:"password" json:"password" yaml:"password"`
	ClientID string   `mapstructure:"client_id" json:"client_id" yaml:"client_id"`
	NeedAuth bool     `mapstructure:"need_auth" json:"need_auth" yaml:"need_auth"`
	Topic    Topic    `mapstructure:"topic" json:"topic" yaml:"topic"`
}

type Topic struct {
	PaymentStatus      string `mapstructure:"payment_status" json:"payment_status" yaml:"payment_status"`
	RefundStatus       string `mapstructure:"refund_status" json:"refund_status" yaml:"refund_status"`
	PaymentCallback    string `mapstructure:"payment_callback" json:"payment_callback" yaml:"payment_callback"`
	RefundCallback     string `mapstructure:"refund_callback" json:"refund_callback" yaml:"refund_callback"`
	UserLimit          string `mapstructure:"user_limit" json:"user_limit" yaml:"user_limit"`
	RefundNotifyRetry  string `mapstructure:"refund_notify_retry" json:"refund_notify_retry" yaml:"refund_notify_retry"`
	PaymentNotifyRetry string `mapstructure:"payment_notify_retry" json:"payment_notify_retry" yaml:"payment_notify_retry"`
}
