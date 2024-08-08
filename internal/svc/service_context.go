package svc

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"

	"github.com/Savvy-Gameing/backend/internal/config"
	"github.com/Savvy-Gameing/backend/internal/model/desc"
	"github.com/Savvy-Gameing/backend/internal/model/sys_invite"
	"github.com/Savvy-Gameing/backend/internal/model/user_invite"
	"github.com/Savvy-Gameing/backend/internal/model/user_invite_count"
	"github.com/Savvy-Gameing/backend/pkg/orm"
)

type ServiceContext struct {
	Config   config.Config
	BizRedis *redis.Redis
	DB       *orm.DB
	// Consumer dq.Consumer
	Cron  *cron.Cron
	TgBot *tgbotapi.BotAPI

	UserInviteModel      user_invite.UserInviteModel
	UserInviteCountModel user_invite_count.UserInviteCountModel
	SysInviteModel       sys_invite.SysInviteModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := orm.MustNewMysql(&orm.Config{
		DSN:          c.Mysql.DataSource,
		MaxOpenConns: c.Mysql.MaxOpenConns,
		MaxIdleConns: c.Mysql.MaxIdleConns,
		MaxLifetime:  c.Mysql.MaxLifetime,
	})

	rds := redis.MustNewRedis(redis.RedisConf{
		Host: c.BizRedis.Host,
		Pass: c.BizRedis.Pass,
		Type: c.BizRedis.Type,
	})

	initTable(db)

	return &ServiceContext{
		Config:   c,
		DB:       db,
		BizRedis: rds,
		// Consumer: dq.NewConsumer(c.DqConf),
		Cron: cron.New(cron.WithSeconds()),

		UserInviteModel:      user_invite.NewUserInviteModel(db.DB, c.CacheRedis),
		UserInviteCountModel: user_invite_count.NewUserInviteCountModel(db.DB, c.CacheRedis),
		SysInviteModel:       sys_invite.NewSysInviteModel(db.DB, c.CacheRedis),
	}
}

func initTable(db *orm.DB) error {
	// debugclose
	// if err := db.Migrator().DropTable(&desc.SysInvite{}); err != nil {
	// 	logx.Errorf("[initTable] DropTable err=%v", err)
	// }
	// if err := db.Migrator().DropTable(&desc.UserInvite{}); err != nil {
	// 	logx.Errorf("[initTable] DropTable err=%v", err)
	// }
	// if err := db.Migrator().DropTable(&desc.UserInviteCount{}); err != nil {
	// 	logx.Errorf("[initTable] DropTable err=%v", err)
	// }

	for _, v := range []interface{}{
		&desc.SysInvite{},
		&desc.UserInvite{},
		&desc.UserInviteCount{},
	} {
		if !db.Migrator().HasTable(v) {
			if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;").Migrator().CreateTable(v); err != nil {
				logx.Errorf("[initTable] CreateTable v=%T, err=%v", v, err)
			}
		} else {

			// switch v.(type) {
			// case *desc.GameReviewLike:
			// 	if !db.Migrator().HasIndex(v, "idx_game_review_like_user_id") {
			// 		db.Migrator().CreateIndex(v, "idx_game_review_like_user_id")
			// 	}
			// }

		}
	}
	return nil
}
