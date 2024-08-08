package tgbot

import (
	"context"
	"fmt"
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/Savvy-Gameing/backend/common/global"
	"github.com/Savvy-Gameing/backend/internal/svc"
)

type TgbotLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTgbotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TgbotLogic {
	return &TgbotLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TgbotLogic) Start() {
	l.Logger.Infof("start Tgbot.")

	bot, err := l.initTgbot()
	if err != nil {
		return
	}
	l.svcCtx.TgBot = bot

	return
}

func (l *TgbotLogic) Stop() {
	l.Logger.Infof("stop Tgbot.")

}

func (l *TgbotLogic) initTgbot() (bot *tgbotapi.BotAPI, err error) {
	l.Logger.Info("init Tgbot.")

	bot, err = tgbotapi.NewBotAPI(global.TG_BOT_TOKEN)
	if err != nil {
		l.Logger.Errorf("[initTgbot] NewBotAPI err: %v", err)
		return bot, err
	}
	bot.Debug = true

	// me, err := bot.GetMe()
	// if err != nil {
	// 	l.Logger.Errorf("[initTgbot] GetMe err: %v", err)
	// 	return bot, err
	// }
	// l.Logger.Infof("me:%+v", me)

	setCommands := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     "/start",
			Description: "go into savvy_game_bot",
		},
		tgbotapi.BotCommand{
			Command:     "/launch",
			Description: "launch savvy web app",
		})

	if _, err = bot.Request(setCommands); err != nil {
		l.Logger.Errorf("Unable to set commands")
		return bot, err
	}

	// commands, err := bot.GetMyCommands()
	// if err != nil {
	// 	l.Logger.Errorf("Unable to get commands")
	// 	return bot, err
	// }
	// l.Logger.Debugf("commands:%+v", commands)

	// 一个Telegram Bot只能设置一个Webhook
	// https://domain 带证书
	whurl, err := url.JoinPath(l.svcCtx.Config.TgWebHook, "webhook")
	if err != nil {
		l.Logger.Errorf("[initTgbot] JoinPath err: %v", err)
		return bot, err
	}

	wh, err := tgbotapi.NewWebhookWithCert(whurl, tgbotapi.FilePath(l.svcCtx.Config.TgPublicPem))
	if err != nil {
		l.Logger.Errorf("[initTgbot] NewWebhookWithCert whurl:%v, TgPublicPem:%v, err: %v", whurl, l.svcCtx.Config.TgPublicPem, err)
		return bot, err
	}
	_, err = bot.Request(wh)
	if err != nil {
		l.Logger.Errorf("[initTgbot] NewWebhook err: %v", err)
		return bot, err
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		l.Logger.Errorf("[initTgbot] GetWebhookInfo  err: %v", err)
		return bot, err
	}
	if info.LastErrorDate != 0 {
		l.Logger.Errorf("[initTgbot] GetWebhookInfo LastErrorMessage: %s", info.LastErrorMessage)
		return bot, fmt.Errorf(info.LastErrorMessage)
	}

	return bot, nil

	// threading.RunSafe(func() {

	// 	// 每天5时3分3秒
	// 	if _, err := l.svcCtx.Tgbot.AddFunc("3 3 5 * * ?", func() { l.userLikeWriteDB() }); err != nil {
	// 		l.Logger.Errorf("[initTgbot] userLikeWriteDB error: %v", err)
	// 		return
	// 	}

	// })
}

func (l *TgbotLogic) verifyInviteCode() {

	// logicx.VerifyInviteCode(l.ctx,l.svcCtx,)

}

func (l *TgbotLogic) InviteHistory() {

	// logicx.VerifyInviteCode(l.ctx,l.svcCtx,)

}
