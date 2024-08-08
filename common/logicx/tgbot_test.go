package logicx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTgBotStart(t *testing.T) {
	require := require.New(t)
	ctx, svcCtx := TestConfig(t)

	// vercel 运行才有效
	svcCtx.Config.TgWebHook = "https://game-invite.vercel.app:8443"
	svcCtx.Config.TgPublicPem = "/home/mart/selfca/YOURPUBLIC.pem"
	svcCtx.Config.TgPrivateKey = "/home/mart/selfca/YOURPRIVATE.key"
	err := TgBotStart(ctx, svcCtx)
	require.NoError(err)

}
