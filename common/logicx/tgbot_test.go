package logicx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTgBotStart(t *testing.T) {
	require := require.New(t)

	ctx, svcCtx := TestConfig(t)
	svcCtx.Config.TgWebHook = "https://ec2-52-77-241-219.ap-southeast-1.compute.amazonaws.com:8443"
	err := TgBotStart(ctx, svcCtx)
	require.NoError(err)

}
