package global

const (
	TimeFormat = "2006-01-02 15:04:05"

	// mysql
	DBFindLimit = 300
	DBBatchSize = 200

	// redis
	CacheExpireCtx           = 300
	CacheListExpireTime      = 3600 * 24 * 2
	CacheListDefaultPageSize = 20
)

// A->B->C->D
const (
	Invite_Level_1 = 1 // C->D
	Invite_Level_2 = 2 // B->D

	InviteCode_Length = 9

	Invite_Credit_Direct_1          = 100 /// 直接邀请
	Invite_Credit_Indirect_1_Child  = 100 /// 分成邀请 完成任务后
	Invite_Credit_Indirect_2_Child  = 0   /// 直接邀请的分成奖励记过了
	Invite_Credit_Indirect_1_Parent = Invite_Credit_Indirect_1_Child / 10
	Invite_Credit_Indirect_2_Parent = Invite_Credit_Indirect_1_Parent / 2
)

const (
	BotToken = "7480704593:AAF96vAvAn_JNSx2gIOt_ppJzF820jisbWE" // aven test

	// 系统
	Sys_Id       = -1
	Sys_Password = "__sq2val&!093inv"
)
