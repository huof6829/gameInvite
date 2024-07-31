package job

import (
	"context"

	"github.com/Savvy-Gameing/backend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type JobConsumeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJobConsumeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JobConsumeLogic {
	return &JobConsumeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JobConsumeLogic) Start() {
	l.Logger.Infof("start consumer \n")

	// threading.GoSafe(func() {
	// 	l.svcCtx.Consumer.Consume(func(body []byte) {
	// 		l.Logger.Infof("consumer job  %s \n", string(body))
	// 	})
	// })
}

func (l *JobConsumeLogic) Stop() {
	l.Logger.Infof("stop consumer \n")
}
