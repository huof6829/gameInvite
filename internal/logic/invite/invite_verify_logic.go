package invite

import (
	"context"

	"github.com/Savvy-Gameing/backend/common/global"
	"github.com/Savvy-Gameing/backend/common/logicx"
	"github.com/Savvy-Gameing/backend/internal/svc"
	"github.com/Savvy-Gameing/backend/internal/types"

	xerrors "github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type InviteVerifyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInviteVerifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InviteVerifyLogic {
	return &InviteVerifyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InviteVerifyLogic) InviteVerify(req *types.InviteVerifyReq) (resp *types.InviteVerifyResp, err error) {
	// todo: add your logic here and delete this line
	defer func() {
		if err != nil {
			l.Logger.Error(xerrors.WithStack(err))
		}
	}()

	/// debugclose
	// if valueJSON, ok := l.ctx.Value(jwt.BUILD_TOKEN_KEY).(string); !ok || valueJSON != req.Wallet {
	// 	l.Logger.Errorf("req.Wallet=%v, walletToken=%v", req.Wallet, valueJSON)
	// 	return nil, response.WalletIncorrect
	// }

	// userBind, err := l.svcCtx.UserBindModel.FindOneByWallet(l.ctx, req.Wallet)
	// if err == gorm.ErrRecordNotFound {
	// 	return nil, response.WalletNotConnected
	// } else if err != nil {
	// 	return nil, err
	// }

	if err = logicx.VerifyInviteCode(l.ctx, l.svcCtx, req.Id, req.InviteCode); err != nil {
		return nil, err
	}

	resp = &types.InviteVerifyResp{
		InviteCredit:    global.Invite_Credit_Direct_1,
		IsInviteSuccess: true,
	}
	return resp, nil

}
