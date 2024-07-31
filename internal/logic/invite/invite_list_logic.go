package invite

import (
	"context"

	"github.com/Savvy-Gameing/backend/common/global"
	"github.com/Savvy-Gameing/backend/internal/model/user_invite"
	"github.com/Savvy-Gameing/backend/internal/model/user_invite_count"
	"github.com/Savvy-Gameing/backend/internal/svc"
	"github.com/Savvy-Gameing/backend/internal/types"
	"gorm.io/gorm"

	xerrors "github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type InviteListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInviteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InviteListLogic {
	return &InviteListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InviteListLogic) InviteList(req *types.InviteListReq) (resp *types.InviteListResp, err error) {
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

	resp = &types.InviteListResp{}

	// userBind, err := l.svcCtx.UserBindModel.FindOneByWallet(l.ctx, req.Wallet)
	// if err == user_bind.ErrNotFound {
	// 	return nil, response.WalletNotConnected
	// } else if err != nil {
	// 	return nil, err
	// }

	// logx.Debugf("userBind=%+v", userBind)

	userId := req.Id

	var userInvite *user_invite.UserInvite
	userInvites, err := l.svcCtx.UserInviteModel.FindByChildIdLevels(l.ctx, l.svcCtx.DB.DB, userId, []int64{global.Invite_Level_1})
	if err == gorm.ErrRecordNotFound {
		logx.Debugf("UserInviteModel child_id=%v, level=%v, err=ErrNotFound", userId, global.Invite_Level_1)
		return resp, nil
	} else if err != nil {
		return nil, err
	} else if len(userInvites) == 0 {
		return resp, nil
	} else {
		userInvite = userInvites[0]
	}

	/// modifyfuture   完成所有任务
	/// debugclose
	// if err = logicx.IndirectInviteCreditSum(l.ctx, l.svcCtx, userId); err != nil {
	// 	return resp, err
	// }

	var (
		inviteCodeParent   string
		inviteCodeSelf     string
		totalCredit        int64
		totalInviteCount   int64
		successInviteCount int64

		userInviteCountParent *user_invite_count.UserInviteCount
		userInviteCounts      []*user_invite_count.UserInviteCount
	)

	userInviteCountChild, err := l.svcCtx.UserInviteCountModel.FindOneByUserId(l.ctx, userId)
	if err == gorm.ErrRecordNotFound {
		l.Logger.Errorf("UserInviteCountModel.FindByUserIds  userBind.Id=%v, err=ErrNotFound", userId)
		return resp, nil
	} else if err != nil {
		return nil, err
	}

	if userInvite.ParentId == global.Sys_Id {
		//  系统发放
		inviteRecord, err := l.svcCtx.UserInviteModel.FindOneByParentIdChildId(l.ctx, global.Sys_Id, userId)
		if err == gorm.ErrRecordNotFound {
			l.Logger.Errorf("SysInviteModel child_id=%v, err=ErrNotFound", userId)
		} else if err != nil {
			return nil, err
		}
		inviteCodeParent = inviteRecord.InviteCodeParent

		totalCredit = userInviteCountChild.TotalCredit
		totalInviteCount = userInviteCountChild.TotalCount
		successInviteCount = userInviteCountChild.SuccessCount
		inviteCodeSelf = userInviteCountChild.InviteCode
	} else {
		userInviteCountParent, err = l.svcCtx.UserInviteCountModel.FindOneByUserId(l.ctx, userInvite.ParentId)
		if err == gorm.ErrRecordNotFound {
			l.Logger.Errorf("UserInviteCountModel.FindByUserIds parent_id=%v, err=ErrNotFound", userInvite.ParentId)
		} else if err != nil {
			return nil, err
		}
		userInviteCounts = append(userInviteCounts, userInviteCountParent, userInviteCountChild)

		if userInviteCounts[0].UserId == userInvite.ParentId {
			inviteCodeParent = userInviteCounts[0].InviteCode

			totalCredit = userInviteCounts[1].TotalCredit
			totalInviteCount = userInviteCounts[1].TotalCount
			successInviteCount = userInviteCounts[1].SuccessCount
			inviteCodeSelf = userInviteCounts[1].InviteCode

		} else if userInviteCounts[1].UserId == userInvite.ParentId {
			inviteCodeParent = userInviteCounts[1].InviteCode

			totalCredit = userInviteCounts[0].TotalCredit
			totalInviteCount = userInviteCounts[0].TotalCount
			successInviteCount = userInviteCounts[0].SuccessCount
			inviteCodeSelf = userInviteCounts[0].InviteCode
		}
	}

	resp = &types.InviteListResp{
		TotalCredit:        int(totalCredit),
		TotalInviteCount:   int(totalInviteCount),
		SuccessInviteCount: int(successInviteCount), /// modifyfuture
		SelfCode:           inviteCodeSelf,
		ParentCode:         inviteCodeParent,
	}
	return resp, nil
}
