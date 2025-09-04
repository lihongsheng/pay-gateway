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
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCul+ogGYraBOK7
bl4nhGH0XHKjAhtm6Fs7Rh/N0B6TRPvprUgYj/OI4lkHPodI6fAE5idt2nPfN3lC
N6z1qTx2YdeZUr8P3nWt/QVZUjXdFJyQg0SdrSbRWaO5YOvIZTUXKI1bBY/Ez+DZ
0pa++RCfELwG64zCX6i/ZbrMMLnVg74ND3OJq6QCoS9GvuNuPijC9hfXIt45UYOV
gDceI1FrvC4hsXM+lK1+XeO1ZlDZZKY+gh2Qq0c1KXSUNrRqc+br8S7HDg5PSZcu
yo3iy3g+0g8hlqYFniN1mrMx+KD/eibNjWJJZ9306LKVRYaLrbP2ltIbyJMnWw6o
dCAAJkUvAgMBAAECggEAG+64CVysx0ai06PLsnzOgx7LRGMQ6TzhY7FLY00ZpywR
NYiAkVVufAbigXMyh6rNuKrtwfiCM0MXxk1MAZ2Ky9HJgYEoiixoHWbulSS+81jZ
15D4yFUsca/qrrzf3Ba9ELBvNRHFt5Ieghrb/T+xthEiU+UZhvtw7vYilYpyIMOH
9bKoHYC+DhnaWMFixLV1HRrMcnFtCBEUb9vTsTtCo/iTKGLh1ZE4wr5i8eIJn1kU
A0DH/V4bLedwmlggE7x4wVeGNVjLpHMPzC3MuOEqxl9Fce/ZoMZ/miswE0XDFNmd
p369sVVAoBs04WGI/SMPnjnjfjvkdzCLBszElDtAgQKBgQDr4H+XIa5arQr4xNjO
DaP52giylEoPASmwcKdapofe25FSO5fFswu6YGCbKbxapW+3rCeOtMS/1nAm+P+e
kG3bxC/O92uhZcx68NmFshcVyPhauM15ofNrqrcJZ314mfxammaWf78X0gXNBEYo
7TqQCozaQmXSoe1spCeKWiXj7wKBgQC9fP9AJBeqlh68PhP/NJl0gXSGEfHyODnh
VL6ZyUFCMJL9RrZ4k34/Xl7OzTiJsIidTsaXmKHubdJ2p5ODbuwgWlBQxik5uwIY
6mA/nzzgpxJUWjPO1kWNzjDJcnEKfxxqPaQFMAMEEgsiZXQ1+cOF8dAoiAXVLMxv
LLZefjtywQKBgFT5yb1umtt736oDcH+7FknarKt5FL0XFCfGTeQwfl5hB5dydJj6
ic0ZD3SWso6NbxPiZ1XX9dGOtVS0+5HqQkmcAgUzyYiZnfLkddjeccadYit8zDl3
iLRPCiXPLLqX7vvNCAxL9VHljcVwNr5jyGdABcJTjU82msw6zyvekivtAoGBAK9m
2nyMLGAa8D5+FFKjZP1ErBFC4Ty/MUzU/k4qzr6vahELePMqTK3p9EboDtLd18gN
2KURg6vKewyc2F7MucTE9R2gIR8wbaOUqR2bkGXAIaZ1jQmErQOz/tMqnVsDCNGL
V535sID/FtFzKlygY45EpdQu/X80JdUvhWz42pzBAoGAHP/7QWzB41Fs7ywgBlmj
ZPJCrMKZlQvfOrN9W0FFHXr6ungeOC0WFX4rQMTE4GK/pVwa9rmhQgnLiEMHesr1
AcKAFrovEWpLfSBsFfatC3nmqmOInQ00pDw0uI7x3hvdBc5DfyxseCaWP+WO452V
AzINQM9HQithUbCUjjpLWhI=
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
