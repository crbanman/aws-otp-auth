package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
	"gopkg.in/ini.v1"
)

// mockSTSCombinedClient implements STSCombinedClient.
type mockSTSCombinedClient struct {
	CheckValid        bool
	SessionTokenValid bool
	SessionTokenError error
}

func (m *mockSTSCombinedClient) GetCallerIdentity(ctx context.Context, input *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error) {
	if m.CheckValid {
		return &sts.GetCallerIdentityOutput{
			Account: aws.String("123456789012"),
			Arn:     aws.String("arn:aws:iam::123456789012:user/test"),
			UserId:  aws.String("AIDEXAMPLE"),
		}, nil
	}
	return nil, fmt.Errorf("invalid credentials")
}

func (m *mockSTSCombinedClient) GetSessionToken(ctx context.Context, input *sts.GetSessionTokenInput, optFns ...func(*sts.Options)) (*sts.GetSessionTokenOutput, error) {
	if m.SessionTokenError != nil {
		return nil, m.SessionTokenError
	}
	if m.SessionTokenValid {
		return &sts.GetSessionTokenOutput{
			Credentials: &types.Credentials{
				AccessKeyId:     aws.String("newAccessKey"),
				SecretAccessKey: aws.String("newSecretKey"),
				SessionToken:    aws.String("newSessionToken"),
				Expiration:      aws.Time(time.Now().Add(1 * time.Hour)),
			},
		}, nil
	}
	return nil, fmt.Errorf("failed to get session token")
}

func TestIntegrationFlow(t *testing.T) {
	// Set up a temporary HOME directory.
	tempHome := t.TempDir()
	os.Setenv("HOME", tempHome)
	defer os.Unsetenv("HOME")

	awsDir := filepath.Join(tempHome, ".aws")
	if err := os.MkdirAll(awsDir, 0755); err != nil {
		t.Fatalf("Failed to create .aws directory: %v", err)
	}

	// Write an initial (invalid) credentials file.
	credsPath := filepath.Join(awsDir, "credentials")
	initialCreds := `[default]
aws_access_key_id = INVALID
aws_secret_access_key = INVALID
`
	if err := os.WriteFile(credsPath, []byte(initialCreds), 0644); err != nil {
		t.Fatalf("Failed to write credentials file: %v", err)
	}

	// Create a mock STS client that forces re-authentication.
	mockClient := &mockSTSCombinedClient{
		CheckValid:        false, // This will force the flow to obtain new credentials.
		SessionTokenValid: true,
	}

	// Simulate user input for OTP using a strings.Reader.
	otpInput := "654321\n"
	otpReader := strings.NewReader(otpInput)

	// Run the authentication flow with the added MFA ARN argument.
	err := RunAuthFlow(context.Background(), mockClient, otpReader, "default", "", false, true, "dummy-mfa-arn", 28800)
	if err != nil {
		t.Fatalf("RunAuthFlow failed: %v", err)
	}

	// Verify that the credentials file was updated with the new session credentials.
	cfg, err := ini.Load(credsPath)
	if err != nil {
		t.Fatalf("Failed to load updated credentials file: %v", err)
	}
	section, err := cfg.GetSection("default")
	if err != nil {
		t.Fatalf("Profile 'default' not found in credentials file")
	}
	if section.Key("aws_access_key_id").String() != "newAccessKey" {
		t.Errorf("Expected aws_access_key_id to be 'newAccessKey', got %s", section.Key("aws_access_key_id").String())
	}
	if section.Key("aws_secret_access_key").String() != "newSecretKey" {
		t.Errorf("Expected aws_secret_access_key to be 'newSecretKey', got %s", section.Key("aws_secret_access_key").String())
	}
	if section.Key("aws_session_token").String() != "newSessionToken" {
		t.Errorf("Expected aws_session_token to be 'newSessionToken', got %s", section.Key("aws_session_token").String())
	}
}

func TestExpiredTokenFlow(t *testing.T) {
	// Create a temporary HOME directory.
	tempHome := t.TempDir()
	os.Setenv("HOME", tempHome)
	defer os.Unsetenv("HOME")

	awsDir := filepath.Join(tempHome, ".aws")
	if err := os.MkdirAll(awsDir, 0755); err != nil {
		t.Fatalf("Failed to create .aws directory: %v", err)
	}

	// Write an initial credentials file with an expired token.
	credsPath := filepath.Join(awsDir, "credentials")
	expiredTime := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
	initialCreds := `[default]
aws_access_key_id = OLD
aws_secret_access_key = OLDSECRET
aws_session_token = OLDSESSION
aws_session_token_expiration = ` + expiredTime + `
`
	if err := os.WriteFile(credsPath, []byte(initialCreds), 0644); err != nil {
		t.Fatalf("Failed to write credentials file: %v", err)
	}
}

func TestRunAuthFlow_ValidTargetCredentials(t *testing.T) {
	// Create a temporary HOME directory and credentials file for the target profile.
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	defer os.Unsetenv("HOME")
	awsDir := filepath.Join(tempDir, ".aws")
	os.MkdirAll(awsDir, 0755)
	credsPath := filepath.Join(awsDir, "credentials")

	// Write a target credentials file with a valid (non-expired) token.
	validExpiry := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
	content := `[default]
aws_access_key_id = DUMMY
aws_secret_access_key = DUMMYSECRET
aws_session_token = DUMMYTOKEN
aws_session_token_expiration = ` + validExpiry + `
`
	if err := os.WriteFile(credsPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write credentials file: %v", err)
	}

	// Create a mock STS client (not used in this flow because token is valid).
	mockSTS := &mockSTSCombinedClient{CheckValid: true}
	// Run the authentication flow with the added MFA ARN argument.
	err := RunAuthFlow(context.Background(), mockSTS, nil, "default", "", false, true, "dummy-mfa-arn", 28800)
	if err != nil {
		t.Errorf("RunAuthFlow failed when token was valid: %v", err)
	}
}
