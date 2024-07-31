package cron

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/Savvy-Gameing/backend/internal/svc"
)

type JobCronLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCronLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JobCronLogic {
	return &JobCronLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JobCronLogic) Start() {
	l.Logger.Infof("start cron.")

	l.initCron()
	l.svcCtx.Cron.Start()

	// threading.GoSafe(func() {
	// 	l.svcCtx.Consumer.Consume(func(body []byte) {
	// 		l.Logger.Infof("consumer job  %s \n", string(body))
	// 	})
	// })
}

func (l *JobCronLogic) Stop() {
	l.Logger.Infof("stop cron.")
	l.svcCtx.Cron.Stop()
}

func (l *JobCronLogic) initCron() {
	l.Logger.Info("init cron.")

	// threading.RunSafe(func() {

	// 	// 每天5时3分3秒
	// 	if _, err := l.svcCtx.Cron.AddFunc("3 3 5 * * ?", func() { l.userLikeWriteDB() }); err != nil {
	// 		l.Logger.Errorf("[initCron] userLikeWriteDB error: %v", err)
	// 		return
	// 	}

	// })
}

func (l *JobCronLogic) userLikeWriteDB() {

	l.Logger.Infof("write user like record/count into DB begin.")
	begin := time.Now()

	delta := time.Now().Sub(begin)
	l.Logger.Infof("write user like record/count into DB end. cost %v ", delta)
}
