package aws

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"gopkg.in/ini.v1"
)

func TestUpdateCredentials(t *testing.T) {
	// Create a temporary directory to simulate the user's HOME.
	tempDir := t.TempDir()
	awsDir := filepath.Join(tempDir, ".aws")
	if err := os.MkdirAll(awsDir, 0755); err != nil {
		t.Fatalf("Failed to create .aws directory: %v", err)
	}

	// Path to the simulated credentials file.
	credsPath := filepath.Join(awsDir, "credentials")
	initialContent := `[default]
aws_access_key_id = OLDACCESSKEY
aws_secret_access_key = OLDSECRETKEY

[profile1]
aws_access_key_id = PROFILE1OLDKEY
aws_secret_access_key = PROFILE1OLDSECRET
aws_session_token = PROFILE1OLDTOKEN
`
	if err := os.WriteFile(credsPath, []byte(initialContent), 0644); err != nil {
		t.Fatalf("Failed to write test credentials file: %v", err)
	}

	// Create new session credentials for the test.
	newCreds := &SessionCredentials{
		AccessKeyID:     "NEWACCESSKEY",
		SecretAccessKey: "NEWSECRETKEY",
		SessionToken:    "NEWSESSIONTOKEN",
		Expiration:      time.Now().Add(1 * time.Hour).Truncate(time.Second), // Truncate to avoid minor differences in formatting
	}

	// Override HOME so that os.UserHomeDir() returns tempDir.
	origHome := os.Getenv("HOME")
	defer os.Setenv("HOME", origHome)
	os.Setenv("HOME", tempDir)

	if err := UpdateCredentials("profile1", newCreds); err != nil {
		t.Fatalf("UpdateCredentials returned error: %v", err)
	}

	backupPath := filepath.Join(awsDir, "credentials.bak")
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		t.Fatalf("Backup file was not created")
	}

	backupContent, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatalf("Failed to read backup file: %v", err)
	}
	if string(backupContent) != initialContent {
		t.Errorf("Backup content mismatch.\nExpected:\n%s\nGot:\n%s", initialContent, string(backupContent))
	}

	cfg, err := ini.Load(credsPath)
	if err != nil {
		t.Fatalf("Failed to load updated credentials file: %v", err)
	}

	section, err := cfg.GetSection("profile1")
	if err != nil {
		t.Fatalf("Profile 'profile1' not found in updated file")
	}
	if section.Key("aws_access_key_id").String() != "NEWACCESSKEY" {
		t.Errorf("aws_access_key_id not updated. Got: %s", section.Key("aws_access_key_id").String())
	}
	if section.Key("aws_secret_access_key").String() != "NEWSECRETKEY" {
		t.Errorf("aws_secret_access_key not updated. Got: %s", section.Key("aws_secret_access_key").String())
	}
	if section.Key("aws_session_token").String() != "NEWSESSIONTOKEN" {
		t.Errorf("aws_session_token not updated. Got: %s", section.Key("aws_session_token").String())
	}
	// Check the expiration key.
	expStr := section.Key("aws_session_token_expiration").String()
	if expStr == "" {
		t.Errorf("aws_session_token_expiration key is missing")
	} else {
		parsedExp, err := time.Parse(time.RFC3339, expStr)
		if err != nil {
			t.Errorf("Failed to parse expiration time: %v", err)
		} else if !parsedExp.Equal(newCreds.Expiration) {
			t.Errorf("Expiration not updated correctly. Expected %v, got %v", newCreds.Expiration, parsedExp)
		}
	}
}
