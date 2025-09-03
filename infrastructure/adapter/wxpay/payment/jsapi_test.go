package payment

import (
	"context"
	"time"

	"github.com/lihongsheng/pay-gateway/infrastructure/config"
	"github.com/lihongsheng/pay-gateway/infrastructure/driver/dto"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestJsapi_Pay(t *testing.T) {
	api, err := NewJsApi(config.Config{
		AppID: "wx2f0b0c0c0c0c0c0c",
		MchID: "123456789",
		CertPrivateKey: `-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQDCsnMIAnRpqu3z
ZjYHmqixlG/ynV0ZfagsWvWwCsnF69bdrW/txUDRThZUUl2J8NUWmy6JyH2ex+FS
nsI05mdr6s296CXNiWu0yGrKdhT5spQEECZhpXIMc6nGSnzp+GnSruqhfxih5K9h
FpVo+gUdFh5qHLR3cv6AvylGGvSTvRFkytoLfPNdu5RwzsY2JuIpsWJCWWGUaV9t
XRuY4fwtlzBbC/jQ+113jHaCKGfm+bMZ1DU4B/cPpia2xd/uIyPQfWDmy/IezXud
xI3uH2JZOtUpw0PGuXV2bgFGzEMmxZ6fCsgzVZQe6ieVsTo6ggkQ7blPe5zKbNYP
MxJbF7exAgMBAAECggEBAIdvhS98LajX26DmaA1QG6s0G8/Egd/almMMfz4Psx54
GUapgGQBRD6VOFk91o2/Nyv7lRsJmcEbP/WuNGCCKk1az/YcCDf7MS5YAFmIXLz5
6ZcN+PUSFszspJwocs57HHoPbW4cMHFl2E4MXLDiwy3hlhSwlSVGnB3JXJfE5n/h
2/cHSZvcErnwJhdsExp7e4RqV4h7nJBAvU1ZnJVUbzhyjyp0wY0/JLl5G2fq0qqn
wnoeQCmqioK6eXH8zrWCnhe1W4ipcPaAkMXST7YcyPLANkxvSXG2CXxb82t671Wr
NqEdHEkly4S532S5ayuuJChwIRR7woWSZ1h43SZWkQUCgYEA+x2SMIlS4fnqWzQQ
yJs5dGfgcqh5J5Kn8J7NOry74LNLf7COL+u6fUcm4rwrBr+Q3YXWrBvyqDw4LUEB
VYhtr8oHJRbYaMjFr9A14QScBhiLjCqf8tBIFQVLqPd8OqxtsN6zSAMxGj7X68xD
t7ztPLbVv8VBk0PbShCKM5o6qpMCgYEAxnvxXCZnsBvcxMbIyq0UJwgmUueUsZfM
D0WRyEgfGKMR+6AwAvVpHaPyEzErXzPlT1TnKvwpUfagfsFsJoUrhN6czey1Np8m
reBaAAXuCRuzavRy1iMjgFuKYH0wORTR3A+fD57ZdNL5VCBN9RsQJwU3JgQZMVkf
pQ6T/oYsSysCgYEApqlxpRT/FUuw5ucfXITpFQD8ThzSjBkhrOk4fItWhkN5ED41
oEhrdUoL3N/WDpyFoQB7Aa9q1Y1iG2bRY9swMUN8inknGCRoT894cueERed0doqz
rYvey1TAalwW7zoRcxnbEyhLJoge9jiTmRaivXD7XFOmuf6HRBjGIIlz9lECgYBs
E/VbPjZbuPA/3hZb9l7w2gk0P5HCGmwtLK6zJkJ4geM65wD9u3AfibQ5Kx742iNV
TWALEf/V97txChW/6+fElAtCPlB2i7beGzompRP2tbS+2pjlbYDZVf9FhyWJD4Mu
lvr/4Hl8mZzWaDjK7I+hD7/13Wlya5tFn2iKwbjAvQKBgQCLiRluLqvhELCyP6VL
3Os3qdb8rLbgMz9ttaxsuI7wDW8g40AkU4iwxzSQRBiCGOMtQHWS3xLxHLSAjUvG
MSdanpimagZ8aHsWo5FtL/oqucHrd1qd7SUb/r/A7HhTctwrV07j+ru2kExnNIGW
/p/L0H4Yz1blO0d6L7EF+3emKA==
-----END PRIVATE KEY-----`,
		CertificateSerialNumber: "123456789",
		APIKey:                  "123456789",
		Proxy:                   config.Proxy{Host: "127.0.0.1", Port: 7897},
	})

	assert.NoError(t, err)
	ctx := context.Background()
	req := &dto.PayOrder{
		Order: dto.Order{
			OrderNo: "123456789",
			PayAmount: dto.Amount{
				Currency: "CNY",
				Total:    1,
			},
			Subject: "test pay",
		},
		NotifyUrl:      "https://www.baidu.com",
		PassbackParams: "123456789",
		Payer: dto.Payer{
			OpenID: "123456789",
		},
		RedirectUrl: "https://www.baidu.com",
		TimeExpire:  time.Now().Unix() + 90,
	}
	resp, err := api.Pay(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
