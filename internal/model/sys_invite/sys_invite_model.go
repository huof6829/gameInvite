package sys_invite

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"gorm.io/gorm"
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
