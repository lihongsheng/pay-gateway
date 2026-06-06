package service

import (
	"context"
	"fmt"
	"github.com/lihongsheng/pay-gateway/global"
	"github.com/lihongsheng/pay-gateway/plugin/payment/api/admin"
	"github.com/lihongsheng/pay-gateway/plugin/payment/domain/entity"
	errors2 "github.com/lihongsheng/pay-gateway/plugin/payment/errors"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	"github.com/lihongsheng/pay-gateway/plugin/payment/service/dto"
	"github.com/lihongsheng/pay-gateway/plugin/payment/svc"
	paySdk "github.com/lihongsheng/payment-sdk"
	params2 "github.com/lihongsheng/payment-sdk/config/params"
	"github.com/lihongsheng/payment-sdk/driver/iface"
	"github.com/lihongsheng/payment-sdk/enum/channel"
	"gorm.io/gorm"
)

type PaymentAccountService interface {
	// Get 后台获取，带有数据验证
	Get(ctx context.Context, id int64) (*admin.PaymentAccountDetail, error)
	// GetByAccountNo 后台获取，带有数据验证
	GetByAccountNo(ctx context.Context, accNo string) (*entity.PaymentAccount, error)
	// Save 后台获取，带有数据验证
	Save(ctx context.Context, app *admin.PaymentAccountCreateRequest) error
	// GetAccountByAppNo 后台获取，带有数据验证-
	GetAccountByAppNo(ctx context.Context, appNo string) ([]*admin.PaymentAccountDetail, error)
	GetApplicationChannelConfig(ctx context.Context, appNo string) ([]*iface.ChannelOption, error)
	GetPaymentProduct(channelCode string) ([]admin.PaymentMethodConfig, error)
}

type paymentAccountService struct {
	svc *svc.ServiceContext
}

func NewPaymentAccountService(svc *svc.ServiceContext) PaymentAccountService {
	return &paymentAccountService{
		svc: svc,
	}
}

func (s *paymentAccountService) GetApplicationAccount(ctx context.Context, appNo string, user *dto.User) ([]*model.PaymentAccount, error) {
	application, err := s.svc.AppRepo.GetByAppNoFormCache(ctx, appNo)
	if err != nil {
		return nil, err
	}
	mch, err := s.svc.MchRepo.GetByMchNo(ctx, application.MchNo)
	if err != nil {
		return nil, err
	}
	err = user.ISHaveMchID(mch.MchNo)
	if err != nil {
		return nil, err
	}
	apps, err := s.svc.PaymentAccountRepo.GetOnlyAppNO(ctx, []string{appNo})
	if err != nil {
		return nil, err
	}
	return apps, nil
}

func (s *paymentAccountService) Save(ctx context.Context, app *admin.PaymentAccountCreateRequest) error {
	if err := app.Validate(); err != nil {
		return err
	}
	application, err := s.svc.AppRepo.GetByAppNo(ctx, app.AppNo)
	if err != nil {
		fmt.Println("-------------------------------------------------------1")
		return err
	}
	if app.MchNo != "" && app.MchNo != application.MchNo {
		return errors2.NewError(errors2.ErrCodeInvalidParam, "无效商户")
	}
	if app.MchNo == "" {
		app.MchNo = application.MchNo
	}
	if err != nil {
		return err
	}
	_, err = s.svc.MchRepo.GetByMchNo(ctx, application.MchNo)
	if err != nil {
		return err
	}
	version := entity.GenAppAccountVersion()
	err = global.GVA_PAY_DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err = s.svc.PaymentAccountRepo.Save(ctx, app, tx)
		if err != nil {
			fmt.Println("-------------------------------------------------------2")
			return err
		}
		return s.svc.AppRepo.UpdateApplicationAccountVersion(ctx, app.AppNo, version, tx)
	})

	if err != nil {
		fmt.Println("-------------------------------------------------------3")
		return err
	}

	err = s.svc.PaymentAccountRepo.RefreshAppAccountsCache(ctx, app.AppNo, version)
	if err != nil {
		fmt.Println("-------------------------------------------------------4")
		return err
	}
	return nil
}

func (s *paymentAccountService) GetByAccountNo(ctx context.Context, accNo string) (*entity.PaymentAccount, error) {
	application, err := s.svc.AppRepo.GetByAppNoFormCache(ctx, accNo)
	if err != nil {
		return nil, err
	}
	_, err = s.svc.MchRepo.GetByMchNo(ctx, application.MchNo)
	if err != nil {
		return nil, err
	}

	return s.svc.PaymentAccountRepo.GetByAccountNoFormCache(ctx, accNo)
}

func (s *paymentAccountService) GetAccountByAppNo(ctx context.Context, appNo string) ([]*admin.PaymentAccountDetail, error) {
	application, err := s.svc.AppRepo.GetByAppNoFormCache(ctx, appNo)
	if err != nil {
		return nil, err
	}
	_, err = s.svc.MchRepo.GetByMchNo(ctx, application.MchNo)
	if err != nil {
		return nil, err
	}
	apps, err := s.svc.PaymentAccountRepo.GetByAppNo(ctx, appNo, nil)
	if err != nil {
		return nil, err
	}
	var result = make([]*admin.PaymentAccountDetail, 0, len(apps))
	for _, app := range apps {
		tmp, err := EntityToParamResponse(app)
		if err != nil {
			return nil, err
		}
		result = append(result, tmp)
	}
	return result, nil
}

func (s *paymentAccountService) Get(ctx context.Context, id int64) (*admin.PaymentAccountDetail, error) {
	app, err := s.svc.PaymentAccountRepo.GetFormCache(ctx, id)
	if err != nil {
		return nil, err
	}
	_, err = s.svc.MchRepo.GetByMchNo(ctx, app.MchNo)
	if err != nil {
		return nil, err
	}
	return EntityToParamResponse(app)
}

func EntityToParamResponse(app *entity.PaymentAccount) (*admin.PaymentAccountDetail, error) {
	payDri, err := paySdk.GetPaymentDriver(channel.Channel(channel.Channel_value[app.Channel]))
	if err != nil {
		return nil, err
	}
	pay, err := paySdk.GetPaymentDriver(channel.Channel(channel.Channel_value[app.Channel]))
	if err != nil {
		return nil, err
	}
	products := payDri.GetSupportProduct()
	options := pay.GetConfigOptions()
	result := &admin.PaymentAccountDetail{
		AppNo:               app.AppNo,
		Channel:             app.Channel,
		ChannelName:         options.Label,
		ChannelConfig:       app.ChannelConfig,
		Extend:              app.Extend,
		ID:                  app.ID,
		MchNo:               app.MchNo,
		AccountNo:           app.AccountNo,
		Name:                app.Name,
		PaymentMethodConfig: BuildPaymentMethod(products, app.PaymentMethod),
		Remark:              app.Remark,
		Status:              app.Status,
		ChannelOption:       options,
		CreatedAt:           app.CreatedAt,
		ValidateStatus:      app.ValidateStatus,
		MaxLimit:            app.MaxLimit,
	}
	return result, nil
}

func BuildPaymentMethod(payDri []iface.PaymentMethod, have []entity.PaymentMethod) []admin.PaymentMethodConfig {
	var result = make([]admin.PaymentMethodConfig, 0, len(have))
	var maps = map[string][]entity.PaymentMethod{}
	for _, item := range have {
		if item2, exists := maps[item.Method]; exists {
			fmt.Println(item2)
			maps[item.Method] = append(maps[item.Method], item)
		} else {
			maps[item.Method] = []entity.PaymentMethod{item}
		}
	}
	for _, item := range payDri {
		r := admin.PaymentMethodConfig{
			Method:  item.Method,
			Label:   item.Label,
			Product: []admin.PaymentProductConfig{},
		}
		for _, item2 := range item.Product {
			tmp := admin.PaymentProductConfig{
				Product: item2.Product,
				Label:   item2.Label,
			}
			if products, exists := maps[item.Method]; exists {
				for _, item3 := range products {
					if item3.Product == item2.Product {
						tmp.Used = true
						break
					}
				}
			}
			r.Product = append(r.Product, tmp)
		}
		result = append(result, r)
	}
	return result
}

func (s *paymentAccountService) GetApplicationChannelConfig(ctx context.Context, appNo string) ([]*iface.ChannelOption, error) {
	var result = []*iface.ChannelOption{}
	wechat, err := paySdk.GetPaymentDriver(channel.Channel_Wechat)
	if err != nil {
		return nil, err
	}
	wechatOptions := wechat.GetConfigOptions()
	result = append(result, wechatOptions)
	alipay, err := paySdk.GetPaymentDriver(channel.Channel_Alipay)
	if err != nil {
		return nil, err
	}
	aliOptions := alipay.GetConfigOptions()
	result = append(result, aliOptions)
	fuiou, err := paySdk.GetPaymentDriver(channel.Channel_Fuiou)
	if err != nil {
		return nil, err
	}
	fuiouOptions := fuiou.GetConfigOptions()
	result = append(result, fuiouOptions)
	return result, nil
}

// GetApplicationExtendConfig 获取商户的扩展配置，一般用于拉卡拉，富友 在收银台页面获取用户openid的微信或者支付宝配置
// 富友，拉卡拉 需要和对应的微信和支付宝一一对应
func (s *paymentAccountService) GetApplicationExtendConfig(ctx context.Context, appNo string, channelCode string) ([]params2.Option, error) {
	if channelCode == channel.Channel_Alipay.String() {
		return nil, nil
	}
	options := []params2.Option{
		{
			Label:        "微信AppID",
			Name:         "app_id",
			Type:         params2.String,
			ValidateReg:  "^wx[0-9a-zA-Z]{16,32}$",
			ValidateType: params2.ValidateReg,
			InputType:    params2.InputText,
			Default:      "",
			Values:       nil,
			Require:      true,
		},
		{
			Label:        "应用Secret",
			Name:         "app_secret",
			Type:         params2.String,
			ValidateReg:  "",
			ValidateType: params2.ValidateString,
			InputType:    params2.InputPassword,
			Default:      "",
			Values:       nil,
			Require:      true,
		},
	}
	return options, nil
}

func (s *paymentAccountService) GetPaymentProduct(channelCode string) ([]admin.PaymentMethodConfig, error) {
	payDri, err := paySdk.GetPaymentDriver(channel.Channel(channel.Channel_value[channelCode]))
	if err != nil {
		return nil, err
	}
	products := payDri.GetSupportProduct()
	var result = make([]admin.PaymentMethodConfig, len(products))
	for i, v := range products {
		result[i] = admin.PaymentMethodConfig{
			Method:  v.Method,
			Label:   v.Label,
			Product: make([]admin.PaymentProductConfig, 0, len(v.Product)),
		}
		for _, v2 := range v.Product {
			result[i].Product = append(result[i].Product, admin.PaymentProductConfig{
				Product: v2.Product,
				Label:   v2.Label,
			})
		}
	}
	return result, nil
}
