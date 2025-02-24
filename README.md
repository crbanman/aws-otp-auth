# AWS OTP Authentication CLI

A command-line tool for securely authenticating with AWS when Multi-Factor Authentication (MFA) or a One-Time Password (OTP) is required. This utility validates existing AWS credentials and, if needed, prompts for an OTP to generate new temporary session credentials via AWS STS. It automatically updates the AWS credentials file to streamline the authentication process.

## Features

- **Credential Validation:** Checks if existing AWS credentials are valid using `aws sts get-caller-identity`.
- **OTP Handling:** Prompts for an OTP if credentials are invalid or expired.
- **Session Token Retrieval:** Uses AWS STS `GetSessionToken` with MFA to obtain temporary credentials.
- **Automatic Credentials Update:** Backs up and updates the AWS credentials file with new session tokens.
- **Multi-Profile Support:** Works with multiple AWS profiles for different environments.
- **Error Handling:** Provides clear error messages and logs.

## Prerequisites

Ensure the following before using this tool:

- **Go 1.24.0 or later** installed.
- **AWS credentials** configured in `~/.aws/credentials`.
- **MFA device** associated with the AWS IAM user.

## Installation

### Pre-built Binaries

Pre-built binaries are available for Linux systems (amd64 and arm64) in the [Releases](https://github.com/crbanman/aws-otp-auth/releases) section.

Move the binary to a location in your PATH for easy access:

```bash
sudo mv aws-otp-auth-linux-* /usr/local/bin/aws-otp-auth
```

### macOS Users

Due to code signing requirements on macOS, pre-built binaries are not provided for macOS. macOS users should build from source following the instructions below.

### Build from Source

#### Clone the Repository

```bash
git clone https://github.com/crbanman/aws-otp-auth.git
cd aws-otp-auth
```

### Build the Application

Using `make`:

```bash
make build
```

Or manually with Go:

```bash
go build -o aws-otp-auth ./cmd/aws-otp-auth
```

## Usage

Run the CLI tool with default settings:

```bash
./aws-otp-auth
```

If additional customization is needed, use the available options:

```bash
./aws-otp-auth --profile-from default-long-term --profile-to default --mfa-arn arn:aws:iam::123456789012:mfa/your-user --otp 123456 --verbose
```

### Command-Line Flags

- `--profile-from` : Source AWS profile for obtaining session credentials (default: `default-long-term`).
- `--profile-to` : Target AWS profile for storing new session credentials (default: `default`).
- `--mfa-arn` : MFA device ARN for authentication. Auto-detects if not provided.
- `--otp` : One-Time Password for MFA authentication. Prompts interactively if omitted.
- `--verbose` : Enables detailed logging.
- `--force` : Forces re-authentication even if credentials are still valid.

### Example Usage

To update the `default` profile with new temporary credentials using MFA:

```bash
./aws-otp-auth --profile-from my-long-term-profile --profile-to default --mfa-arn arn:aws:iam::123456789012:mfa/my-mfa-device
```

If `--otp` is not supplied, the tool will prompt for it interactively.

## AWS Credentials File Format

Ensure your `~/.aws/credentials` file follows the standard INI format:

```ini
[default-long-term]
aws_access_key_id = YOUR_LONG_TERM_ACCESS_KEY
aws_secret_access_key = YOUR_LONG_TERM_SECRET_KEY

[default]
# Updated dynamically by the CLI
aws_access_key_id = NEW_TEMPORARY_ACCESS_KEY
aws_secret_access_key = NEW_TEMPORARY_SECRET_KEY
aws_session_token = NEW_TEMPORARY_SESSION_TOKEN
aws_session_token_expiration = 2025-02-24T15:04:05Z
```

## Development & Testing

### Running Tests

Execute tests using:

```bash
make test
```

Or manually:

```bash
go test -v ./...
```

## Troubleshooting

- **Invalid Credentials:** Ensure `~/.aws/credentials` contains valid long-term access keys.
- **MFA Device Not Found:** Verify the MFA device ARN is correct or allow the tool to auto-detect.
- **Session Token Not Updated:** Check if an expired session token is still in use and use `--force` to override.
