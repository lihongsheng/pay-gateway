package config

import "github.com/lihongsheng/pay-gateway/config"

var Config = &Con{}

type Con struct {
	// 域名
	ApiHost string `mapstructure:"api_host" json:"api_host" yaml:"api_host"`
	// 代理前缀地址
	ProxyNotifyPrefix string `mapstructure:"proxy_notify_prefix" json:"proxy_notify_prefix" yaml:"proxy_notify_prefix"`
	// 前端域名
	WebHost string `mapstructure:"web_host" json:"web_host" yaml:"web_host"`
	// salt
	Salt  string `mapstructure:"salt" json:"salt" yaml:"salt"`
	Env   string `mapstructure:"env" json:"env" yaml:"env"`
	Topic config.Topic
}
