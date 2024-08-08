package handler

import (
	"context"

	"github.com/zeromicro/go-zero/core/service"

	// job "github.com/Savvy-Gameing/backend/internal/logic/job"
	cron "github.com/Savvy-Gameing/backend/internal/logic/cron"
	"github.com/Savvy-Gameing/backend/internal/logic/tgbot"
	"github.com/Savvy-Gameing/backend/internal/svc"
)

// func JobRun(serverCtx *svc.ServiceContext) {
// 	threading.GoSafe(func() {
// 		job.JobProduceHandler(serverCtx)
// 		job.JobConsumeHandler(serverCtx)
// 		//...many job
// 	})
// }

func RegisterJob(serverCtx *svc.ServiceContext, group *service.ServiceGroup) {
	// group.Add(job.NewJobProduceLogic(context.Background(), serverCtx))
	// group.Add(job.NewJobConsumeLogic(context.Background(), serverCtx))

	group.Add(cron.NewCronLogic(context.Background(), serverCtx))
	group.Add(tgbot.NewTgbotLogic(context.Background(), serverCtx))
}
