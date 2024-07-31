package sys

import (
	"context"
	"fmt"

	"github.com/Savvy-Gameing/backend/common/global"
	"github.com/Savvy-Gameing/backend/common/utils"
	"github.com/Savvy-Gameing/backend/internal/model/sys_invite"
	"github.com/Savvy-Gameing/backend/internal/svc"
	"github.com/Savvy-Gameing/backend/internal/types"

	xerrors "github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type InviteCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInviteCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InviteCreateLogic {
	return &InviteCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InviteCreateLogic) InviteCreate(req *types.InviteCreateReq) (resp *types.InviteCreateResp, err error) {
	// todo: add your logic here and delete this line
	defer func() {
		if err != nil {
			l.Logger.Error(xerrors.WithStack(err))
		}
	}()

	if req.Password != global.Sys_Password {
		return nil, fmt.Errorf("Sys password is wrong. req.Password=%v", req.Password)
	}

	inviteCodes := make([]string, 0, req.Count)
	for i := 0; i < req.Count; i++ {
		inviteCode := utils.GetRandomString(global.InviteCode_Length)
		if err = l.svcCtx.SysInviteModel.Insert(l.ctx, l.svcCtx.DB.DB, &sys_invite.SysInvite{
			InviteCode: inviteCode,
		}); err != nil {
			l.Logger.Errorf("SysInviteModel:InviteCreate  err=%v", err)
		}
		inviteCodes = append(inviteCodes, inviteCode)
	}

	resp = &types.InviteCreateResp{
		InviteCodes: inviteCodes,
	}

	return resp, nil
}
