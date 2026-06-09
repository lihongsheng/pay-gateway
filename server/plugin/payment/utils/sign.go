package utils

import (
	"crypto/md5"
	"fmt"
	"github.com/lihongsheng/pay-gateway/plugin/payment/enum"
	"github.com/lihongsheng/pay-gateway/plugin/payment/errors"
	"io"
	"net/http"
	"time"
)

func GenMd5Meta(key string, appNo string) map[string]string {
	meta := map[string]string{
		enum.MetaTimestamp:  fmt.Sprintf("%d", time.Now().Unix()),
		enum.MetaNonce:      RandomString(12),
		enum.MetaAppNo:      appNo,
		enum.MetaSignMethod: "MD5",
	}
	str := fmt.Sprintf("%s\n%s\n%s\n%s\n", meta[enum.MetaTimestamp], meta[enum.MetaNonce], meta[enum.MetaAppNo], meta[enum.MetaSignMethod])
	meta[enum.MetaSign] = Md5(str + key)
	return meta
}

func SignValidate(key string, req *http.Request) error {
	var meta = map[string]string{}
	meta[enum.MetaTimestamp] = req.Header.Get(enum.MetaTimestamp)
	meta[enum.MetaNonce] = req.Header.Get(enum.MetaNonce)
	meta[enum.MetaAppNo] = req.Header.Get(enum.MetaAppNo)
	meta[enum.MetaSignMethod] = req.Header.Get(enum.MetaSignMethod)
	meta[enum.MetaSign] = req.Header.Get(enum.MetaSign)
	str := fmt.Sprintf("%s\n%s\n%s\n%s\n", meta[enum.MetaTimestamp], meta[enum.MetaNonce], meta[enum.MetaAppNo], meta[enum.MetaSignMethod])
	sign := Md5(str + key)
	if !(sign == meta[enum.MetaSign]) {
		return errors.NewError(errors.ErrCodeSignValidateFailed, "签名验证失败")
	}
	return nil
}

func RandomString(length int) string {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(result)
}

func Md5(content string) (md string) {
	h := md5.New()
	_, _ = io.WriteString(h, content)
	md = fmt.Sprintf("%x", h.Sum(nil))
	return
}
