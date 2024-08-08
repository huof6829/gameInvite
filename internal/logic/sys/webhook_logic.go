package sys

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/Savvy-Gameing/backend/internal/svc"
)

type WebhookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WebhookLogic {
	return &WebhookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WebhookLogic) Webhook(r *http.Request) error {
	// todo: add your logic here and delete this line

	if l.svcCtx.TgBot == nil {
		return fmt.Errorf("tg bot start failed.")
	}

	update, err := l.svcCtx.TgBot.HandleUpdate(r)

	if update.Message == nil { // 忽略非消息更新
		return nil
	}

	if strings.HasPrefix(update.Message.Text, "/start") {
		code, _ := strings.CutPrefix(update.Message.Text, "/start")
		code = strings.TrimSpace(code)

		// 获取用户信息
		user := update.Message.From

		l.Logger.Debugf("user:%+v, message:%+v, code:%v", user, update.Message, code)

		// 记录邀请码，用户信息

		// 回复用户
		// reply := "Hello, @" + user.FirstName + "! I received your message."
		reply := fmt.Sprintf(`Account ID: %v
Account Name: %v
Account activated successfully!
You're in!`,
			user.ID, user.FirstName)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		_, err := l.svcCtx.TgBot.Send(msg)
		if err != nil {
			l.Logger.Errorf("[initTgbot] bot.Send chatid:%v, err: %v", update.Message.Chat.ID, err)
		}

	} else if update.Message.Text == "/launch" {
		webAppButton := tgbotapi.NewInlineKeyboardButtonURL("Open Web App", "https://t.me/mart_jim_bot/drabapp") // "https://t.me/eden_savvy_game_bot/eden_savvy_game"
		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(webAppButton),
		)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, `https://t.me/mart_jim_bot/drabapp
Click the button to open the Web App!`)
		msg.ReplyMarkup = inlineKeyboard
		_, err = l.svcCtx.TgBot.Send(msg)
		if err != nil {
			l.Logger.Errorf("[initTgbot] bot.Send chatid:%v, err: %v", update.Message.Chat.ID, err)
		}
	}

	return nil
}
