package modules

import (
	"context"
	"time"

	"github.com/ashwinath/moneybags/pbgo/configpb"
	telegramprocessor "github.com/ashwinath/moneybags/pkg/telegram/processor"
	"github.com/ashwinath/simple/framework"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type TelegramModule struct {
	fw               framework.FW
	bot              *bot.Bot
	processorManager *telegramprocessor.ProcessorManager
}

func NewTelegramModule(fw framework.FW) (framework.Module, error) {
	pm, err := telegramprocessor.NewManager(fw)
	if err != nil {
		return nil, err
	}
	tm := &TelegramModule{
		fw:               fw,
		processorManager: pm,
	}
	opts := []bot.Option{
		bot.WithDefaultHandler(tm.handler),
	}

	b, err := bot.New(fw.GetConfig().(*configpb.Config).TelegramConfig.ApiKey, opts...)
	if err != nil {
		return nil, err
	}

	tm.bot = b
	self, err := b.GetMe(context.Background())
	if err != nil {
		fw.GetLogger().Errorf("unable to get telegram user information: %s", err)
	}

	fw.GetLogger().Infof("Authorized telegram bot on account %s", self.Username)

	return tm, nil
}

func (m *TelegramModule) Name() string {
	return "telegram"
}

func (m *TelegramModule) Start(ctx context.Context) {
	m.bot.Start(ctx)
}

func (m *TelegramModule) handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	if update.Message.From == nil {
		return
	}

	if update.Message.From.Username == m.fw.GetConfig().(*configpb.Config).TelegramConfig.AllowedUser { // If we got a message
		m.fw.GetLogger().Infof("[telegram] [%s to bot] %s", update.Message.From.Username, update.Message.Text)

		reply := m.processorManager.ProcessMessage(update.Message.Text, time.Unix(int64(update.Message.Date), 0))
		if reply == nil {
			return
		}

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      *reply,
			ParseMode: models.ParseModeMarkdown,
			ReplyParameters: &models.ReplyParameters{
				MessageID: update.Message.ID,
			},
		})
	}
}
