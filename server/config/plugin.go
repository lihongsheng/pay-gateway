package config

import "github.com/lihongsheng/pay-gateway/plugin/payment/config"

type Plugin struct {
  PaymentConfig config.Config `json:"payment_config" mapstructure:"payment_config" yaml:"payment_config"`
}
