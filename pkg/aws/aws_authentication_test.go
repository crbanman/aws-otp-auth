package aws

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type mockSTSClient struct {
	Valid bool
	Err   error
}

func (m *mockSTSClient) GetCallerIdentity(ctx context.Context, input *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	if m.Valid {
		return &sts.GetCallerIdentityOutput{
			Account: aws.String("123456789012"),
			Arn:     aws.String("arn:aws:iam::123456789012:user/test"),
			UserId:  aws.String("AIDEXAMPLE"),
		}, nil
	}
	return nil, errors.New("invalid credentials")
}

func TestCheckAuthentication_Valid(t *testing.T) {
	ctx := context.Background()
	mockClient := &mockSTSClient{Valid: true}

	if err := CheckAuthentication(ctx, mockClient); err != nil {
		t.Fatalf("Expected valid credentials to pass authentication, got error: %v", err)
	}
}

func TestCheckAuthentication_Invalid(t *testing.T) {
	ctx := context.Background()
	mockClient := &mockSTSClient{Valid: false}

	err := CheckAuthentication(ctx, mockClient)
	if err == nil {
		t.Fatal("Expected authentication to fail for invalid credentials, got no error")
	}
	if !strings.Contains(err.Error(), "invalid credentials") {
		t.Fatalf("Expected error message to contain 'invalid credentials', got: %v", err)
	}
}

func TestCheckAuthentication_Error(t *testing.T) {
	ctx := context.Background()
	mockClient := &mockSTSClient{Err: errors.New("network error")}

	err := CheckAuthentication(ctx, mockClient)
	if err == nil {
		t.Fatal("Expected authentication to fail due to network error, got no error")
	}
	if !strings.Contains(err.Error(), "network error") {
		t.Fatalf("Expected error message to contain 'network error', got: %v", err)
	}
}
