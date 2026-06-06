package enum

import "time"

const (
	ApplicationInfoCacheKeys   = "app:%s"
	ApplicationInfoCacheExpire = time.Hour * 24
)

const (
	AccountIDCacheKeys   = "account:%d"
	AccountIDCacheExpire = time.Hour * 24
	AccountNoCacheKeys   = "account:%s"
	AppAccounts          = "app:accounts:%s"
	AppAccountsVersion   = "app:accounts:version:%s"
)

const (
	// RouterStatisticCacheKeys 路由统计缓存, key: app_no, t
	RouterStatisticCacheKeys     = "router:%s%s"
	RouterStatisticLockCacheKeys = "routerLock:%s%s"
	RouterStatisticCacheExpire   = time.Minute * 5
)
