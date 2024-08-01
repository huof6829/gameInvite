package logicx

import (
	"context"

	xerrors "github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/Savvy-Gameing/backend/common/global"
	"github.com/Savvy-Gameing/backend/common/response"
	"github.com/Savvy-Gameing/backend/common/utils"
	"github.com/Savvy-Gameing/backend/internal/model/sys_invite"
	"github.com/Savvy-Gameing/backend/internal/model/user_invite"
	"github.com/Savvy-Gameing/backend/internal/model/user_invite_count"
	"github.com/Savvy-Gameing/backend/internal/svc"

	"github.com/zeromicro/go-zero/core/logc"
)

// modifyfuture
// 成功购买 balanceid , 才算成功推荐
// 直接邀请关系，A B各获得100
// 间接： C完成所有积分任务（绑定、关注、）之后，获得100，B获得10，A获得5
// 一个用户一个邀请码，可以多次邀请
// 1. B带上 A invite_code 连接钱包
// 2. B填写 A invite_code  点领取

// A->B 邀请进来，B获得100，A获得100；B完成（授权）任务，获得对应积分，上级获得分成10%，上上级5%；所有（授权）任务完成获得100，上级10%，上上级5%，B总共获得800

func VerifyInviteCode(ctx context.Context, svcCtx *svc.ServiceContext, childId int64, inviteCode string) (err error) {
	defer func() {
		if err != nil {
			logc.Error(ctx, xerrors.WithStack(err))
		}
	}()

	/// 不带邀请码
	if inviteCode == "" {
		return nil
	}

	var (
		parentId           int64
		totalCreditParent  int64
		totalCountParent   int64
		successCountParent int64
		inviteCodeParent   string = inviteCode

		totalCreditChild  int64
		totalCountChild   int64
		successCountChild int64
	)

	logc.Debugf(ctx, "childId=%v, inviteCode=%v", childId, inviteCode)

	// B的
	_, err = svcCtx.UserInviteCountModel.FindOneByUserId(ctx, childId)
	if err == gorm.ErrRecordNotFound {
		totalCreditChild = global.Invite_Credit_Direct_1
	} else if err != nil {
		return err
	} else {
		/// 已经邀请过
		return response.InvitedAlready
	}

	// A->B A的
	userInviteCountParent, err := svcCtx.UserInviteCountModel.FindOneByInviteCode(ctx, inviteCode)
	if err == gorm.ErrRecordNotFound {
		_, err = svcCtx.SysInviteModel.FindOneByInviteCode(ctx, inviteCode)
		if err == gorm.ErrRecordNotFound {
			logc.Errorf(ctx, "SysInviteModel.FindOneByInviteCode invite_code=%v, err=ErrNotFound", inviteCode)
			return response.InviteCodeIncorrect
		} else if err != nil {
			return err
		}
		parentId = global.Sys_Id
	} else if err != nil {
		return err
	} else {
		parentId = userInviteCountParent.UserId
		totalCreditParent = userInviteCountParent.TotalCredit
		totalCountParent = userInviteCountParent.TotalCount
		successCountParent = userInviteCountParent.SuccessCount
	}

	///
	totalCreditParent += global.Invite_Credit_Direct_1
	totalCountParent += 1
	successCountParent += 1 /// modifyfuture

	var userInviteCounts []*user_invite_count.UserInviteCount
	inviteCodeChild := utils.GetRandomString(global.InviteCode_Length)

	logc.Debugf(ctx, "parentId=%v", parentId)

	if parentId == global.Sys_Id {
		userInviteCounts = append(userInviteCounts,
			&user_invite_count.UserInviteCount{
				UserId:       childId,
				InviteCode:   inviteCodeChild, /// 领取parent，生成
				TotalCredit:  totalCreditChild,
				TotalCount:   totalCountChild,
				SuccessCount: successCountChild,
			})

		if err = svcCtx.DB.Transaction(func(tx *gorm.DB) error {
			if err = svcCtx.SysInviteModel.InsertOnUpdate(ctx, tx, &sys_invite.SysInvite{
				InviteCode: inviteCodeParent,
			}); err != nil {
				return err
			}

			if err = svcCtx.UserInviteModel.InsertOnUpdate(ctx, tx, &user_invite.UserInvite{
				ParentId:                 parentId,
				ChildId:                  childId,
				Level:                    global.Invite_Level_1,
				InviteCreditDirectParent: global.Invite_Credit_Direct_1,
				InviteCreditDirectChild:  global.Invite_Credit_Direct_1,
				InviteCodeParent:         inviteCodeParent,
			}); err != nil {
				return err
			}

			/// self
			return svcCtx.UserInviteCountModel.InsertOnUpdate(ctx, tx, userInviteCounts[0])
		}); err != nil {
			logc.Errorf(ctx, "[VerifyInviteCode] Transaction err=%v", err)
			return err
		}

	} else {
		userInviteCounts = append(userInviteCounts,
			&user_invite_count.UserInviteCount{
				UserId:       childId,
				InviteCode:   inviteCodeChild, /// 领取parent，生成
				TotalCredit:  totalCreditChild,
				TotalCount:   totalCountChild,
				SuccessCount: successCountChild,
			},
			&user_invite_count.UserInviteCount{
				UserId:       parentId,
				InviteCode:   inviteCodeParent,
				TotalCredit:  totalCreditParent,
				TotalCount:   totalCountParent,
				SuccessCount: successCountParent,
			})

		var (
			parentId_2         int64
			userInvitePraent_2 *user_invite.UserInvite
		)
		userInvitePraents, err := svcCtx.UserInviteModel.FindByChildIdLevels(ctx, svcCtx.DB.DB, parentId, []int64{global.Invite_Level_1})
		if err == gorm.ErrRecordNotFound {
			logc.Errorf(ctx, " [VerifyInviteCode] UserInviteModel.FindByChildIdLevels parentId=%v, level=1, err=%v", parentId, err)
			return response.InvitedNever
		} else if err != nil {
			return err
		} else if len(userInvitePraents) == 0 {
			return response.InvitedNever
		} else {
			userInvitePraent_2 = userInvitePraents[0]
			parentId_2 = userInvitePraent_2.ParentId
		}

		if err = svcCtx.DB.Transaction(func(tx *gorm.DB) error {
			if parentId_2 != 0 { // 存在父父级
				if err = svcCtx.UserInviteModel.Insert(ctx, tx, &user_invite.UserInvite{
					ParentId: parentId_2,
					ChildId:  childId,
					Level:    global.Invite_Level_2,
				}); err != nil {
					return err
				}
			}

			if err = svcCtx.UserInviteModel.InsertOnUpdate(ctx, tx, &user_invite.UserInvite{
				ParentId:                 parentId,
				ChildId:                  childId,
				Level:                    global.Invite_Level_1,
				InviteCreditDirectParent: global.Invite_Credit_Direct_1,
				InviteCreditDirectChild:  global.Invite_Credit_Direct_1,
			}); err != nil {
				return err
			}

			if err = svcCtx.UserInviteCountModel.InsertOnUpdate(ctx, tx, userInviteCounts[0]); err != nil {
				return err
			}
			return svcCtx.UserInviteCountModel.InsertOnUpdate(ctx, tx, userInviteCounts[1])
		}); err != nil {
			logc.Errorf(ctx, "[VerifyInviteCode] Transaction err=%v", err)
			return err
		}
	}

	return err
}

// 完成所有任务，通过邀请进来的，获得分成奖励
func IndirectInviteCreditSum(ctx context.Context, svcCtx *svc.ServiceContext, childId int64) (err error) {

	var (
		parentId_1, parentId_2                 int64
		userInvitePraent_1, userInvitePraent_2 *user_invite.UserInvite
	)

	userInvites, err := svcCtx.UserInviteModel.FindByChildIdLevels(ctx, svcCtx.DB.DB, childId, []int64{global.Invite_Level_1, global.Invite_Level_2})
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	} else if err == gorm.ErrRecordNotFound {
		return nil
	} else if len(userInvites) == 0 {
		return nil
	}

	for k := 0; k < len(userInvites); k++ {
		if userInvites[k].Level == global.Invite_Level_1 {
			userInvitePraent_1 = userInvites[k]
			parentId_1 = userInvitePraent_1.ParentId
		} else if userInvites[k].Level == global.Invite_Level_2 {
			userInvitePraent_2 = userInvites[k]
		}
	}

	if userInvitePraent_1 == nil {
		logc.Errorf(ctx, "UserInviteModel.FindByChildIdLevels userInvites=%+v, err=ErrNotFound", userInvites)
		return response.InvitedNever
	}

	if userInvitePraent_1.InviteCreditIndirectChild == global.Invite_Credit_Indirect_1_Child {
		return nil /// 已经获得
	}

	logc.Debugf(ctx, "userInvitePraent_1=%+v", userInvitePraent_1)
	logc.Debugf(ctx, "userInvitePraent_2=%+v", userInvitePraent_2)

	userInvite_1 := &user_invite.UserInvite{
		Id:                         userInvitePraent_1.Id,
		ChildId:                    childId,
		Level:                      global.Invite_Level_1,
		InviteCreditIndirectParent: global.Invite_Credit_Indirect_1_Parent,
		InviteCreditIndirectChild:  global.Invite_Credit_Indirect_1_Child,
	}

	var userInvite_2 *user_invite.UserInvite
	if userInvitePraent_2 != nil {
		userInvite_2 = &user_invite.UserInvite{
			Id:                         userInvitePraent_2.Id,
			ChildId:                    childId,
			Level:                      global.Invite_Level_2,
			InviteCreditIndirectParent: global.Invite_Credit_Indirect_2_Parent,
			InviteCreditIndirectChild:  global.Invite_Credit_Indirect_2_Child,
		}
		parentId_2 = userInvitePraent_2.ParentId
	}

	if err = svcCtx.DB.Transaction(func(tx *gorm.DB) error {

		if err = svcCtx.UserInviteModel.UpdateByChildIdLevel(ctx, tx, userInvite_1); err != nil {
			return err
		}

		if userInvitePraent_2 != nil {
			if err = svcCtx.UserInviteModel.UpdateByChildIdLevel(ctx, tx, userInvite_2); err != nil {
				return err
			}
		}

		if err = svcCtx.UserInviteCountModel.UpdateByUserId(ctx, tx, &user_invite_count.UserInviteCount{
			UserId: childId,
		}, global.Invite_Credit_Indirect_1_Child); err != nil {
			return err
		}

		if parentId_1 > 0 {
			if err = svcCtx.UserInviteCountModel.UpdateByUserId(ctx, tx, &user_invite_count.UserInviteCount{
				UserId: parentId_1,
			}, global.Invite_Credit_Indirect_1_Parent); err != nil {
				return err
			}
		}

		if parentId_2 > 0 {
			if err = svcCtx.UserInviteCountModel.UpdateByUserId(ctx, tx, &user_invite_count.UserInviteCount{
				UserId: parentId_2,
			}, global.Invite_Credit_Indirect_2_Parent); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		logc.Errorf(ctx, "[IndirectInviteCreditSum] Transaction err=%v", err)
		return err
	}

	return err
}

// userInvite, err := l.svcCtx.UserInviteModel.FindOneByParentIdChildId(l.ctx, parentId, childId)
// if err == gorm.ErrRecordNotFound {
// 	logc.Errorf(ctx, "UserInviteModel.FindOneByParentIdChildId parentId=%v, childId=%v, err=ErrNotFound", parentId, childId)
// 	return nil, response.InviteCodeIncorrect
// } else if err != nil {
// 	return nil, err
// }
// if userInvite.InviteCreditDirectChild == utils.Invite_Credit_Direct_1 && userInvite.Level == utils.Invite_Level_1 {
// 	return nil, response.InvitedAlready
// }

// totalCreditChild = userInviteCountChild.TotalCredit
// totalCountChild = userInviteCountChild.TotalCount
// successCountChild = userInviteCountChild.SuccessCount
