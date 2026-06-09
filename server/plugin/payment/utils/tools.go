package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/VictoriaMetrics/fastcache"
	"net/http"
	"sync"
	"time"
)

func HttpJsonPost(ctx context.Context, url string, data interface{}, expireTime time.Duration) (*http.Response, error) {
	body, _ := json.Marshal(data)
	reqCtx, cancel := context.WithTimeout(ctx, expireTime)
	defer cancel()
	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

var fastCache *fastcache.Cache
var fastCacheSync = sync.Once{}

func GetBigCache() *fastcache.Cache {
	fastCacheSync.Do(func() {
		fastCache = fastcache.New(1024 * 1024 * 100) // 100 MB
	})
	return fastCache
}
