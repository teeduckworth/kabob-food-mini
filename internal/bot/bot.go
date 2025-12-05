package bot

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
	"unicode"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Config collects dependencies required to start the bot loop.
type Config struct {
	Token      string
	BackendURL string
	MiniAppURL string
	Debug      bool
	HTTPClient *http.Client
}

// Bot drives Telegram interactions and backend registration requests.
type Bot struct {
	api        *tgbotapi.BotAPI
	backendURL string
	miniAppURL string
	httpClient *http.Client
	sessions   map[int64]*session
	mu         sync.Mutex
}

// New creates a configured Bot instance.
func New(cfg Config) (*Bot, error) {
	if strings.TrimSpace(cfg.Token) == "" {
		return nil, errors.New("telegram bot token is required")
	}
	backend := strings.TrimRight(strings.TrimSpace(cfg.BackendURL), "/")
	if backend == "" {
		backend = "http://localhost:8080"
	}
	mini := strings.TrimSpace(cfg.MiniAppURL)
	if mini == "" {
		mini = "https://kabob-food-mini.vercel.app"
	}
	api, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, err
	}
	api.Debug = cfg.Debug
	client := cfg.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return &Bot{
		api:        api,
		backendURL: backend,
		miniAppURL: mini,
		httpClient: client,
		sessions:   make(map[int64]*session),
	}, nil
}

// Run consumes Telegram updates until the provided context is cancelled.
func (b *Bot) Run(ctx context.Context) error {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := b.api.GetUpdatesChan(updateConfig)
	for {
		select {
		case <-ctx.Done():
			b.api.StopReceivingUpdates()
			return ctx.Err()
		case update := <-updates:
			b.handleUpdate(update)
		}
	}
}

func (b *Bot) handleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}
	msg := update.Message
	if msg.From == nil {
		return
	}
	userID := msg.From.ID
	if msg.IsCommand() {
		b.handleCommand(msg)
		return
	}
	if msg.Contact != nil {
		b.handleContact(msg)
		return
	}
	if msg.Location != nil {
		b.handleLocation(msg)
		return
	}
	if msg.Text != "" {
		b.handleText(msg)
		return
	}
	b.reply(msg.Chat.ID, "Пока не понимаю это сообщение. Нажмите /start, чтобы начать заново.")
	b.resetSession(userID)
}

func (b *Bot) handleCommand(msg *tgbotapi.Message) {
	switch msg.Command() {
	case "start":
		b.resetSession(msg.From.ID)
		b.promptContact(msg.Chat.ID)
	default:
		b.reply(msg.Chat.ID, "Команда не поддерживается. Отправьте /start, чтобы начать регистрацию.")
	}
}

func (b *Bot) handleContact(msg *tgbotapi.Message) {
	sess := b.ensureSession(msg.From.ID)
	sess.Phone = strings.TrimSpace(msg.Contact.PhoneNumber)
	sess.FirstName = strings.TrimSpace(msg.Contact.FirstName)
	sess.LastName = strings.TrimSpace(msg.Contact.LastName)
	if sess.Phone == "" {
		b.reply(msg.Chat.ID, "Не удалось прочитать номер телефона. Попробуйте снова отправить контакт.")
		return
	}
	b.askForName(msg.Chat.ID, sess)
}

func (b *Bot) handleText(msg *tgbotapi.Message) {
	sess := b.ensureSession(msg.From.ID)
	switch sess.Stage {
	case stageNeedContact:
		phone := normalizePhone(msg.Text)
		if phone == "" {
			b.reply(msg.Chat.ID, "Пожалуйста, отправьте номер телефона (можно вручную) или воспользуйтесь кнопкой.")
			return
		}
		sess.Phone = phone
		if sess.FirstName == "" {
			sess.FirstName = strings.TrimSpace(msg.From.FirstName)
		}
		if sess.LastName == "" {
			sess.LastName = strings.TrimSpace(msg.From.LastName)
		}
		b.askForName(msg.Chat.ID, sess)
	case stageNeedName:
		name := strings.TrimSpace(msg.Text)
		if name == "" {
			b.reply(msg.Chat.ID, "Пожалуйста, введите имя или напишите /start, чтобы начать заново.")
			return
		}
		sess.CustomName = name
		sess.Stage = stageNeedLocation
		b.promptLocation(msg.Chat.ID)
	default:
		b.reply(msg.Chat.ID, "Отправьте /start, чтобы начать регистрацию.")
	}
}

func (b *Bot) handleLocation(msg *tgbotapi.Message) {
	sess := b.ensureSession(msg.From.ID)
	if sess.Stage != stageNeedLocation {
		b.reply(msg.Chat.ID, "Сначала отправьте контакт и имя. Нажмите /start, чтобы начать заново.")
		return
	}
	sess.Latitude = msg.Location.Latitude
	sess.Longitude = msg.Location.Longitude
	sess.Stage = stageReady
	if err := b.registerUser(msg); err != nil {
		b.reply(msg.Chat.ID, fmt.Sprintf("Не удалось завершить регистрацию: %v", err))
		return
	}
	deleteReply := tgbotapi.NewRemoveKeyboard(true)
	send := tgbotapi.NewMessage(msg.Chat.ID, "Отлично! Нажмите кнопку ниже, чтобы открыть мини-апп.")
	send.ReplyMarkup = deleteReply
	b.api.Send(send)
}

func (b *Bot) registerUser(msg *tgbotapi.Message) error {
	sess := b.ensureSession(msg.From.ID)
	name := sess.CustomName
	if name == "" {
		name = sess.FirstName
	}
	if name == "" {
		return errors.New("не указано имя")
	}
	if sess.Phone == "" {
		return errors.New("не указан телефон")
	}
	reqBody := botRegisterRequest{
		TelegramID: msg.From.ID,
		Phone:      sess.Phone,
		FirstName:  name,
		LastName:   sess.LastName,
		Name:       name,
		Location: locationPayload{
			Latitude:  sess.Latitude,
			Longitude: sess.Longitude,
		},
	}
	data, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	endpoint := fmt.Sprintf("%s/bot/register", b.backendURL)
	resp, err := b.httpClient.Post(endpoint, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("backend ответил %d", resp.StatusCode)
	}
	var regResp registerResponse
	if err := json.NewDecoder(resp.Body).Decode(&regResp); err != nil {
		return err
	}
	if regResp.Token == "" {
		return errors.New("получен пустой токен")
	}
	link := b.buildMiniAppLink(regResp.Token)
	message := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Если кнопка не появилась, воспользуйтесь ссылкой: %s", link))
	message.ReplyMarkup = newMiniAppKeyboard(link)
	if _, err := b.api.Send(message); err != nil {
		return err
	}
	return nil
}

func (b *Bot) promptContact(chatID int64) {
	button := tgbotapi.KeyboardButton{Text: "Отправить телефон", RequestContact: true}
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(button),
	)
	keyboard.ResizeKeyboard = true
	keyboard.OneTimeKeyboard = true
	msg := tgbotapi.NewMessage(chatID, "Привет! Нажмите кнопку ниже, чтобы поделиться номером телефона.")
	msg.ReplyMarkup = keyboard
	b.api.Send(msg)
}

func (b *Bot) promptLocation(chatID int64) {
	button := tgbotapi.KeyboardButton{Text: "Поделиться геолокацией", RequestLocation: true}
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(button),
	)
	keyboard.ResizeKeyboard = true
	keyboard.OneTimeKeyboard = true
	msg := tgbotapi.NewMessage(chatID, "Спасибо! Осталось отправить локацию, чтобы мы знали, куда доставлять.")
	msg.ReplyMarkup = keyboard
	b.api.Send(msg)
}

func (b *Bot) reply(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	b.api.Send(msg)
}

func (b *Bot) askForName(chatID int64, sess *session) {
	sess.Stage = stageNeedName
	text := "Спасибо! Теперь напишите, как к вам обращаться?"
	b.reply(chatID, text)
}

func (b *Bot) buildMiniAppLink(token string) string {
	sep := "?"
	if strings.Contains(b.miniAppURL, "?") {
		sep = "&"
	}
	return fmt.Sprintf("%s%stoken=%s", b.miniAppURL, sep, token)
}

func (b *Bot) ensureSession(userID int64) *session {
	b.mu.Lock()
	defer b.mu.Unlock()
	sess, ok := b.sessions[userID]
	if !ok {
		sess = &session{Stage: stageNeedContact}
		b.sessions[userID] = sess
	}
	return sess
}

func (b *Bot) resetSession(userID int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.sessions[userID] = &session{Stage: stageNeedContact}
}

type session struct {
	Stage      stage
	Phone      string
	FirstName  string
	LastName   string
	CustomName string
	Latitude   float64
	Longitude  float64
}

type stage int

const (
	stageNeedContact stage = iota
	stageNeedName
	stageNeedLocation
	stageReady
)

type botRegisterRequest struct {
	TelegramID int64           `json:"telegram_id"`
	Phone      string          `json:"phone"`
	FirstName  string          `json:"first_name"`
	LastName   string          `json:"last_name"`
	Name       string          `json:"name"`
	Location   locationPayload `json:"location"`
}

type locationPayload struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type registerResponse struct {
	Token string `json:"token"`
}

func normalizePhone(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	var b strings.Builder
	var wroteDigits bool
	for _, r := range raw {
		if r == '+' && !wroteDigits && b.Len() == 0 {
			b.WriteRune(r)
			continue
		}
		if unicode.IsDigit(r) {
			b.WriteRune(r)
			wroteDigits = true
		}
	}
	phone := b.String()
	if phone == "" || phone == "+" {
		return ""
	}
	return phone
}

// inline keyboard payload with web_app support (older tgbotapi lacks this type).
type miniAppKeyboard struct {
	InlineKeyboard [][]miniAppButton `json:"inline_keyboard"`
}

type miniAppButton struct {
	Text   string          `json:"text"`
	WebApp *miniAppWebInfo `json:"web_app,omitempty"`
	URL    string          `json:"url,omitempty"`
}

type miniAppWebInfo struct {
	URL string `json:"url"`
}

func newMiniAppKeyboard(link string) miniAppKeyboard {
	return miniAppKeyboard{
		InlineKeyboard: [][]miniAppButton{
			{
				{Text: "Открыть мини-апп", WebApp: &miniAppWebInfo{URL: link}},
			},
			{
				{Text: "Открыть в браузере", URL: link},
			},
		},
	}
}
