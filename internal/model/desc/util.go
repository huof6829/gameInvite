package desc

import (
	"context"
	"testing"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/redis/redistest"

	"github.com/Savvy-Gameing/backend/pkg/orm"
)

func TestConfig(t *testing.T) (ctx context.Context, db *orm.DB, confCache cache.ClusterConf) {
	ctx = context.Background()

	r := redistest.CreateRedis(t)
	confCache = cache.ClusterConf{
		{
			RedisConf: redis.RedisConf{
				Host: r.Addr,
				Type: redis.NodeType,
			},
			Weight: 100,
		},
	}
	db = orm.MustNewMysql(&orm.Config{
		DSN:          "root:root123456@tcp(127.0.0.1:3306)/savvy_gameing?timeout=1s&readTimeout=1s&writeTimeout=1s&charset=utf8mb4&parseTime=true&loc=Local",
		MaxIdleConns: 10,
		MaxOpenConns: 100,
		MaxLifetime:  3600,
	})

	return
}
