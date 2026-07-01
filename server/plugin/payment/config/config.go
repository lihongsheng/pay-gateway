package config

type Config struct {
	// 域名
	ApiHost string `mapstructure:"api_host" json:"api_host" yaml:"api_host"`
	// 代理前缀地址
	ProxyNotifyPrefix string `mapstructure:"proxy_notify_prefix" json:"proxy_notify_prefix" yaml:"proxy_notify_prefix"`
	// 前端域名
	WebHost string `mapstructure:"web_host" json:"web_host" yaml:"web_host"`
	// salt
	Salt  string `mapstructure:"salt" json:"salt" yaml:"salt"`
	Env   string `mapstructure:"env" json:"env" yaml:"env"`
	Topic Topic  `mapstructure:"topic" json:"topic" yaml:"topic"`
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
