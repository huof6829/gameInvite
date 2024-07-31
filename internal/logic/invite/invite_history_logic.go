package invite

import (
	"context"

	"github.com/Savvy-Gameing/backend/common/global"
	"github.com/Savvy-Gameing/backend/internal/svc"
	"github.com/Savvy-Gameing/backend/internal/types"
	"gorm.io/gorm"

	xerrors "github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type InviteHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInviteHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InviteHistoryLogic {
	return &InviteHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InviteHistoryLogic) InviteHistory(req *types.InviteHistoryReq) (resp *types.InviteHistoryResp, err error) {
	// todo: add your logic here and delete this line
	defer func() {
		if err != nil {
			l.Logger.Error(xerrors.WithStack(err))
		}
	}()

	// if valueJSON, ok := l.ctx.Value(jwt.BUILD_TOKEN_KEY).(string); !ok || valueJSON != req.Wallet {
	// 	l.Logger.Errorf("req.Wallet=%v, walletToken=%v", req.Wallet, valueJSON)
	// 	return nil, response.WalletIncorrect
	// }

	// userBind, err := l.svcCtx.UserBindModel.FindOneByWallet(l.ctx, req.Wallet)
	// if err == user_bind.ErrNotFound {
	// 	return nil, response.WalletNotConnected
	// } else if err != nil {
	// 	return nil, err
	// }

	resp = &types.InviteHistoryResp{
		InviteHistorys: []types.InviteHistory{},
	}

	userInvites, err := l.svcCtx.UserInviteModel.FindByParentLevels(l.ctx, l.svcCtx.DB.DB, req.Id, []int64{global.Invite_Level_1}, global.DBFindLimit, "id desc")
	if err == gorm.ErrRecordNotFound {
		return resp, nil
	} else if err != nil {
		return nil, err
	}

	resp.TotalInviteCount = len(userInvites)
	resp.SuccessInviteCount = resp.TotalInviteCount // modifyfuture

	for _, userInvite := range userInvites {
		// modify
		var childId int64
		var childUserName string
		var avatar string

		// childBind, err := l.svcCtx.UserBindModel.FindOne(l.ctx, userInvite.ChildId)
		// if err != nil {
		// 	l.Logger.Errorf("UserBindModel.FindOne id=%v, err=%v", userInvite.ChildId, err)
		// 	continue
		// }

		resp.InviteHistorys = append(resp.InviteHistorys, types.InviteHistory{
			InviteCredit:  int(userInvite.InviteCreditDirectChild),
			CreateTime:    userInvite.CreatedAt.Unix(),
			ChildUserName: childUserName,
			ChildId:       childId,
			Avatar:        avatar,
		})
	}

	return resp, nil
}
