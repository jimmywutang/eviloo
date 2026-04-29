package core

import (
	"strings"
	"testing"
)

// TestHandleConfigTelegramBotToken tests setting telegram bot token via config command
func TestHandleConfigTelegramBotToken(t *testing.T) {
	cfg, err := NewConfig("", "")
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	term := &Terminal{
		cfg: cfg,
	}

	// Test setting bot token
	testToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	err = term.handleConfig([]string{"telegram", "bot_token", testToken})
	if err != nil {
		t.Errorf("handleConfig returned error: %v", err)
	}

	// Verify token was set
	if cfg.GetTelegramBotToken() != testToken {
		t.Errorf("Expected bot token %s, got %s", testToken, cfg.GetTelegramBotToken())
	}
}

// TestHandleConfigTelegramChatIDs tests setting telegram chat IDs via config command
func TestHandleConfigTelegramChatIDs(t *testing.T) {
	cfg, err := NewConfig("", "")
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	term := &Terminal{
		cfg: cfg,
	}

	// Test setting chat IDs
	testChatIDs := "123456789,987654321"
	err = term.handleConfig([]string{"telegram", "chat_ids", testChatIDs})
	if err != nil {
		t.Errorf("handleConfig returned error: %v", err)
	}

	// Verify chat IDs were set
	chatIDs := cfg.GetTelegramChatIDs()
	if len(chatIDs) != 2 || chatIDs[0] != "123456789" || chatIDs[1] != "987654321" {
		t.Errorf("Expected chat IDs [%s], got %v", testChatIDs, chatIDs)
	}
}

// TestHandleConfigTelegramAddChatID tests adding a single chat ID
func TestHandleConfigTelegramAddChatID(t *testing.T) {
	cfg, err := NewConfig("", "")
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	term := &Terminal{
		cfg: cfg,
	}

	// First set initial chat IDs
	cfg.SetTelegramChatIDs([]string{"111111111"})

	// Add a chat ID
	testChatID := "222222222"
	err = term.handleConfig([]string{"telegram", "add_chat_id", testChatID})
	if err != nil {
		t.Errorf("handleConfig returned error: %v", err)
	}

	// Verify chat ID was added
	chatIDs := cfg.GetTelegramChatIDs()
	if len(chatIDs) != 2 {
		t.Errorf("Expected 2 chat IDs, got %d", len(chatIDs))
	}
}

// TestHandleConfigTelegramInvalidOption tests error handling for invalid telegram options
func TestHandleConfigTelegramInvalidOption(t *testing.T) {
	cfg, err := NewConfig("", "")
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	term := &Terminal{
		cfg: cfg,
	}

	// Test invalid option
	err = term.handleConfig([]string{"telegram", "invalid_option", "value"})
	if err == nil {
		t.Error("Expected error for invalid telegram option, got nil")
	}

	if !strings.Contains(err.Error(), "invalid telegram config option") {
		t.Errorf("Expected error message to contain 'invalid telegram config option', got: %v", err)
	}
}

// TestHandleConfigDisplayIncludesTelegram tests that config display includes telegram settings
func TestHandleConfigDisplayIncludesTelegram(t *testing.T) {
	cfg, err := NewConfig("", "")
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	term := &Terminal{
		cfg: cfg,
	}

	// Set telegram values
	testToken := "test_token_123"
	testChatIDs := "123456789,987654321"
	cfg.SetTelegramBotToken(testToken)
	cfg.SetTelegramChatIDs(strings.Split(testChatIDs, ","))

	// Call handleConfig with no args (display mode)
	// This will print to log, so we just verify it doesn't error
	err = term.handleConfig([]string{})
	if err != nil {
		t.Errorf("handleConfig display returned error: %v", err)
	}

	// Verify values are still set (they should be retrievable)
	if cfg.GetTelegramBotToken() != testToken {
		t.Errorf("Expected bot token %s, got %s", testToken, cfg.GetTelegramBotToken())
	}
	chatIDs := cfg.GetTelegramChatIDs()
	if len(chatIDs) != 2 || chatIDs[0] != "123456789" || chatIDs[1] != "987654321" {
		t.Errorf("Expected chat IDs, got %v", chatIDs)
	}
}