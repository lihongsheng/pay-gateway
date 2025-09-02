package driver

import (
	"context"

	"github.com/lihongsheng/pay-gateway/driver/dto"
)

type Pay interface {
	Pay(ctx context.Context, req *dto.PayOrder) (*dto.PayResponse, error)
	Query()
	Close()
}
