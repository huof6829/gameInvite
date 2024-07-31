package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string

		// AccessExpire: 604800   ## 7天
		// RefreshAfter: 86400   ## 客户端提前1天刷新
		// AccessExpire int64
		// RefreshAfter int64
	}
	Mysql struct {
		DataSource   string
		MaxOpenConns int `json:",default=10"`
		MaxIdleConns int `json:",default=100"`
		MaxLifetime  int `json:",default=3600"`
	}
	CacheRedis cache.CacheConf
	BizRedis   redis.RedisConf
}
