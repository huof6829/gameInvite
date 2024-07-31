package job

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/Savvy-Gameing/backend/internal/svc"
)

type JobProduceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJobProduceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JobProduceLogic {
	return &JobProduceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JobProduceLogic) Start() {
	l.Logger.Infof("start  Producer \n")

	// threading.GoSafe(func() {
	// 	producer := dq.NewProducer(l.svcCtx.Config.DqConf.Beanstalks)

	// 	// for i := 1000; i < 1005; i++ {
	// 	// 	_, err := producer.Delay([]byte(strconv.Itoa(i)), time.Second*1)
	// 	// 	if err != nil {
	// 	// 		l.Logger.Error(err)
	// 	// 	}
	// 	// }

	// 	str, err := producer.At([]byte("hello1111111111"), time.Now().Add(time.Second*10))
	// 	l.Logger.Debugf("JobProduce str: ", str)
	// 	if err != nil {
	// 		l.Logger.Error(err)
	// 	}
	// })
}

func (l *JobProduceLogic) Stop() {
	l.Logger.Infof("stop Producer \n")
}
