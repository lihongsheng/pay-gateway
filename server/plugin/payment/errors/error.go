package errors

import (
	"errors"
	"fmt"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func IsUserLimitError(err error) bool {
	if err == nil {
		return false
	}
	userErr, errOk := err.(*Error)
	if errOk && userErr.Code == ErrUserLimit {
		return true
	}
	return false
}

func (e *Error) Error() string {
	if e.Err == nil {
		// 无包装错误，直接返回错误信息
		return fmt.Sprintf("[%d] %s", e.Code, e.Message)
	}
	// 有包装错误，拼接【当前错误+原始错误链】，完整展示错误链路
	return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) IsCode(code int) bool {
	return e.Code == code
}

func WrapCode(code int, err error) error {
	return &Error{
		Code:    code,
		Message: GetErrMsg(code),
		Err:     err, // 核心：把原始错误挂载到当前错误的 Err 字段
	}
}

func NewError(code int, msg string) error {
	return &Error{
		Code:    code,
		Message: msg,
		Err:     nil,
	}
}

func NewErrorDefault(code int) error {
	return &Error{
		Code:    code,
		Message: GetErrMsg(code),
		Err:     nil,
	}
}

func WrapError(code int, msg string, err error) error {
	return &Error{
		Code:    code,
		Message: msg,
		Err:     err,
	}
}

func ParseError(err error) *Error {
	if err == nil {
		return nil
	}
	// 如果是自定义的Error，直接类型断言返回
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr
	}
	// 兜底：非自定义错误，统一转为【系统内部错误】
	return &Error{
		Code:    50000,
		Message: "系统内部错误",
		Err:     err,
	}
}
