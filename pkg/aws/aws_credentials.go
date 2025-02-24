package aws

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/ini.v1"
)

// Credentials holds AWS credentials along with an optional expiration timestamp.
type Credentials struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	Expiration      time.Time
}

// ReadAWSCredentials reads the AWS credentials file from ~/.aws/credentials
// and returns the credentials for the specified profile.
func ReadAWSCredentials(profile string) (*Credentials, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("unable to determine user home directory: %w", err)
	}
	filePath := filepath.Join(home, ".aws", "credentials")
	return readAWSCredentialsFromFile(filePath, profile)
}

// readAWSCredentialsFromFile reads and parses the credentials from the given file path.
func readAWSCredentialsFromFile(filePath, profile string) (*Credentials, error) {
	cfg, err := ini.Load(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load credentials file: %w", err)
	}
	section, err := cfg.GetSection(profile)
	if err != nil {
		return nil, fmt.Errorf("profile %s not found in credentials file", profile)
	}

	cred := &Credentials{
		AccessKeyID:     section.Key("aws_access_key_id").String(),
		SecretAccessKey: section.Key("aws_secret_access_key").String(),
		SessionToken:    section.Key("aws_session_token").String(),
	}
	// Parse expiration time if present.
	if expStr := section.Key("aws_session_token_expiration").String(); expStr != "" {
		if expTime, err := time.Parse(time.RFC3339, expStr); err == nil {
			cred.Expiration = expTime
		}
	}

	// Ensure required fields are present.
	if cred.AccessKeyID == "" || cred.SecretAccessKey == "" {
		return nil, fmt.Errorf("incomplete credentials for profile %s", profile)
	}
	return cred, nil
}

// CleanExpiredTokenFromCredentials checks the specified profile in the credentials file.
// If a session token and its expiration exist and the token is expired, they are removed.
func CleanExpiredTokenFromCredentials(profile string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("unable to determine home directory: %w", err)
	}
	credsPath := filepath.Join(home, ".aws", "credentials")
	cfg, err := ini.Load(credsPath)
	if err != nil {
		return fmt.Errorf("failed to load credentials file: %w", err)
	}
	section, err := cfg.GetSection(profile)
	if err != nil {
		// Profile not found, nothing to clean.
		return nil
	}

	token := section.Key("aws_session_token").String()
	expStr := section.Key("aws_session_token_expiration").String()
	if token != "" && expStr != "" {
		expTime, err := time.Parse(time.RFC3339, expStr)
		if err == nil && time.Now().After(expTime) {
			// Remove expired session token and expiration keys.
			section.DeleteKey("aws_session_token")
			section.DeleteKey("aws_session_token_expiration")
			if err := cfg.SaveTo(credsPath); err != nil {
				return fmt.Errorf("failed to save cleaned credentials: %w", err)
			}
		}
	}
	return nil
}
