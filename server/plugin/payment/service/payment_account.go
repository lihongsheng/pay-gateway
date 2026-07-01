package service

import (
	"context"
	"fmt"
	"github.com/lihongsheng/pay-gateway/repo/system"

	"github.com/lihongsheng/pay-gateway/plugin/payment/domain/entity"
	"github.com/lihongsheng/pay-gateway/plugin/payment/dto"
	errors2 "github.com/lihongsheng/pay-gateway/plugin/payment/errors"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo"
	"github.com/lihongsheng/pay-gateway/plugin/payment/repo/model"
	paySdk "github.com/lihongsheng/payment-sdk"
	"github.com/lihongsheng/payment-sdk/driver/iface"
	"github.com/lihongsheng/payment-sdk/enum/channel"
	"gorm.io/gorm"
)

type PaymentAccountService interface {
	// Get 后台获取，带有数据验证
	Get(ctx context.Context, id int64) (*dto.PaymentAccountDetail, error)
	// GetByAccountNo 后台获取，带有数据验证
	GetByAccountNo(ctx context.Context, accNo string) (*entity.PaymentAccount, error)
	// Save 后台获取，带有数据验证
	Save(ctx context.Context, app *dto.PaymentAccountCreateRequest) error
	// GetAccountByAppNo 后台获取，带有数据验证-
	GetAccountByAppNo(ctx context.Context, appNo string) ([]*dto.PaymentAccountDetail, error)
	GetApplicationChannelConfig(ctx context.Context, appNo string) ([]*iface.ChannelOption, error)
	GetPaymentProduct(channelCode string) ([]dto.PaymentMethodConfig, error)
}

type paymentAccountService struct {
	appRepo            repo.ApplicationRepo
	mchRepo            system.MchRepo
	paymentAccountRepo repo.PaymentAccountRepo
	db                 *gorm.DB
}

func NewPaymentAccountService(appRepo repo.ApplicationRepo, mchRepo system.MchRepo, paymentAccountRepo repo.PaymentAccountRepo, db *gorm.DB) PaymentAccountService {
	return &paymentAccountService{
		appRepo:            appRepo,
		mchRepo:            mchRepo,
		paymentAccountRepo: paymentAccountRepo,
		db:                 db,
	}
}

func (s *paymentAccountService) GetApplicationAccount(ctx context.Context, appNo string, user *dto.User) ([]*model.PaymentAccount, error) {
	application, err := s.appRepo.GetByAppNoFormCache(ctx, appNo)
	if err != nil {
		return nil, err
	}
	mch, err := s.mchRepo.GetByMchNo(ctx, application.MchNo)
	if err != nil {
		return nil, err
	}
	err = user.ISHaveMchID(mch.MchNo)
	if err != nil {
		return nil, err
	}
	apps, err := s.paymentAccountRepo.GetOnlyAppNO(ctx, []string{appNo})
	if err != nil {
		return nil, err
	}
	return apps, nil
}

func (s *paymentAccountService) Save(ctx context.Context, app *dto.PaymentAccountCreateRequest) error {
	if err := app.Validate(); err != nil {
		return err
	}
	application, err := s.appRepo.GetByAppNo(ctx, app.AppNo)
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
	_, err = s.mchRepo.GetByMchNo(ctx, application.MchNo)
	if err != nil {
		return err
	}
	version := entity.GenAppAccountVersion()
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err = s.paymentAccountRepo.Save(ctx, app, tx)
		return s.appRepo.UpdateApplicationAccountVersion(ctx, app.AppNo, version, tx)
	})
	if err != nil {
		return err
	}
	err = s.paymentAccountRepo.RefreshAppAccountsCache(ctx, app.AppNo, version)
	if err != nil {
		return err
	}
	return nil
}

func (s *paymentAccountService) GetByAccountNo(ctx context.Context, accNo string) (*entity.PaymentAccount, error) {
	application, err := s.appRepo.GetByAppNoFormCache(ctx, accNo)
	if err != nil {
		return nil, err
	}
	_, err = s.mchRepo.GetByMchNo(ctx, application.MchNo)
	if err != nil {
		return nil, err
	}

	return s.paymentAccountRepo.GetByAccountNoFormCache(ctx, accNo)
}

func (s *paymentAccountService) GetAccountByAppNo(ctx context.Context, appNo string) ([]*dto.PaymentAccountDetail, error) {
	application, err := s.appRepo.GetByAppNoFormCache(ctx, appNo)
	if err != nil {
		return nil, err
	}
	_, err = s.mchRepo.GetByMchNo(ctx, application.MchNo)
	if err != nil {
		return nil, err
	}
	apps, err := s.paymentAccountRepo.GetByAppNo(ctx, appNo, nil)
	if err != nil {
		return nil, err
	}
	var result = make([]*dto.PaymentAccountDetail, 0, len(apps))
	for _, app := range apps {
		tmp, err := EntityToParamResponse(app)
		if err != nil {
			return nil, err
		}
		result = append(result, tmp)
	}
	return result, nil
}

func (s *paymentAccountService) Get(ctx context.Context, id int64) (*dto.PaymentAccountDetail, error) {
	app, err := s.paymentAccountRepo.GetFormCache(ctx, id)
	if err != nil {
		return nil, err
	}
	_, err = s.mchRepo.GetByMchNo(ctx, app.MchNo)
	if err != nil {
		return nil, err
	}
	return EntityToParamResponse(app)
}

func EntityToParamResponse(app *entity.PaymentAccount) (*dto.PaymentAccountDetail, error) {
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
	result := &dto.PaymentAccountDetail{
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
		Status:              int(app.Status),
		ChannelOption:       options,
		CreatedAt:           app.CreatedAt,
		ValidateStatus:      app.ValidateStatus,
		MaxLimit:            app.MaxLimit,
	}
	return result, nil
}

func BuildPaymentMethod(payDri []iface.PaymentMethod, have []entity.PaymentMethod) []dto.PaymentMethodConfig {
	var result = make([]dto.PaymentMethodConfig, 0, len(have))
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
		r := dto.PaymentMethodConfig{
			Method:  item.Method,
			Label:   item.Label,
			Product: []dto.PaymentProductConfig{},
		}
		for _, item2 := range item.Product {
			tmp := dto.PaymentProductConfig{
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

func (s *paymentAccountService) GetPaymentProduct(channelCode string) ([]dto.PaymentMethodConfig, error) {
	payDri, err := paySdk.GetPaymentDriver(channel.Channel(channel.Channel_value[channelCode]))
	if err != nil {
		return nil, err
	}
	products := payDri.GetSupportProduct()
	var result = make([]dto.PaymentMethodConfig, len(products))
	for i, v := range products {
		result[i] = dto.PaymentMethodConfig{
			Method:  v.Method,
			Label:   v.Label,
			Product: make([]dto.PaymentProductConfig, 0, len(v.Product)),
		}
		for _, v2 := range v.Product {
			result[i].Product = append(result[i].Product, dto.PaymentProductConfig{
				Product: v2.Product,
				Label:   v2.Label,
			})
		}
	}
	return result, nil
}
