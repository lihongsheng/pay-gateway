package dto

import (
	"errors"
	"github.com/lihongsheng/pay-gateway/enum"
)

// User 当前登录用户信息（从 JWT 解析）
type User struct {
	UserID    int64           `json:"user_id"`
	UserType  enum.SystemType `json:"user_type"`
	HaveMchNo string          `json:"have_mch_no"`
}

// ISMch 是否为商户用户
func (u User) ISMch() bool {
	return u.UserType == enum.SystemTypeMch
}

// ISPlatform 是否为平台用户
func (u User) ISPlatform() bool {
	return u.UserType == enum.SystemTypePlatform
}

func (u User) Validate() error {
	if u.ISMch() {
		if u.HaveMchNo == "" {
			return errors.New("商户用户必须指定商户ID")
		}
	}
	return nil
}

func (u User) ISHaveMchID(mchNo string) error {
	if u.ISPlatform() {
		return nil
	}
	if u.HaveMchNo != mchNo {
		return errors.New("商户用户没有此商户权限")
	}
	return nil
}

func (u User) ISHaveMchIDs(mchNos []string) error {
	if u.ISPlatform() {
		return nil
	}
	for _, id := range mchNos {
		if id == u.HaveMchNo {
			return nil
		}
	}
	return errors.New("商户用户没有此商户权限")
}
