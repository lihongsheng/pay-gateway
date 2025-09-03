package tools

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// isTimeActive 判断时间是否在一个范围内
func isTimeActive(now, start, stop time.Time) bool {
	return (start.Before(now) || start.Equal(now)) && now.Before(stop)
}

func UnixToTime(t int64) time.Time {
	return time.Unix(t, 0)
}

func EndTime(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d+1, 0, 0, 0, 0, time.Local)
}

func StartTime(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}

func Md5(content string) (md string) {
	h := md5.New()
	_, _ = io.WriteString(h, content)
	md = fmt.Sprintf("%x", h.Sum(nil))
	return
}

func VerifyMobileFormat(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

func IsURI(fl string) bool {
	if !strings.Contains(fl, "http://") && !strings.Contains(fl, "https://") {
		fl = "http://" + fl
	}
	s := fl
	_, err := url.ParseRequestURI(s)
	return err == nil
}

func GenerateID() string {
	return uuid.New().String()
}

const time33Hash = uint64(5381)

func Time33(str string) uint64 {
	hash := time33Hash
	for _, char := range str {
		tmp, _ := strconv.ParseInt(fmt.Sprintf("%d", char), 10, 0)
		hash = (hash << 5) + uint64(tmp)
	}
	return hash & 0x7FFFFFFF
}

func HmacSha256(key []byte, data []byte) []byte {
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write(data)
	return mac.Sum(nil)
}
