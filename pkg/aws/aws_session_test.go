package aws

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
)

type mockSTSGetSessionTokenClient struct {
	Valid bool
	Err   error
}

func (m *mockSTSGetSessionTokenClient) GetSessionToken(ctx context.Context, input *sts.GetSessionTokenInput, optFns ...func(*sts.Options)) (*sts.GetSessionTokenOutput, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	if m.Valid {
		return &sts.GetSessionTokenOutput{
			Credentials: &types.Credentials{
				AccessKeyId:     aws.String("newAccessKey"),
				SecretAccessKey: aws.String("newSecretKey"),
				SessionToken:    aws.String("newSessionToken"),
				Expiration:      aws.Time(time.Now().Add(1 * time.Hour)),
			},
		}, nil
	}
	return nil, errors.New("failed to get session token")
}

func TestGetSessionToken_Success(t *testing.T) {
	ctx := context.Background()
	mfaArn := "arn:aws:iam::123456789012:mfa/user"
	otp := "123456"
	duration := int32(3600)

	mockClient := &mockSTSGetSessionTokenClient{Valid: true}
	creds, err := GetSessionToken(ctx, mockClient, mfaArn, otp, duration)
	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}
	if creds.AccessKeyID != "newAccessKey" {
		t.Errorf("Expected AccessKeyID 'newAccessKey', got '%s'", creds.AccessKeyID)
	}
	if creds.SecretAccessKey != "newSecretKey" {
		t.Errorf("Expected SecretAccessKey 'newSecretKey', got '%s'", creds.SecretAccessKey)
	}
	if creds.SessionToken != "newSessionToken" {
		t.Errorf("Expected SessionToken 'newSessionToken', got '%s'", creds.SessionToken)
	}
	if time.Until(creds.Expiration) <= 0 {
		t.Errorf("Expected Expiration to be in the future, got %v", creds.Expiration)
	}
}

func TestGetSessionToken_Failure(t *testing.T) {
	ctx := context.Background()
	mfaArn := "arn:aws:iam::123456789012:mfa/user"
	otp := "123456"
	duration := int32(3600)

	mockClientErr := &mockSTSGetSessionTokenClient{Err: fmt.Errorf("network error")}
	_, err := GetSessionToken(ctx, mockClientErr, mfaArn, otp, duration)
	if err == nil || err.Error() != "failed to get session token: network error" {
		t.Errorf("Expected network error, got %v", err)
	}

	mockClientInvalid := &mockSTSGetSessionTokenClient{Valid: false}
	_, err = GetSessionToken(ctx, mockClientInvalid, mfaArn, otp, duration)
	if err == nil {
		t.Errorf("Expected failure due to invalid response, got nil")
	} else if !strings.Contains(err.Error(), "failed to get session token") {
		t.Errorf("Expected error message to contain 'failed to get session token', got %v", err)
	}
}
