package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// SessionCredentials holds temporary AWS session credentials along with their expiration.
type SessionCredentials struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	Expiration      time.Time
}

// STSGetSessionTokenClient defines the subset of the AWS STS client's methods needed to call GetSessionToken.
type STSGetSessionTokenClient interface {
	GetSessionToken(ctx context.Context, params *sts.GetSessionTokenInput, optFns ...func(*sts.Options)) (*sts.GetSessionTokenOutput, error)
}

// GetSessionToken calls AWS STS's GetSessionToken API using the provided MFA ARN, OTP code, and desired session duration.
// It returns the new session credentials on success.
func GetSessionToken(ctx context.Context, client STSGetSessionTokenClient, mfaArn, tokenCode string, durationSeconds int32) (*SessionCredentials, error) {
	input := &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int32(durationSeconds),
		SerialNumber:    aws.String(mfaArn),
		TokenCode:       aws.String(tokenCode),
	}
	result, err := client.GetSessionToken(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get session token: %w", err)
	}
	if result.Credentials == nil {
		return nil, fmt.Errorf("no credentials returned")
	}
	creds := result.Credentials
	return &SessionCredentials{
		AccessKeyID:     aws.ToString(creds.AccessKeyId),
		SecretAccessKey: aws.ToString(creds.SecretAccessKey),
		SessionToken:    aws.ToString(creds.SessionToken),
		Expiration:      aws.ToTime(creds.Expiration),
	}, nil
}
