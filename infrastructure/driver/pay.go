package driver

import (
	"context"

	"github.com/lihongsheng/pay-gateway/infrastructure/driver/dto"
)

type Pay interface {
	Pay(ctx context.Context, req *dto.PayOrder) (*dto.PayResponse, error)
	Query(ctx context.Context, req dto.Query) (*dto.PayDetail, error)
	Close()
}
