package logicx

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/Savvy-Gameing/backend/common/global"
	"github.com/Savvy-Gameing/backend/internal/svc"

	"github.com/zeromicro/go-zero/core/logc"
)

func TgBotStart(ctx context.Context, svcCtx *svc.ServiceContext) error {

	logc.Infof(ctx, "svcCtx.Config.TgWebHook:%v", svcCtx.Config.TgWebHook)

	bot, err := tgbotapi.NewBotAPI(global.TG_BOT_TOKEN)
	if err != nil {
		logc.Errorf(ctx, "[TgBotStart] NewBotAPI err: %v", err)
		return err
	}

	// 设置 Webhook 公网地址
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(svcCtx.Config.TgWebHook))
	if err != nil {
		logc.Errorf(ctx, "[TgBotStart] NewWebhook tgwebhook: %v, err: %v", svcCtx.Config.TgWebHook, err)
		return err
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		logc.Errorf(ctx, "[TgBotStart] GetWebhookInfo tgwebhook: %v, err: %v", svcCtx.Config.TgWebHook, err)
		return err
	}
	if info.LastErrorDate != 0 {
		logc.Errorf(ctx, "[TgBotStart] GetWebhookInfo LastErrorMessage: %s", info.LastErrorMessage)
		return fmt.Errorf(info.LastErrorMessage)
	}

	// 处理 Telegram 发来的更新
	updates := bot.ListenForWebhook("/")

	uurl, err := url.Parse(svcCtx.Config.TgWebHook)
	if err != nil {
		logc.Errorf(ctx, "[TgBotStart] url.Parse: url: %v err: %v", svcCtx.Config.TgWebHook, err)
		return err
	}
	strs := strings.Split(uurl.Host, ":")
	var port string
	if len(strs) > 1 {
		port = strs[1]
	} else {
		return fmt.Errorf("[TgBotStart] yaml's TgWebHook: %v hasnot port.", svcCtx.Config.TgWebHook)
	}
	go http.ListenAndServe("0.0.0.0:"+port, nil)

	logc.Info(ctx, "ListenAndServe:", port)

	for update := range updates {

		logc.Debugf(ctx, "update: %+v", update)

		if update.Message == nil { // 忽略非消息更新
			continue
		}

		switch update.Message.Text {
		case "/start":
			// 生成一个随机邀请链接
			inviteLink, err := bot.GetInviteLink(tgbotapi.ChatConfig{
				ChatID: update.Message.Chat.ID,
			})
			if err != nil {
				logc.Errorf(ctx, "[TgBotStart] bot.GetInviteLink  chatid:%v err: %v", update.Message.Chat.ID, err)
				continue
			}

			logc.Info(ctx, "GetInviteLink:", inviteLink)

			// 回复用户包含邀请链接的消息
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, inviteLink)
			_, err = bot.Send(msg)
			if err != nil {
				logc.Errorf(ctx, "[TgBotStart] tgbotapi.NewMessage  chatid:%v err: %v", update.Message.Chat.ID, err)
			}
		}
	}

	return nil
}
