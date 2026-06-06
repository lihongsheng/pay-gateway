package request

import (
	"github.com/lihongsheng/pay-gateway/model/common/request"
	"github.com/lihongsheng/pay-gateway/model/system"
)

type SysOperationRecordSearch struct {
	system.SysOperationRecord
	request.PageInfo
}
