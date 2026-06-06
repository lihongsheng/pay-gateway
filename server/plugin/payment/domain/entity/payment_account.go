package entity

import (
	"fmt"
	"github.com/lihongsheng/pay-gateway/plugin/payment/enum"
	"time"
)

type PaymentAccount struct {
	ID             int64           `gorm:"column:id;type:int;primaryKey;autoIncrement:true" json:"id"`
	Name           string          `gorm:"column:name;type:varchar(100);not null;comment:名字标识" json:"name"`                                                     // 名字标识
	Remark         string          `gorm:"column:remark;type:varchar(200);not null;comment:描述" json:"remark"`                                                   // 描述
	AppNo          string          `gorm:"column:app_no;type:varchar(100);not null;comment:应用编号" json:"app_no"`                                                 // 应用编号
	MchNo          string          `gorm:"column:mch_no;type:varchar(100);not null;comment:商户编号" json:"mch_no"`                                                 // 商户编号
	Channel        string          `gorm:"column:channel;type:varchar(60);not null;comment:支付渠道：Wechat 微信 | Alipay 支付宝 | Lakala 拉卡拉 | Fuiou 富有" json:"channel"` // 支付渠道：Wechat 微信 | Alipay 支付宝 | Lakala 拉卡拉 | Fuiou 富有
	Status         enum.MchStatus  `gorm:"column:status;type:tinyint;not null;default:1;comment:1 上线 2 下线" json:"status"`                                       // 1 上线 2 下线
	Extend         string          `gorm:"column:extend;type:json;not null;comment:各个渠道的私有配置" json:"extend"`                                                    // 各个渠道的私有配置
	MaxLimit       int64           `gorm:"column:max_limit;type:int;not null;default:10;comment:权重" json:"max_limit"`                                           // 权重
	ValidateStatus int64           `gorm:"column:validate_status;type:tinyint;not null;comment:是否验证 0 没有 1 验证通过" json:"validate_status"`                        // 是否验证 0 没有 1 验证通过
	AccountNo      string          `gorm:"column:account_no;type:varchar(100);not null;comment:账户编号" json:"account_no"`                                         // 账户编号
	CreatedAt      time.Time       `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time       `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	ChannelConfig  string          `gorm:"column:channel_config;type:json;not null;comment:各个渠道的私有配置" json:"channel_config"` // 各个渠道的私有配置
	PaymentMethod  []PaymentMethod `json:"payment_method"`
}

type PaymentMethod struct {
	Method  string
	Product string
	Extend  string
}

func (p PaymentAccount) Support(method, product string) bool {
	for _, item := range p.PaymentMethod {
		if item.Method == method && item.Product == product {
			return true
		}
	}
	return false
}

func GenAppAccountVersion() string {
	return fmt.Sprintf("%d", time.Now().UnixMilli())
	//version := ""
	//versiongList := make([]string, 0, len(accounts))
	//sort.Slice(accounts, func(i, j int) bool {
	//	return accounts[i].ID < accounts[j].ID
	//})
	//for _, account := range accounts {
	//	versiongList = append(versiongList, fmt.Sprintf("%d", account.ID))
	//}
	//version = strings.Join(versiongList, ",")
	//version = fmt.Sprintf("%x", version)
}
