package request

import (
	"github.com/lihongsheng/pay-gateway/model/common/request"
)

type ExaAttachmentCategorySearch struct {
	ClassId int `json:"classId" form:"classId"`
	request.PageInfo
}
