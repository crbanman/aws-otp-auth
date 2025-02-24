package aws

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/ini.v1"
)

// UpdateCredentials backs up the current credentials file and updates the specified profile
// with the new session credentials.
func UpdateCredentials(profile string, newCreds *SessionCredentials) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("unable to determine user home directory: %w", err)
	}
	credsPath := filepath.Join(home, ".aws", "credentials")
	backupPath := filepath.Join(home, ".aws", "credentials.bak")

	if err := copyFile(credsPath, backupPath); err != nil {
		return fmt.Errorf("failed to backup credentials file: %w", err)
	}

	cfg, err := ini.Load(credsPath)
	if err != nil {
		return fmt.Errorf("failed to load credentials file: %w", err)
	}

	section, err := cfg.GetSection(profile)
	if err != nil {
		section, err = cfg.NewSection(profile)
		if err != nil {
			return fmt.Errorf("failed to create profile section: %w", err)
		}
	}

	section.Key("aws_access_key_id").SetValue(newCreds.AccessKeyID)
	section.Key("aws_secret_access_key").SetValue(newCreds.SecretAccessKey)
	section.Key("aws_session_token").SetValue(newCreds.SessionToken)
	section.Key("aws_session_token_expiration").SetValue(newCreds.Expiration.Format(time.RFC3339))

	if err := cfg.SaveTo(credsPath); err != nil {
		return fmt.Errorf("failed to save updated credentials file: %w", err)
	}

	return nil
}

// copyFile copies a file from src to dst.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Close()
	}()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}
