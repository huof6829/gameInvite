package user_invite_count

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ UserInviteCountModel = (*customUserInviteCountModel)(nil)

type (
	// UserInviteCountModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserInviteCountModel.
	UserInviteCountModel interface {
		userInviteCountModel
		customUserInviteCountLogicModel
	}

	customUserInviteCountModel struct {
		*defaultUserInviteCountModel
	}

	customUserInviteCountLogicModel interface {
		InsertOnUpdate(ctx context.Context, tx *gorm.DB, data *UserInviteCount) error
		FindByUserIds(ctx context.Context, tx *gorm.DB, userIds []int64) ([]*UserInviteCount, error)
		UpdateByUserId(ctx context.Context, tx *gorm.DB, data *UserInviteCount, credit int64) error
	}
)

// NewUserInviteCountModel returns a model for the database table.
func NewUserInviteCountModel(conn *gorm.DB, c cache.CacheConf) UserInviteCountModel {
	return &customUserInviteCountModel{
		defaultUserInviteCountModel: newUserInviteCountModel(conn, c),
	}
}

func (m *defaultUserInviteCountModel) getNewModelNeedReloadCacheKeys(data *UserInviteCount) []string {
	if data == nil {
		return []string{}
	}
	return []string{}
}
func (m *defaultUserInviteCountModel) customCacheKeys(data *UserInviteCount) []string {
	if data == nil {
		return []string{}
	}
	return []string{}
}

func (m *customUserInviteCountModel) InsertOnUpdate(ctx context.Context, tx *gorm.DB, data *UserInviteCount) error {
	err := m.ExecCtx(ctx, func(conn *gorm.DB) error {
		db := conn
		if tx != nil {
			db = tx
		}

		return db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"user_id", "invite_code", "total_credit", "total_count", "success_count"}),
		}).Create(data).Error
	}, m.getCacheKeys(data)...)
	return err
}

func (m *customUserInviteCountModel) FindByUserIds(ctx context.Context, tx *gorm.DB, userIds []int64) ([]*UserInviteCount, error) {
	var result []*UserInviteCount
	err := m.ExecNoCache(func(conn *gorm.DB) error {
		db := conn
		if tx != nil {
			db = tx
		}
		return db.WithContext(ctx).
			Where("user_id in (?)", userIds).
			Find(&result).Error
	})
	return result, err
}

func (m *customUserInviteCountModel) UpdateByUserId(ctx context.Context, tx *gorm.DB, data *UserInviteCount, credit int64) error {
	olda, err := m.FindOneByUserId(ctx, data.UserId)
	if err != nil {
		return err
	}
	old, err := m.FindOne(ctx, olda.Id)
	if err != nil && err != ErrNotFound {
		return err
	}
	clearKeys := append(m.getCacheKeys(old), m.getNewModelNeedReloadCacheKeys(data)...)
	err = m.ExecCtx(ctx, func(conn *gorm.DB) error {
		db := conn
		if tx != nil {
			db = tx
		}
		return db.Model(data).Where("user_id = ?", data.UserId).Updates(map[string]interface{}{
			"total_credit": gorm.Expr("total_credit + ?", credit),
		}).Error
	}, clearKeys...)
	return err
}
