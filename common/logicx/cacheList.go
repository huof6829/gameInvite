package logicx

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/redis/redistest"

	"github.com/Savvy-Gameing/backend/internal/model/sys_invite"
	"github.com/Savvy-Gameing/backend/internal/model/user_invite"
	"github.com/Savvy-Gameing/backend/internal/model/user_invite_count"
	"github.com/Savvy-Gameing/backend/internal/svc"
	"github.com/Savvy-Gameing/backend/pkg/orm"
)

func CacheListDefault(ctx context.Context, svcCtx *svc.ServiceContext, bizKey string, start, cursor, page, pageSize int64) ([]redis.Pair, error) {
	// b, err := svcCtx.BizRedis.ExistsCtx(ctx, bizKey)
	// if err != nil {
	// 	logc.Errorf(ctx, "[CacheListDefault] BizRedis.ExistsCtx error: %v", err)
	// }
	// if b {
	// 	err = svcCtx.BizRedis.ExpireCtx(ctx, bizKey, common.CacheListExpireTime)
	// 	if err != nil {
	// 		logc.Errorf(ctx, "[CacheListDefault] BizRedis.ExpireCtx error: %v", err)
	// 	}
	// }

	/// score 介于 [start, cursor]
	pairs, err := svcCtx.BizRedis.ZrangebyscoreWithScoresAndLimitCtx(ctx, bizKey, start, cursor, int(page), int(pageSize))
	if err != nil {
		logc.Errorf(ctx, "[CacheListDefault] BizRedis.ZrevrangebyscoreWithScoresAndLimitCtx error: %v", err)
		return nil, err
	}

	// var keys []string
	// var scores []int64
	// for _, pair := range pairs {

	// 	logc.Debugf(ctx, "pair.Score=%T, pair.Score=%v, pair.Key=%v", pair.Score, pair.Score, pair.Key)

	// 	scores = append(scores, pair.Score)
	// 	keys = append(keys, pair.Key)

	// 	// createTime, err := strconv.ParseInt(strconv.FormatInt(pair.Score, 10), 10, 64)
	// 	// if err != nil {
	// 	// logc.Errorf(ctx, "[CacheListDefault] strconv.ParseInt error: %v", err)
	// 	// continue
	// 	// }

	// }
	return pairs, nil
}

/***************************************************8***************/

func BizGet(ctx context.Context, svcCtx *svc.ServiceContext, bizKey string,
	model interface{},
	seconds int,
) (isCache bool,
	bizrlt string,
	err error,
) {

	if isExist, err := CheckBizKey(ctx, svcCtx, bizKey, seconds); !isExist || err != nil {
		return false, "", err
	}
	bizrlt, err = svcCtx.BizRedis.GetCtx(ctx, bizKey)
	if err != nil {
		logc.Errorf(ctx, "[BizGet] BizRedis.GetCtx bizKey=%v, err=%v", bizKey, err)
		return false, "", err
	}

	return true, bizrlt, nil
}

func BizSet(ctx context.Context, svcCtx *svc.ServiceContext, bizKey string,
	model interface{},
	seconds int,
) error {

	val, err := json.Marshal(model)
	if err != nil {
		logc.Errorf(ctx, "[BizSet] json.Marshal bizKey=%v, err=%v", bizKey, err)
		return err
	}
	if err = svcCtx.BizRedis.SetCtx(ctx, bizKey, string(val)); err != nil {
		logc.Errorf(ctx, "[BizSet] BizRedis.SetCtx bizKey=%v, err=%v", bizKey, err)
	}
	if seconds > 0 {
		return svcCtx.BizRedis.ExpireCtx(ctx, bizKey, seconds)
	}
	return err
}

// 分页 pgae从 0 开始
func BizZrevrange(ctx context.Context, svcCtx *svc.ServiceContext, bizKey string,
	start, stop float64,
	page, size int,
	cb func([]string) []interface{},
	seconds int,
) (isCache bool,
	isEnd bool,
	count int,
	rlts []interface{},
	err error,
) {

	if isExist, err := CheckBizKey(ctx, svcCtx, bizKey, seconds); !isExist || err != nil {
		return isCache, isEnd, count, nil, err
	}

	count, err = svcCtx.BizRedis.ZcountCtx(ctx, bizKey, -999999999, 999999999)
	if err != nil {
		logc.Errorf(ctx, "[GameReviewList] BizRedis.ZcountCtx bizKey=%s, err=%v", bizKey, err)
		return isCache, isEnd, count, nil, err
	}
	count -= 1 // 去掉最后一条

	var pairs []redis.FloatPair

	if size > 0 {
		if count <= (page+1)*size {
			isEnd = true
		}
		if count/size < page {
			return true, true, count, nil, nil
		}

		pairs, err = svcCtx.BizRedis.ZrevrangebyscoreWithScoresByFloatAndLimitCtx(ctx, bizKey, start, stop, page, size)

	} else {
		pairs, err = svcCtx.BizRedis.ZrevrangebyscoreWithScoresByFloatCtx(ctx, bizKey, start, stop)
	}
	if err != nil {
		logc.Errorf(ctx, "[BizZrevrange] BizRedis.ZrevrangebyscoreWithScoresByFloatAndLimitCtx bizKey=%v, start=%v, stop=%v, page=%v, size=%v, err=%v", bizKey, start, stop, page, size, err)
		return isCache, isEnd, count, nil, err
	}

	var (
		keys   []string
		scores []float64
	)
	for _, pair := range pairs {
		scores = append(scores, pair.Score)
		keys = append(keys, pair.Key)
	}

	if len(scores) > 0 {
		isCache = true
		if scores[len(scores)-1] == 0 { // 去掉最后一条
			scores = scores[:len(scores)-1]
			keys = keys[:len(keys)-1]
			isEnd = true
		}
		if len(scores) == 0 {
			return isCache, isEnd, count, nil, nil
		}

		return isCache, isEnd, count, cb(keys), nil
	}

	return isCache, isEnd, count, nil, nil
}

func BizZadd(ctx context.Context, svcCtx *svc.ServiceContext, bizKey string,
	datas []interface{},
	cb func(interface{}) int64,
	seconds int,
) error {

	for _, data := range datas {
		score := cb(data)
		bytes, err := json.Marshal(data)
		if err != nil {
			logc.Errorf(ctx, "[BizZadd] json.Marshal bizKey=%v, err=%v", bizKey, err)
			continue
		}

		_, err = svcCtx.BizRedis.ZaddCtx(ctx, bizKey, score, string(bytes))
		if err != nil {
			logc.Errorf(ctx, "[BizZadd] BizRedis.ZaddCtx bizKey=%v, score=%v, err=%v", bizKey, score, err)
		}
	}

	if seconds > 0 {
		return svcCtx.BizRedis.ExpireCtx(ctx, bizKey, seconds)
	}
	return nil
}

func CheckBizKey(ctx context.Context, svcCtx *svc.ServiceContext, bizKey string, seconds int) (isExist bool, err error) {
	b, err := svcCtx.BizRedis.ExistsCtx(ctx, bizKey)
	if err != nil {
		logc.Errorf(ctx, "[CheckBizKey] BizRedis.ExistsCtx bizKey=%s, err=%v", bizKey, err)
		return false, err
	}
	if !b {
		return false, nil
	}
	if seconds > 0 {
		if err = svcCtx.BizRedis.ExpireCtx(ctx, bizKey, seconds); err != nil {
			logc.Errorf(ctx, "[CheckBizKey] BizRedis.ExpireCtx bizKey=%v, err=%v", bizKey, err)
			return true, err
		}
	}
	return true, nil
}

/***************************************************8***************/

func TestConfig(t *testing.T) (ctx context.Context, svcCtx *svc.ServiceContext) {
	ctx = context.Background()

	r := redistest.CreateRedis(t)
	confCache := cache.ClusterConf{
		{
			RedisConf: redis.RedisConf{
				Host: r.Addr,
				Type: redis.NodeType,
			},
			Weight: 100,
		},
	}
	db := orm.MustNewMysql(&orm.Config{
		DSN:          "root:root123456@tcp(127.0.0.1:3306)/savvy_gameing?timeout=1s&readTimeout=1s&writeTimeout=1s&charset=utf8mb4&parseTime=true&loc=Local",
		MaxIdleConns: 10,
		MaxOpenConns: 100,
		MaxLifetime:  3600,
	})

	svcCtx = &svc.ServiceContext{
		DB: db,

		BizRedis: redis.MustNewRedis(redis.RedisConf{
			Host: "127.0.0.1:6379",
			Pass: "",
			Type: "node",
		}),

		UserInviteModel:      user_invite.NewUserInviteModel(db.DB, confCache),
		UserInviteCountModel: user_invite_count.NewUserInviteCountModel(db.DB, confCache),
		SysInviteModel:       sys_invite.NewSysInviteModel(db.DB, confCache),
	}

	return
}
