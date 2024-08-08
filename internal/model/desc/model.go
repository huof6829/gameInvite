package desc

import (
	"time"
)

type SysInvite struct {
	Id         int64     `gorm:"column:id;type:bigint(20);autoIncrement;not null;primaryKey"` // auto id
	InviteCode string    `gorm:"column:invite_code;type:varchar(200);not null;unique"`
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:timestamp;not null"`
}

type UserInvite struct {
	Id                         int64     `gorm:"column:id;type:bigint(20);autoIncrement;not null;primaryKey"` // auto id
	ParentId                   int64     `gorm:"column:parent_id;type:bigint(20);not null;index;uniqueIndex:uk_user_invite_parent_child_id"`
	ChildId                    int64     `gorm:"column:child_id;type:bigint(20);not null;index;uniqueIndex:uk_user_invite_parent_child_id"`
	Level                      int64     `gorm:"column:level;type:tinyint(4);not null;index;comment:1-直接,2-间接"`
	InviteCreditDirectParent   int64     `gorm:"column:invite_credit_direct_parent;type:int(11);not null;comment:直接邀请积分"`
	InviteCreditDirectChild    int64     `gorm:"column:invite_credit_direct_child;type:int(11);not null;comment:间接邀请积分"`
	InviteCreditIndirectParent int64     `gorm:"column:invite_credit_indirect_parent;type:int(11);not null;comment:分成邀请积分,直接10%,间接5%"`
	InviteCreditIndirectChild  int64     `gorm:"column:invite_credit_indirect_child;type:int(11);not null;comment:分成邀请积分，直接邀请填入，间接邀请不填"`
	InviteCodeParent           string    `gorm:"column:invite_code_parent;type:varchar(20);not null;comment:系统邀请码"`
	CreatedAt                  time.Time `gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt                  time.Time `gorm:"column:updated_at;type:timestamp;not null"`
}

type UserInviteCount struct {
	Id           int64     `gorm:"column:id;type:bigint(20);autoIncrement;not null;primaryKey"` // auto id
	UserId       int64     `gorm:"column:user_id;type:bigint(20);not null;unique"`
	InviteCode   string    `gorm:"column:invite_code;type:varchar(20);not null;unique"`
	TotalCredit  int64     `gorm:"column:total_credit;type:int(11);not null;comment:邀请总分"`
	TotalCount   int64     `gorm:"column:total_count;type:int(11);not null;comment:要去总数"`
	SuccessCount int64     `gorm:"column:success_count;type:int(11);not null;comment:成功邀请总数"`
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:timestamp;not null"`
}
