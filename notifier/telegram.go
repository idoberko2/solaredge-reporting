package notifier

import (
	"github.com/pkg/errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

func NewTelegram(cfg TelegramConfig) Notifier {
	return &telegramNotifier{
		cfg: cfg,
	}
}

type telegramNotifier struct {
	cfg TelegramConfig
	bot *tgbotapi.BotAPI
}

func (t *telegramNotifier) Init() error {
	bot, err := tgbotapi.NewBotAPI(t.cfg.Token)
	if err != nil {
		return err
	}

	if t.cfg.IsDebug {
		bot.Debug = true
	}

	t.bot = bot

	if t.cfg.IsListen {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		updates := bot.GetUpdatesChan(u)

		go func() {
			for update := range updates {
				if update.Message != nil {
					log.WithField("chat_id", update.Message.Chat.ID).Info("received message")
					return
				}
			}
		}()
	}

	return nil
}

var ErrNotInitialized = errors.New("not initialized")
