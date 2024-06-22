package notification

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageOption interface {
	Apply(*tgbotapi.MessageConfig)
}

type withButton struct {
	text string
	url  string
}

func (o withButton) Apply(m *tgbotapi.MessageConfig) {
	button := tgbotapi.NewInlineKeyboardButtonURL(o.text, o.url)
	markup := tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{button})

	m.BaseChat.ReplyMarkup = markup
}

func WithButton(text, url string) MessageOption {
	return withButton{text: text, url: url}
}

type MessageService interface {
	Send(context.Context, string, ...MessageOption) error
}

type telegramService struct {
	chatID int64
	client *tgbotapi.BotAPI
}

func NewTelegramService(token string, chatID int64) (*telegramService, error) {
	client, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram bot: %w", err)
	}

	return &telegramService{chatID: chatID, client: client}, nil
}

func (t *telegramService) Send(ctx context.Context, msg string, opts ...MessageOption) error {
	m := tgbotapi.NewMessage(t.chatID, msg)

	for _, opt := range opts {
		opt.Apply(&m)
	}

	_, err := t.client.Send(m)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
