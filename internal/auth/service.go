package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/rashidmailru/kabobfood/internal/users"
)

// Service handles Telegram WebApp auth + JWT issuance.
type Service struct {
	userRepo    UserRepository
	botToken    string
	jwtSecret   []byte
	jwtExpiry   time.Duration
	initDataTTL time.Duration
}

// Config gathers dependencies for Service.
type Config struct {
	UserRepo    UserRepository
	BotToken    string
	JWTSecret   string
	JWTExpiry   time.Duration
	InitDataTTL time.Duration
}

// UserRepository abstracts user persistence for auth service.
type UserRepository interface {
	UpsertTelegramUser(ctx context.Context, input users.UpsertTelegramUserInput) (*users.User, error)
}

var (
	// ErrInvalidInitData indicates the initData string could not be parsed or verified.
	ErrInvalidInitData = errors.New("invalid telegram init data")
	// ErrExpiredInitData indicates the payload is too old based on config.
	ErrExpiredInitData = errors.New("telegram init data expired")
	// ErrMissingUserPayload indicates the user JSON chunk is missing.
	ErrMissingUserPayload = errors.New("telegram user payload missing")
)

// NewService builds a Service instance.
func NewService(cfg Config) (*Service, error) {
	if cfg.UserRepo == nil {
		return nil, errors.New("nil user repository")
	}
	if cfg.BotToken == "" {
		return nil, errors.New("telegram bot token is required")
	}
	if cfg.JWTSecret == "" {
		return nil, errors.New("jwt secret is required")
	}
	if cfg.JWTExpiry <= 0 {
		cfg.JWTExpiry = 24 * time.Hour
	}
	if cfg.InitDataTTL <= 0 {
		cfg.InitDataTTL = time.Hour
	}

	return &Service{
		userRepo:    cfg.UserRepo,
		botToken:    cfg.BotToken,
		jwtSecret:   []byte(cfg.JWTSecret),
		jwtExpiry:   cfg.JWTExpiry,
		initDataTTL: cfg.InitDataTTL,
	}, nil
}

// AuthResult contains a JWT token plus user profile.
type AuthResult struct {
	Token   string      `json:"token"`
	Profile *users.User `json:"profile"`
}

// Authenticate verifies Telegram initData, upserts user, and issues JWT.
func (s *Service) Authenticate(ctx context.Context, initData string) (*AuthResult, error) {
	values, err := url.ParseQuery(initData)
	if err != nil {
		return nil, ErrInvalidInitData
	}

	if err := s.verifyPayload(values); err != nil {
		return nil, err
	}

	userJSON := values.Get("user")
	if userJSON == "" {
		return nil, ErrMissingUserPayload
	}

	var telegramUser telegramUserPayload
	if err := json.Unmarshal([]byte(userJSON), &telegramUser); err != nil {
		return nil, ErrInvalidInitData
	}

	profile, err := s.userRepo.UpsertTelegramUser(ctx, users.UpsertTelegramUserInput{
		TelegramID: telegramUser.ID,
		FirstName:  telegramUser.FirstName,
		LastName:   telegramUser.LastName,
		Username:   telegramUser.Username,
		Phone:      telegramUser.PhoneNumber,
		Language:   telegramUser.LanguageCode,
	})
	if err != nil {
		return nil, err
	}

	token, err := s.issueJWT(profile)
	if err != nil {
		return nil, err
	}

	return &AuthResult{Token: token, Profile: profile}, nil
}

func (s *Service) verifyPayload(values url.Values) error {
	hash := values.Get("hash")
	if hash == "" {
		return ErrInvalidInitData
	}

	authDateStr := values.Get("auth_date")
	if authDateStr == "" {
		return ErrInvalidInitData
	}

	authUnix, err := strconv.ParseInt(authDateStr, 10, 64)
	if err != nil {
		return ErrInvalidInitData
	}

	if s.initDataTTL > 0 {
		authTime := time.Unix(authUnix, 0)
		if time.Since(authTime) > s.initDataTTL {
			return ErrExpiredInitData
		}
	}

	dataCheckString := buildDataCheckString(values)
	secret := sha256.Sum256([]byte(s.botToken))

	mac := hmac.New(sha256.New, secret[:])
	mac.Write([]byte(dataCheckString))
	expected := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(expected), []byte(hash)) {
		return ErrInvalidInitData
	}

	return nil
}

func (s *Service) issueJWT(user *users.User) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":         user.ID,
		"telegram_id": user.TelegramID,
		"role":        "user",
		"iat":         now.Unix(),
		"exp":         now.Add(s.jwtExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func buildDataCheckString(values url.Values) string {
	pairs := make([]string, 0, len(values))
	for key, vals := range values {
		if key == "hash" {
			continue
		}
		if len(vals) == 0 {
			continue
		}
		pairs = append(pairs, key+"="+vals[0])
	}
	sort.Strings(pairs)
	return strings.Join(pairs, "\n")
}

// telegramUserPayload mirrors Telegram user JSON embedded inside initData.
type telegramUserPayload struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
	PhoneNumber  string `json:"phone_number"`
}
