package notifications

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// OrderInfo contains data needed for notification text.
type OrderInfo struct {
	OrderID       int64
	Status        string
	Total         float64
	CustomerName  string
	CustomerPhone string
}

// TelegramNotifier sends messages to Telegram chats.
type TelegramNotifier struct {
	botToken    string
	client      *http.Client
	adminChatID string
}

// TelegramConfig holds notifier settings.
type TelegramConfig struct {
	BotToken    string
	AdminChatID string
}

// NewTelegramNotifier builds notifier.
func NewTelegramNotifier(cfg TelegramConfig) *TelegramNotifier {
	return &TelegramNotifier{
		botToken:    cfg.BotToken,
		client:      &http.Client{},
		adminChatID: cfg.AdminChatID,
	}
}

// NotifyOrderCreated notifies user/admin about new order.
func (n *TelegramNotifier) NotifyOrderCreated(ctx context.Context, info OrderInfo, userChatID int64) {
	if n.botToken == "" {
		return
	}
	msg := fmt.Sprintf("Новый заказ #%d от %s (%s) на сумму %.2f", info.OrderID, info.CustomerName, info.CustomerPhone, info.Total)
	if n.adminChatID != "" {
		n.sendMessage(ctx, n.adminChatID, msg)
	}
	if userChatID != 0 {
		n.sendMessage(ctx, strconv.FormatInt(userChatID, 10), fmt.Sprintf("Ваш заказ #%d принят. Статус: %s", info.OrderID, strings.Title(info.Status)))
	}
}

// NotifyStatusChanged notifies user about status update.
func (n *TelegramNotifier) NotifyStatusChanged(ctx context.Context, info OrderInfo, userChatID int64) {
	if n.botToken == "" {
		return
	}
	msg := fmt.Sprintf("Статус заказа #%d изменён на %s", info.OrderID, info.Status)
	if userChatID != 0 {
		n.sendMessage(ctx, strconv.FormatInt(userChatID, 10), msg)
	}
	if n.adminChatID != "" {
		n.sendMessage(ctx, n.adminChatID, msg)
	}
}

func (n *TelegramNotifier) sendMessage(ctx context.Context, chatID, message string) {
	if chatID == "" || n.botToken == "" {
		return
	}
	endpoint := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", n.botToken)
	data := url.Values{}
	data.Set("chat_id", chatID)
	data.Set("text", message)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := n.client.Do(req)
	if err != nil {
		return
	}
	resp.Body.Close()
}
