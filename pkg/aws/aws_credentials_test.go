package aws

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"gopkg.in/ini.v1"
)

func TestReadAWSCredentialsFromFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "credentials")
	content := `[default]
aws_access_key_id = ABCDEFGHIJKLMNOP
aws_secret_access_key = abcdefghijklmnopqrstuvwxyz1234567890

[myprofile]
aws_access_key_id = MYACCESSKEY
aws_secret_access_key = mysecretkey
aws_session_token = mysessiontoken
`
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write temp credentials file: %v", err)
	}

	creds, err := readAWSCredentialsFromFile(filePath, "default")
	if err != nil {
		t.Errorf("Expected no error for default profile, got %v", err)
	}
	if creds.AccessKeyID != "ABCDEFGHIJKLMNOP" {
		t.Errorf("Expected AccessKeyID 'ABCDEFGHIJKLMNOP', got '%s'", creds.AccessKeyID)
	}
	if creds.SecretAccessKey != "abcdefghijklmnopqrstuvwxyz1234567890" {
		t.Errorf("Expected SecretAccessKey 'abcdefghijklmnopqrstuvwxyz1234567890', got '%s'", creds.SecretAccessKey)
	}
	if creds.SessionToken != "" {
		t.Errorf("Expected empty SessionToken for default profile, got '%s'", creds.SessionToken)
	}

	creds, err = readAWSCredentialsFromFile(filePath, "myprofile")
	if err != nil {
		t.Errorf("Expected no error for myprofile, got %v", err)
	}
	if creds.AccessKeyID != "MYACCESSKEY" {
		t.Errorf("Expected AccessKeyID 'MYACCESSKEY', got '%s'", creds.AccessKeyID)
	}
	if creds.SecretAccessKey != "mysecretkey" {
		t.Errorf("Expected SecretAccessKey 'mysecretkey', got '%s'", creds.SecretAccessKey)
	}
	if creds.SessionToken != "mysessiontoken" {
		t.Errorf("Expected SessionToken 'mysessiontoken', got '%s'", creds.SessionToken)
	}

	_, err = readAWSCredentialsFromFile(filePath, "nonexistent")
	if err == nil {
		t.Errorf("Expected error for non-existent profile, got nil")
	}

	malformedPath := filepath.Join(tempDir, "malformed")
	badContent := "This is not a valid INI content"
	if err := os.WriteFile(malformedPath, []byte(badContent), 0644); err != nil {
		t.Fatalf("Failed to write temp malformed file: %v", err)
	}
	_, err = readAWSCredentialsFromFile(malformedPath, "default")
	if err == nil {
		t.Errorf("Expected error for malformed file, got nil")
	}

	incompletePath := filepath.Join(tempDir, "incomplete")
	incompleteContent := `[default]
aws_access_key_id = SOMEKEY
`
	if err := os.WriteFile(incompletePath, []byte(incompleteContent), 0644); err != nil {
		t.Fatalf("Failed to write temp incomplete file: %v", err)
	}
	_, err = readAWSCredentialsFromFile(incompletePath, "default")
	if err == nil {
		t.Errorf("Expected error for incomplete credentials, got nil")
	}
}

func TestCleanExpiredToken(t *testing.T) {
	// Create a temporary HOME directory.
	tempDir := t.TempDir()
	awsDir := filepath.Join(tempDir, ".aws")
	if err := os.MkdirAll(awsDir, 0755); err != nil {
		t.Fatalf("Failed to create .aws directory: %v", err)
	}
	credsPath := filepath.Join(awsDir, "credentials")

	// Write a credentials file with an expired token.
	expiredTime := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
	content := `[default]
aws_access_key_id = ABCD
aws_secret_access_key = SECRET
aws_session_token = EXPIREDTOKEN
aws_session_token_expiration = ` + expiredTime + `
`
	if err := os.WriteFile(credsPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write credentials file: %v", err)
	}

	// Override HOME so os.UserHomeDir() returns tempDir.
	origHome := os.Getenv("HOME")
	defer os.Setenv("HOME", origHome)
	os.Setenv("HOME", tempDir)

	// Call the cleaning function.
	if err := CleanExpiredTokenFromCredentials("default"); err != nil {
		t.Fatalf("CleanExpiredTokenFromCredentials failed: %v", err)
	}

	// Reload the credentials file and verify the keys are removed.
	cfg, err := ini.Load(credsPath)
	if err != nil {
		t.Fatalf("Failed to load cleaned credentials file: %v", err)
	}
	section, err := cfg.GetSection("default")
	if err != nil {
		t.Fatalf("Default section not found")
	}
	if section.HasKey("aws_session_token") {
		t.Errorf("aws_session_token was not removed")
	}
	if section.HasKey("aws_session_token_expiration") {
		t.Errorf("aws_session_token_expiration was not removed")
	}
}
