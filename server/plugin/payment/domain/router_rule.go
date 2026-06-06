package domain

import (
	"context"
	errors2 "github.com/lihongsheng/pay-gateway/plugin/payment/errors"
	"github.com/lihongsheng/pay-gateway/plugin/payment/service/dto"
)

type RuleEngine interface {
	Evaluate(ctx context.Context, conditions []*dto.RuleConditions) ([]string, error)
}

type RuleEngineService struct {
	ruleMap map[string]Rule
}

func NewRuleEngineService() RuleEngine {
	ruleMap := make(map[string]Rule)
	ruleMap["BaseLimitRule"] = &BaseLimitRule{}
	ruleMap["RequestLimitRule"] = &RequestLimitRule{}
	return &RuleEngineService{
		ruleMap: ruleMap,
	}
}
func (r *RuleEngineService) Evaluate(ctx context.Context, conditions []*dto.RuleConditions) ([]string, error) {
	var accountNo = make([]string, 0, len(conditions))
	for _, item := range conditions {
		exists := true
		for _, rule := range r.ruleMap {
			if !(rule.IsEnabled() && rule.Match(item)) {
				exists = false
			}
		}
		if exists {
			accountNo = append(accountNo, item.AccountNO)
		}
	}
	if len(accountNo) == 0 {
		return nil, errors2.NewError(errors2.ErrCodeInvalidParam, "未找到可用的支付账户")
	}
	return accountNo, nil
}

type Rule interface {
	Name() string
	IsEnabled() bool
	Match(req *dto.RuleConditions) bool
}

type BaseLimitRule struct {
}

func (r *BaseLimitRule) Name() string {
	return "BaseLimitRule"
}

func (r *BaseLimitRule) IsEnabled() bool {
	return true
}

func (r *BaseLimitRule) Match(condition *dto.RuleConditions) bool {
	if condition.UserLimit || condition.PaymentLimit {
		return false
	}
	return true
}

type RequestLimitRule struct {
}

func (r *RequestLimitRule) Name() string {
	return "RequestLimitRule"
}

func (r *RequestLimitRule) IsEnabled() bool {
	return true
}

func (r *RequestLimitRule) Match(condition *dto.RuleConditions) bool {
	// 过滤
	if len(condition.FilterAccountNO) > 0 {
		for _, item := range condition.FilterAccountNO {
			if item == condition.AccountNO {
				return false
			}
		}
	}
	if condition.MaxLimit > 0 && condition.RequestSuccess >= condition.MaxLimit {
		return false
	}

	return true
}
