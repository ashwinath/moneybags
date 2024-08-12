package modules

import (
	"context"
	"log"
	"time"

	"github.com/ashwinath/moneybags/pkg/framework"
	telegramprocessor "github.com/ashwinath/moneybags/pkg/telegram/processor"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramModule struct {
	fw               framework.FW
	bot              *tgbotapi.BotAPI
	processorManager *telegramprocessor.ProcessorManager
}

func NewTelegramModule(fw framework.FW) (framework.Module, error) {
	bot, err := tgbotapi.NewBotAPI(fw.GetConfig().TelegramConfig.ApiKey)

	if err != nil {
		return nil, err
	}

	bot.Debug = fw.GetConfig().TelegramConfig.Debug
	fw.GetLogger().Infof("Authorized telegram bot on account %s", bot.Self.UserName)

	pm, err := telegramprocessor.NewManager(fw)
	if err != nil {
		return nil, err
	}

	return &TelegramModule{
		fw:               fw,
		bot:              bot,
		processorManager: pm,
	}, nil
}

func (m *TelegramModule) Name() string {
	return "telegram"
}

func (m *TelegramModule) Start(ctx context.Context) {
	// Don't need to use context here
	m.fw.GetLogger().Infof("starting telegram module")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := m.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message.From.UserName == m.fw.GetConfig().TelegramConfig.AllowedUser && update.Message != nil { // If we got a message
			m.fw.GetLogger().Infof("[telegram] [%s to bot] %s", update.Message.From.UserName, update.Message.Text)

			reply := m.processorManager.ProcessMessage(update.Message.Text, time.Unix(int64(update.Message.Date), 0))
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, *reply)
			msg.ReplyToMessageID = update.Message.MessageID
			msg.ParseMode = "Markdown"
			m.fw.GetLogger().Infof("[telegram] [bot to %s] %s", update.Message.From.UserName, *reply)

			_, err := m.bot.Send(msg)
			if err != nil {
				log.Printf("[telegram] [bot to %s] error: %s", update.Message.From.UserName, err)
			}
		}
	}
}
