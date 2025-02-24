package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/user"
	"time"

	"github.com/crbanman/aws-otp-auth/pkg/aws"
	"github.com/crbanman/aws-otp-auth/pkg/otp"

	awsPkg "github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	awsIam "github.com/aws/aws-sdk-go-v2/service/iam"
	awsSts "github.com/aws/aws-sdk-go-v2/service/sts"
)

// STSCombinedClient combines the methods needed for both authentication and session token retrieval.
type STSCombinedClient interface {
	GetCallerIdentity(ctx context.Context, input *awsSts.GetCallerIdentityInput, optFns ...func(*awsSts.Options)) (*awsSts.GetCallerIdentityOutput, error)
	GetSessionToken(ctx context.Context, input *awsSts.GetSessionTokenInput, optFns ...func(*awsSts.Options)) (*awsSts.GetSessionTokenOutput, error)
}

// RunAuthFlow performs the complete authentication flow.
// It reads the target profile’s credentials and if the token is present and not expired, it exits early.
// RunAuthFlow performs the complete authentication flow.
func RunAuthFlow(ctx context.Context, stsClient STSCombinedClient, inReader io.Reader, profile string, providedOTP string, force bool, verbose bool, mfaArn string) error {
	// Read current target credentials.
	creds, err := aws.ReadAWSCredentials(profile)
	if err != nil && verbose {
		fmt.Printf("Warning: failed to read credentials: %v\n", err)
	}

	// If credentials have an expiration and are still valid, exit.
	if creds != nil && !creds.Expiration.IsZero() && time.Now().Before(creds.Expiration) {
		fmt.Println("Existing credentials are valid. No update necessary.")
		return nil
	}

	// If force is not set and token is still valid, exit.
	if !force && creds != nil && !creds.Expiration.IsZero() && time.Now().Before(creds.Expiration) {
		if verbose {
			fmt.Println("Existing credentials are valid. No update necessary.")
		}
		return nil
	}

	// Obtain OTP.
	userOTP, err := otp.GetOTP(providedOTP, inReader)
	if err != nil {
		return fmt.Errorf("failed to obtain OTP: %w", err)
	}

	// Retrieve new session credentials using the provided MFA ARN.
	newCreds, err := aws.GetSessionToken(ctx, stsClient, mfaArn, userOTP, 3600)
	if err != nil {
		return fmt.Errorf("failed to get new session token: %w", err)
	}

	// Update the credentials file.
	if err = aws.UpdateCredentials(profile, newCreds); err != nil {
		return fmt.Errorf("failed to update credentials file: %w", err)
	}

	if verbose {
		fmt.Println("AWS credentials successfully updated.")
	}
	return nil
}

func CreateSTSClient(ctx context.Context, profile, region string) (STSCombinedClient, error) {
	opts := []func(*awsConfig.LoadOptions) error{
		awsConfig.WithSharedConfigProfile(profile),
	}
	if region != "" {
		opts = append(opts, awsConfig.WithRegion(region))
	}
	cfg, err := awsConfig.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return awsSts.NewFromConfig(cfg), nil
}

func main() {
	// Define flags.
	profileFrom := flag.String("profile-from", "default-long-term", "AWS profile to use for obtaining session credentials")
	profileTo := flag.String("profile-to", "default", "AWS profile to update with new session credentials")
	region := flag.String("region", "", "AWS region to use (auto-detected if not provided)")
	mfaArn := flag.String("mfa-arn", "", "MFA device ARN to use for authentication (if not provided, will auto lookup)")
	awsUser := flag.String("user", "", "AWS username (if not provided, defaults to current OS user)")
	otp := flag.String("otp", "", "One Time Password for authentication")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	force := flag.Bool("force", false, "Force re-authentication even if credentials are valid")
	flag.Parse()

	// Determine the AWS username if not provided.
	if *awsUser == "" {
		u, err := user.Current()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: unable to determine current OS user; please provide --user")
			os.Exit(1)
		}
		*awsUser = u.Username
	}

	ctx := context.Background()

	// Load AWS config using the source profile and region.
	opts := []func(*awsConfig.LoadOptions) error{
		awsConfig.WithSharedConfigProfile(*profileFrom),
	}
	var regionOpt string
	if *region != "" {
		regionOpt = *region
	} else if envRegion := os.Getenv("AWS_REGION"); envRegion != "" {
		regionOpt = envRegion
	} else if envDefaultRegion := os.Getenv("AWS_DEFAULT_REGION"); envDefaultRegion != "" {
		regionOpt = envDefaultRegion
	} else {
		regionOpt = "us-east-1"
	}
	opts = append(opts, awsConfig.WithRegion(regionOpt))

	cfg, err := awsConfig.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading AWS config: %v\n", err)
		os.Exit(1)
	}

	// Auto lookup MFA ARN if not provided.
	if *mfaArn == "" {
		iamClient := awsIam.NewFromConfig(cfg)
		out, err := iamClient.ListMFADevices(ctx, &awsIam.ListMFADevicesInput{
			UserName: awsPkg.String(*awsUser),
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing MFA devices for user %s: %v\n", *awsUser, err)
			os.Exit(1)
		}
		if len(out.MFADevices) == 0 {
			fmt.Fprintf(os.Stderr, "No MFA devices found for user %s\n", *awsUser)
			os.Exit(1)
		} else if len(out.MFADevices) == 1 {
			*mfaArn = *out.MFADevices[0].SerialNumber
			if *verbose {
				fmt.Fprintf(os.Stderr, "Using MFA device ARN: %s\n", *mfaArn)
			}
		} else {
			fmt.Fprintf(os.Stderr, "Multiple MFA devices found for user %s. Please specify one with --mfa-arn. Devices:\n", *awsUser)
			for _, device := range out.MFADevices {
				fmt.Fprintf(os.Stderr, "  %s\n", *device.SerialNumber)
			}
			os.Exit(1)
		}
	}

	// Clean expired tokens from the target profile.
	if err := aws.CleanExpiredTokenFromCredentials(*profileTo); err != nil {
		fmt.Fprintf(os.Stderr, "Error cleaning expired token: %v\n", err)
		os.Exit(1)
	}

	// Create the STS client using the source config.
	stsClient := awsSts.NewFromConfig(cfg)

	// Run the authentication flow.
	// Pass in the MFA ARN we determined.
	if err = RunAuthFlow(ctx, stsClient, nil, *profileTo, *otp, *force, *verbose, *mfaArn); err != nil {
		fmt.Fprintf(os.Stderr, "Authentication flow failed: %v\n", err)
		os.Exit(1)
	}
}
