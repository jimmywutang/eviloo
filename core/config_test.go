package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTelegramConfigSettersAndGetters(t *testing.T) {
	// Create a temporary directory for the test config
	tmpDir, err := os.MkdirTemp("", "evilginx-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a new config
	configPath := filepath.Join(tmpDir, "config.json")
	cfg, err := NewConfig(tmpDir, configPath)
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	// Test SetTelegramBotToken and GetTelegramBotToken
	testToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	cfg.SetTelegramBotToken(testToken)
	
	gotToken := cfg.GetTelegramBotToken()
	if gotToken != testToken {
		t.Errorf("GetTelegramBotToken() = %v, want %v", gotToken, testToken)
	}

	// Test SetTelegramChatIDs and GetTelegramChatIDs
	testChatIDs := []string{"-1001234567890", "-1009876543210"}
	cfg.SetTelegramChatIDs(testChatIDs)
	
	gotChatIDs := cfg.GetTelegramChatIDs()
	if len(gotChatIDs) != len(testChatIDs) || gotChatIDs[0] != testChatIDs[0] || gotChatIDs[1] != testChatIDs[1] {
		t.Errorf("GetTelegramChatIDs() = %v, want %v", gotChatIDs, testChatIDs)
	}
}

func TestTelegramConfigPersistence(t *testing.T) {
	// Create a temporary directory for the test config
	tmpDir, err := os.MkdirTemp("", "evilginx-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a new config
	configPath := filepath.Join(tmpDir, "config.json")
	cfg, err := NewConfig(tmpDir, configPath)
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	// Set telegram config values
	testToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	testChatIDs := []string{"-1001234567890", "-1009876543210"}
	cfg.SetTelegramBotToken(testToken)
	cfg.SetTelegramChatIDs(testChatIDs)

	// Load the config again from the same file
	cfg2, err := NewConfig(tmpDir, configPath)
	if err != nil {
		t.Fatalf("Failed to reload config: %v", err)
	}

	// Verify the values persisted
	gotToken := cfg2.GetTelegramBotToken()
	if gotToken != testToken {
		t.Errorf("After reload, GetTelegramBotToken() = %v, want %v", gotToken, testToken)
	}

	gotChatIDs := cfg2.GetTelegramChatIDs()
	if len(gotChatIDs) != len(testChatIDs) || gotChatIDs[0] != testChatIDs[0] || gotChatIDs[1] != testChatIDs[1] {
		t.Errorf("After reload, GetTelegramChatIDs() = %v, want %v", gotChatIDs, testChatIDs)
	}
}
