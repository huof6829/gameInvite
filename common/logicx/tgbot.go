package logicx

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Savvy-Gameing/backend/common/global"
	"github.com/Savvy-Gameing/backend/common/utils"
	"github.com/Savvy-Gameing/backend/internal/svc"

	"github.com/zeromicro/go-zero/core/logc"
)

func TgBotStart(ctx context.Context, svcCtx *svc.ServiceContext) error {

	logc.Infof(ctx, "svcCtx.Config.TgWebHook:%v, svcCtx.Config.TgPublicPem: %v, svcCtx.Config.TgPrivateKey: %v", svcCtx.Config.TgWebHook, svcCtx.Config.TgPublicPem, svcCtx.Config.TgPrivateKey)

	// 生成一个随机邀请链接
	ms := time.Now().UnixMilli()
	code := utils.GetRandomString(12) + fmt.Sprintf("-%v", ms)
	logc.Infof(ctx, "code:%v", code)

	bot, err := tgbotapi.NewBotAPI(global.TG_BOT_TOKEN)
	if err != nil {
		logc.Errorf(ctx, "[TgBotStart] NewBotAPI err: %v", err)
		return err
	}
	bot.Debug = true

	me, err := bot.GetMe()
	if err != nil {
		logc.Errorf(ctx, "[TgBotStart] NewWebhook tgwebhook: %v, err: %v", svcCtx.Config.TgWebHook, err)
		return err
	}
	logc.Infof(ctx, "me:%+v", me)

	setCommands := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     "/start",
			Description: "go into savvy_game_bot",
		},
		tgbotapi.BotCommand{
			Command:     "/launch",
			Description: "launch savvy web app",
		})

	if _, err := bot.Request(setCommands); err != nil {
		logc.Errorf(ctx, "Unable to set commands")
		return err
	}

	commands, err := bot.GetMyCommands()
	if err != nil {
		logc.Errorf(ctx, "Unable to get commands")
		return err
	}

	logc.Debugf(ctx, "commands:%+v", commands)
	log.Println("")
	log.Println("")

	// 一个Telegram Bot只能设置一个Webhook
	// https://domain 带证书
	whurl, err := url.JoinPath(svcCtx.Config.TgWebHook, "webhook")
	if err != nil {
		logc.Errorf(ctx, "[TgBotStart] JoinPath err: %v", err)
		return err
	}

	wh, err := tgbotapi.NewWebhookWithCert(whurl, tgbotapi.FilePath(svcCtx.Config.TgPublicPem))
	if err != nil {
		logc.Errorf(ctx, "[TgBotStart] NewWebhookWithCert tgwebhook: %v, err: %v", svcCtx.Config.TgWebHook, err)
		return err
	}
	resp, err := bot.Request(wh)
	if err != nil {
		logc.Errorf(ctx, "[TgBotStart] NewWebhook tgwebhook: %v, err: %v", svcCtx.Config.TgWebHook, err)
		return err
	}

	logc.Debugf(ctx, "resp:%+v", resp)
	log.Println("")
	log.Println("")

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
	updates := bot.ListenForWebhook("/webhook")
	go func() {
		// err = http.ListenAndServeTLS("0.0.0.0:8443", svcCtx.Config.TgPublicPem, svcCtx.Config.TgPrivateKey, nil)
		err = http.ListenAndServe("0.0.0.0:8443", nil)
		if err != nil {
			logc.Errorf(ctx, "[TgBotStart] ListenAndServeTLS err: %v", err)
		}
	}()

	// _, err = http.Get("https://api.telegram.org/bot" + bot.Token + "/deleteWebhook")
	// if err != nil {
	// 	logc.Errorf(ctx, "[TgBotStart] deleteWebhook err: %v", err)
	// 	return err
	// }

	// // 轮询
	// u := tgbotapi.NewUpdate(0)
	// u.Timeout = 3600
	// updates := bot.GetUpdatesChan(u)

	for update := range updates {

		logc.Debugf(ctx, "update:%+v", update)
		log.Println("")
		log.Println("")

		if update.Message == nil { // 忽略非消息更新
			continue
		}

		logc.Debugf(ctx, "message:%+v", update.Message)
		log.Println("")
		log.Println("")

		if strings.HasPrefix(update.Message.Text, "/start") {
			code, _ := strings.CutPrefix(update.Message.Text, "/start")
			code = strings.TrimSpace(code)

			// 获取用户信息
			user := update.Message.From

			logc.Debugf(ctx, "user:%+v, message:%+v, code:%v", user, update.Message, code)
			log.Println("")
			log.Println("")

			// 记录邀请码，用户信息

			// 回复用户
			// reply := "Hello, @" + user.FirstName + "! I received your message."
			reply := fmt.Sprintf(`Account ID: %v
Account Name: %v
Account activated successfully!
You're in!`,
				user.ID, user.FirstName)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			_, err := bot.Send(msg)
			if err != nil {
				logc.Errorf(ctx, "[TgBotStart] bot.Send chatid:%v, err: %v", update.Message.Chat.ID, err)
			}

		} else if update.Message.Text == "/launch" {
			webAppButton := tgbotapi.NewInlineKeyboardButtonURL("Open Web App", "https://t.me/mart_jim_bot/drabapp") // "https://t.me/eden_savvy_game_bot/eden_savvy_game"
			inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(webAppButton),
			)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, `https://t.me/mart_jim_bot/drabapp
Click the button to open the Web App!`)
			msg.ReplyMarkup = inlineKeyboard
			_, err = bot.Send(msg)
			if err != nil {
				logc.Errorf(ctx, "[TgBotStart] bot.Send chatid:%v, err: %v", update.Message.Chat.ID, err)
			}
		}

	}

	return nil
}
