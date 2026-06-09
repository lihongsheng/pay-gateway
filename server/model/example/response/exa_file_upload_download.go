package response

import "github.com/lihongsheng/pay-gateway/model/example"

type ExaFileResponse struct {
	File example.ExaFileUploadAndDownload `json:"file"`
}
