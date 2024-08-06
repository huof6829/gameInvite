package sys

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/Savvy-Gameing/backend/common/logicx"
	"github.com/Savvy-Gameing/backend/internal/svc"
)

type TelegramLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTelegramLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TelegramLogic {
	return &TelegramLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TelegramLogic) Telegram() error {
	// todo: add your logic here and delete this line

	err := logicx.TgBotStart(l.ctx, l.svcCtx)
	if err != nil {
		l.Logger.Errorf("err: %v", err)
	}

	return nil
}
