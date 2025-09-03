package config

type Config struct {
	AppID          string
	MchID          string
	APIKey         string
	CertPrivateKey string
	//KeyPath string
	CertificateSerialNumber string
	PaymentName             string
	Proxy                   Proxy
}

type Proxy struct {
	Host     string
	Port     int
	UserName string
	Password string
}
