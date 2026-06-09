package system

import (
	"github.com/google/uuid"
	"github.com/lihongsheng/pay-gateway/enum"
	"github.com/lihongsheng/pay-gateway/global"
	"github.com/lihongsheng/pay-gateway/model/common"
)

type Login interface {
	GetUsername() string
	GetNickname() string
	GetUUID() uuid.UUID
	GetUserId() uint
	GetAuthorityId() uint
	GetUserInfo() any
	GetUserType() enum.UserType
	GetMchNo() string
}

var _ Login = new(SysUser)

type SysUser struct {
	global.GVA_MODEL
	UUID          uuid.UUID      `json:"uuid" gorm:"index;comment:用户UUID"`                                                                   // 用户UUID
	Username      string         `json:"userName" gorm:"index;comment:用户登录名"`                                                                // 用户登录名
	Password      string         `json:"-"  gorm:"comment:用户登录密码"`                                                                           // 用户登录密码
	NickName      string         `json:"nickName" gorm:"default:系统用户;comment:用户昵称"`                                                          // 用户昵称
	HeaderImg     string         `json:"headerImg" gorm:"default:https://qmplusimg.henrongyi.top/gva_header.jpg;comment:用户头像"`               // 用户头像
	AuthorityId   uint           `json:"authorityId" gorm:"default:888;comment:用户角色ID"`                                                      // 用户角色ID
	Authority     SysAuthority   `json:"authority" gorm:"foreignKey:AuthorityId;references:AuthorityId;comment:用户角色"`                        // 用户角色
	Authorities   []SysAuthority `json:"authorities" gorm:"many2many:sys_user_authority;"`                                                   // 多用户角色
	Phone         string         `json:"phone"  gorm:"comment:用户手机号"`                                                                        // 用户手机号
	Email         string         `json:"email"  gorm:"comment:用户邮箱"`                                                                         // 用户邮箱
	Enable        int            `json:"enable" gorm:"default:1;comment:用户是否被冻结 1正常 2冻结"`                                                    //用户是否被冻结 1正常 2冻结
	OriginSetting common.JSONMap `json:"originSetting" form:"originSetting" gorm:"type:text;default:null;column:origin_setting;comment:配置;"` //配置
	UserType      enum.UserType  `json:"userType" gorm:"column:user_type;default:1;comment:用户类型 1系统用户 2商户用户"`
	MchNo         string         `json:"mchNo" gorm:"column:mch_no;default:0;comment:商户编号"`
}

type SysMch struct {
	ID     uint   `gorm:"primarykey" json:"ID"`                                                          // 主键ID
	MchNo  string `gorm:"column:mch_no;type:varchar(100);not null;comment:编号" json:"mch_no"`             // 编号
	Status int64  `gorm:"column:status;type:tinyint;not null;default:1;comment:2 停用 1 正常" json:"status"` // 2 停用 1 正常
}

func (*SysMch) TableName() string {
	return "merchant"
}

func (SysUser) TableName() string {
	return "sys_users"
}

func (s *SysUser) GetUsername() string {
	return s.Username
}

func (s *SysUser) GetNickname() string {
	return s.NickName
}

func (s *SysUser) GetUUID() uuid.UUID {
	return s.UUID
}

func (s *SysUser) GetUserId() uint {
	return s.ID
}

func (s *SysUser) GetAuthorityId() uint {
	return s.AuthorityId
}

func (s *SysUser) GetUserInfo() any {
	return *s
}
func (s *SysUser) GetUserType() enum.UserType {
	return s.UserType
}

func (s *SysUser) GetMchNo() string {
	return s.MchNo
}
