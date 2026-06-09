package dto

import (
	"errors"
	"github.com/lihongsheng/pay-gateway/enum"
)

type User struct {
	UserID    int64         `json:"user_id"`
	UserType  enum.UserType `json:"user_type"`
	HaveMchNo string        `json:"have_mch_no"`
}

func (u User) Validate() error {
	if u.UserType.ISMch() {
		if u.HaveMchNo == "" {
			return errors.New("商户用户必须指定商户ID")
		}
	}
	return nil
}

func (u User) ISHaveMchID(mchNo string) error {
	if u.UserType.ISPlatform() {
		return nil
	}
	if u.HaveMchNo != mchNo {
		return errors.New("商户用户没有此商户权限")
	}

	return nil
}

func (u User) ISHaveMchIDs(mchNos []string) error {
	if u.UserType.ISPlatform() {
		return nil
	}
	for _, id := range mchNos {
		if id == u.HaveMchNo {
			return nil
		}

	}
	return errors.New("商户用户没有此商户权限")
}
