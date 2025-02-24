package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// STSClient defines the subset of the AWS STS client's methods used by CheckAuthentication.
type STSClient interface {
	GetCallerIdentity(ctx context.Context, params *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error)
}

// CheckAuthentication calls AWS STS's GetCallerIdentity API to validate the current credentials.
// It returns nil if the credentials are valid or an error if the call fails.
func CheckAuthentication(ctx context.Context, client STSClient) error {
	_, err := client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return fmt.Errorf("authentication check failed: %w", err)
	}
	return nil
}
