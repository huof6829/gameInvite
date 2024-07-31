package sys_invite

import (
	"context"

	"github.com/SpectatorNan/gorm-zero/gormc"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ SysInviteModel = (*customSysInviteModel)(nil)

type (
	// SysInviteModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysInviteModel.
	SysInviteModel interface {
		sysInviteModel
		customSysInviteLogicModel
	}

	customSysInviteModel struct {
		*defaultSysInviteModel
	}

	customSysInviteLogicModel interface {
		FindOneByChildId(ctx context.Context, tx *gorm.DB, childId int64) (*SysInvite, error)
		InsertOnUpdate(ctx context.Context, tx *gorm.DB, data *SysInvite) error
	}
)

// NewSysInviteModel returns a model for the database table.
func NewSysInviteModel(conn *gorm.DB, c cache.CacheConf) SysInviteModel {
	return &customSysInviteModel{
		defaultSysInviteModel: newSysInviteModel(conn, c),
	}
}

func (m *defaultSysInviteModel) getNewModelNeedReloadCacheKeys(data *SysInvite) []string {
	if data == nil {
		return []string{}
	}
	return []string{}
}
func (m *defaultSysInviteModel) customCacheKeys(data *SysInvite) []string {
	if data == nil {
		return []string{}
	}
	return []string{}
}

func (m *customSysInviteModel) FindOneByChildId(ctx context.Context, tx *gorm.DB, childId int64) (*SysInvite, error) {
	var resp SysInvite
	err := m.ExecNoCache(func(conn *gorm.DB) error {
		db := conn
		if tx != nil {
			db = tx
		}
		return db.WithContext(ctx).
			Where("child_id = ?", childId).
			Find(&resp).Error
	})
	switch err {
	case nil:
		return &resp, nil
	case gormc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customSysInviteModel) InsertOnUpdate(ctx context.Context, tx *gorm.DB, data *SysInvite) error {
	err := m.ExecCtx(ctx, func(conn *gorm.DB) error {
		db := conn
		if tx != nil {
			db = tx
		}
		return db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"child_id", "invite_code", "invite_credit_direct_child"}),
		}).Create(data).Error
	}, m.getCacheKeys(data)...)
	return err
}
