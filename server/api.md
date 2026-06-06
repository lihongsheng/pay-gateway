### 支付系统
### 接口签名
#### 接口签名参数
接口签名，使用X-App-No、X-Method、X-Sign、X-Timestamp、X-Nonce进行签名，签名结果放在X-Sign中。
1. X-App-No：应用编号；
2. X-Method：签名方法，目前仅支持MD5签名方式；
3. X-Timestamp：时间戳；
4. X-Nonce：随机数;
5. X-Sign：签名结果。

#### 签名算法
签名算法：
以下方式进行签名：
生成X-App-No、X-Method、X-Timestamp、X-Nonce 对应的值，把对应的值按顺序进行拼接，然后进行MD5加密，得到签名结果。
拼接顺序：
str = X-Timestamp\n X-Nonce\n X-App-No\n X-Method\n
// key 是应用的密钥
X-Sign = md5(str + key)
```go
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
```
### 接口相关
#### 接口支付
1. 支付单;接口支付，需要商户自身维护收银台页面进行支付；在应用配置多微信账户支付情况下，在支付的时候需要传入订单以及微信账户信息（account_no），微信账户对应的openid；在应用配置单微信账户支付情况下，在支付的时候需要传入订单信息，支付会查询对应的微信单账户信息进行支付。
```shell
curl --location 'http://localhost:8888/public/v1/payment/pay' \
--header 'X-App-No: A4966c7e946c00' \
--header 'X-Method: MD5' \
--header 'X-Sign: 9499d21b5057af08b56ae7f7c503e7a6' \
--header 'X-Timestamp: 1772159912' \
--header 'X-Nonce: 8IeIbE8C4y1v' \
--header 'Content-Type: application/json' \
--data '{
    "order": {
        "order_no": "T1772159914",
        "subject": "测试支付",
        "amount": {
            "total":1,
            "currency": "CNY"
        },
        "create_at": "2026-03-05T16:43:46Z"
    },
    "pass_back_params": "order",
    "app_no": "A4966c7e946c00",
    "request_id": "191772159912",
    "payer": {
       "open_id":"ogdvH6h9jPp5R3f1fyLsQjdB-fAc" 
    },
    "account_no": "P4966c7e946c00",
    "payment_method": 2,
    "payment_product": 3
    }'
```
2. 订单查询
```shell
curl --location 'http://localhost:8888/public/v1/payment/query?app_no=A4966c7e946c00&order_no=PT177268013213600783993' \
--header 'X-App-No: A4966c7e946c00' \
--header 'X-Method: MD5' \
--header 'X-Sign: 9499d21b5057af08b56ae7f7c503e7a6' \
--header 'X-Timestamp: 1772159912' \
--header 'X-Nonce: 8IeIbE8C4y1v'
```
3. 订单关闭
```shell
curl --location 'http://localhost:8888/public/v1/payment/close' \
--header 'X-App-No: A4966c7e946c00' \
--header 'X-Method: MD5' \
--header 'X-Sign: 9499d21b5057af08b56ae7f7c503e7a6' \
--header 'X-Timestamp: 1772159912' \
--header 'X-Nonce: 8IeIbE8C4y1v' \
--header 'Content-Type: application/json' \
--data '{
   "order_no": "T1772159914",
    "app_no": "A4966c7e946c00"
    
    }'
```
#### 收银台支付
1. 收银台支付单
    1. 请求参数；创建收银台支付单
        ```shell
        curl --location 'https://proxy.test.jianxindianzi.com/sslab/public/v1/cashier/create' \
        --header 'X-App-No: A4966c7e946c00' \
        --header 'X-Method: MD5' \
        --header 'X-Sign: 9499d21b5057af08b56ae7f7c503e7a6' \
        --header 'X-Timestamp: 1772159912' \
        --header 'X-Nonce: 8IeIbE8C4y1v' \
        --header 'Content-Type: application/json' \
        --data '{
            "order": {
                "order_no": "T17721599111",
                "subject": "测试支付",
                "amount": {
                    "total":1,
                    "currency": "CNY"
                },
                "create_at": "2026-03-05T10:43:46Z"
            },
            "pass_back_params": "order",
            "app_no": "A4966c7e946c00",
            "request_id": "91772159912"
        }'
        ```
    2. 响应参数；基于返回的open_url，跳转到收银台页面,收银台页面自动支付，如果应用配置的是多微信支持，会在失败支付的时候自动使用其他微信支付账户进行支付。
         ```json
         {
             "code": 0,
             "data": {
                 "order_no": "T17721599112",
                 "open_url": "https://proxy.test.jianxindianzi.com/sslab/payment/web/public/cashier.html?action=hub&action_token=5KAW2O244VWIVWOFJXTORNLCQNO7NPO3SHVSWSKNVMUZ7QNJM4CZEXRY367NJ33W56AEBBAA7E5RMVWR&order_no=T17721599112",
                 "amount": {
                     "total": 1,
                     "currency": "CNY"
                 },
                 "trade_no": "195925959338291208250204",
                 "order_expire_time": "2026-03-12T17:00:50+08:00",
                 "status": 1
             },
             "msg": "成功"
         }
         ```
#### 退款
1. 退款单
```shell
curl --location 'http://localhost:8888/public/v1/refund' \
--header 'X-App-No: A4966c7e946c00' \
--header 'X-Method: MD5' \
--header 'X-Sign: 9499d21b5057af08b56ae7f7c503e7a6' \
--header 'X-Timestamp: 1772159912' \
--header 'X-Nonce: 8IeIbE8C4y1v' \
--header 'Content-Type: application/json' \
--data '{
   "order_no": "PT177270249285809578808",
    "app_no": "A4966c7e946c00",
    "refund_no": "R177270249285809578809",
    "amount": {
        "total":1
    },
    "reason": "接口测试",
    "pass_back_params": "refund"
    }'
```
2. 退款查询
```shell
curl --location 'http://localhost:8888/public/v1/refund?refund_no=R177270249285809578809&app_no=A4966c7e946c00' \
--header 'X-App-No: A4966c7e946c00' \
--header 'X-Method: MD5' \
--header 'X-Sign: 9499d21b5057af08b56ae7f7c503e7a6' \
--header 'X-Timestamp: 1772159912' \
--header 'X-Nonce: 8IeIbE8C4y1v'
```