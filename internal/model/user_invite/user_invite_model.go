package user_invite

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ UserInviteModel = (*customUserInviteModel)(nil)

type (
	// UserInviteModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserInviteModel.
	UserInviteModel interface {
		userInviteModel
		customUserInviteLogicModel
	}

	customUserInviteModel struct {
		*defaultUserInviteModel
	}

	customUserInviteLogicModel interface {
		FindByChildIdLevels(ctx context.Context, tx *gorm.DB, childId int64, levels []int64) ([]*UserInvite, error)
		FindByParentLevels(ctx context.Context, tx *gorm.DB, parentId int64, levels []int64, limit int, order string) ([]*UserInvite, error)
		InsertOnUpdate(ctx context.Context, tx *gorm.DB, data *UserInvite) error
		UpdateByChildIdLevel(ctx context.Context, tx *gorm.DB, data *UserInvite) error
	}
)

// NewUserInviteModel returns a model for the database table.
func NewUserInviteModel(conn *gorm.DB, c cache.CacheConf) UserInviteModel {
	return &customUserInviteModel{
		defaultUserInviteModel: newUserInviteModel(conn, c),
	}
}

func (m *defaultUserInviteModel) getNewModelNeedReloadCacheKeys(data *UserInvite) []string {
	if data == nil {
		return []string{}
	}
	return []string{}
}
func (m *defaultUserInviteModel) customCacheKeys(data *UserInvite) []string {
	if data == nil {
		return []string{}
	}
	return []string{}
}

func (m *customUserInviteModel) FindByChildIdLevels(ctx context.Context, tx *gorm.DB, childId int64, levels []int64) ([]*UserInvite, error) {
	var result []*UserInvite
	err := m.ExecNoCache(func(conn *gorm.DB) error {
		db := conn
		if tx != nil {
			db = tx
		}
		return db.WithContext(ctx).
			Where("child_id = ? and level in (?)", childId, levels).
			Find(&result).Error
	})
	return result, err
}

func (m *customUserInviteModel) FindByParentLevels(ctx context.Context, tx *gorm.DB, parentId int64, levels []int64, limit int, order string) ([]*UserInvite, error) {
	var result []*UserInvite
	err := m.ExecNoCache(func(conn *gorm.DB) error {
		db := conn
		if tx != nil {
			db = tx
		}
		return db.WithContext(ctx).
			Where("parent_id = ? and level in (?)", parentId, levels).
			Order(order).Limit(limit).
			Find(&result).Error
	})
	return result, err
}

func (m *customUserInviteModel) InsertOnUpdate(ctx context.Context, tx *gorm.DB, data *UserInvite) error {
	err := m.ExecCtx(ctx, func(conn *gorm.DB) error {
		db := conn
		if tx != nil {
			db = tx
		}
		return db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"parent_id", "child_id", "level", "invite_credit_direct_parent", "invite_credit_direct_child", "invite_credit_indirect_parent", "invite_credit_indirect_child"}),
		}).Create(data).Error
	}, m.getCacheKeys(data)...)
	return err
}

func (m *customUserInviteModel) UpdateByChildIdLevel(ctx context.Context, tx *gorm.DB, data *UserInvite) error {
	old, err := m.FindOne(ctx, data.Id)
	if err != nil && err != ErrNotFound {
		return err
	}
	clearKeys := append(m.getCacheKeys(old), m.getNewModelNeedReloadCacheKeys(data)...)
	err = m.ExecCtx(ctx, func(conn *gorm.DB) error {
		db := conn
		if tx != nil {
			db = tx
		}

		/// data带id值，where 自动带id，返回计算值
		return db.Model(data).Clauses(
			// clause.Returning{Columns: []clause.Column{{Name: "invite_credit_indirect_parent"}}},
			clause.Returning{},
		).Where("child_id = ? and level = ?", data.ChildId, data.Level).Updates(map[string]interface{}{
			"invite_credit_indirect_parent": data.InviteCreditIndirectParent,
			"invite_credit_indirect_child":  data.InviteCreditIndirectChild,
		}).Error
	}, clearKeys...)
	return err
}
