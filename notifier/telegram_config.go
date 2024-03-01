package notifier

type TelegramConfig struct {
	Token    string `required:"true"`
	ChatID   int64  `split_words:"true" required:"true"`
	IsDebug  bool
	IsListen bool `split_words:"true"`
}
