package notification_test

import (
	"context"
	"errors"
	"testing"

	"github.com/glocurrency/commons/notification" // Adjust import path
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockTelegramSender implements notification.TelegramSender
type mockTelegramSender struct {
	CapturedChattable tgbotapi.Chattable
	SendErr           error
}

func (m *mockTelegramSender) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	m.CapturedChattable = c
	return tgbotapi.Message{}, m.SendErr
}

func TestTelegramService_Send(t *testing.T) {
	chatID := int64(123456789)

	t.Run("successfully sends basic message", func(t *testing.T) {
		mockClient := &mockTelegramSender{}
		svc := notification.NewTelegramServiceWithClient(mockClient, chatID)

		err := svc.Send(context.Background(), "Hello World")
		require.NoError(t, err)

		// Assert that the message was formed correctly
		msg, ok := mockClient.CapturedChattable.(tgbotapi.MessageConfig)
		require.True(t, ok)
		assert.Equal(t, chatID, msg.ChatID)
		assert.Equal(t, "Hello World", msg.Text)
		assert.Empty(t, msg.ParseMode)
	})

	t.Run("fails immediately if context is cancelled", func(t *testing.T) {
		mockClient := &mockTelegramSender{}
		svc := notification.NewTelegramServiceWithClient(mockClient, chatID)

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		err := svc.Send(ctx, "Too late")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "context cancelled")

		// Ensure the mock was never actually called
		assert.Nil(t, mockClient.CapturedChattable)
	})

	t.Run("returns underlying client error", func(t *testing.T) {
		mockClient := &mockTelegramSender{SendErr: errors.New("telegram API down")}
		svc := notification.NewTelegramServiceWithClient(mockClient, chatID)

		err := svc.Send(context.Background(), "Fail me")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to send message")
		assert.Contains(t, err.Error(), "telegram API down")
	})
}
