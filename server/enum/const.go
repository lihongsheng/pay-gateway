package enum

// 用户类型

type UserType int

const (
	// 平台商户
	UserType_Platform UserType = 1
	// 商户
	UserType_Mch UserType = 2
)

func (u UserType) ISPlatform() bool {
	return u == UserType_Platform
}

func (u UserType) ISMch() bool {
	return u == UserType_Mch
}
