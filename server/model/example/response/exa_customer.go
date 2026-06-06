package response

import "github.com/lihongsheng/pay-gateway/model/example"

type ExaCustomerResponse struct {
	Customer example.ExaCustomer `json:"customer"`
}
