package request

import (
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lihongsheng/pay-gateway/enum"
)

// CustomClaims structure
type CustomClaims struct {
	BaseClaims
	BufferTime int64
	jwt.RegisteredClaims
}

type BaseClaims struct {
	UUID        uuid.UUID
	ID          uint
	Username    string
	NickName    string
	AuthorityId uint
	UserType    enum.UserType
	HaveMchNO   string
}
