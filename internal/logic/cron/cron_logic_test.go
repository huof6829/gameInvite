package cron

import (
	"testing"

	"github.com/Savvy-Gameing/backend/common/logicx"
)

func TestUserLikeWriteDB(t *testing.T) {
	// require := require.New(t)

	NewCronLogic(logicx.TestConfig(t)).userLikeWriteDB()

}
