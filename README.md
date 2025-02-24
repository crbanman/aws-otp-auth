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

// ... existing code ...
