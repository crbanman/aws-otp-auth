# AWS OTP Authentication CLI - Developer Specification

## Overview

This CLI application is designed to help individual developers authenticate with AWS CLI when an OTP is required. It automates checking authentication status, prompting the user for an OTP when necessary, and updating AWS credentials accordingly.

## Functional Requirements

### User Flow

1. Check if the user is already authenticated by:
   - Inspecting `~/.aws/credentials` for an active session.
   - Running `aws sts get-caller-identity` to verify session validity.
2. If authentication is valid, exit successfully.
3. If authentication is required:
   - Prompt the user for an OTP.
   - Authenticate using `aws sts get-session-token --serial-number <MFA_ARN> --token-code <OTP>`.
   - Update `~/.aws/credentials` with the new session token.
   - Print a success message.
4. Handle errors appropriately and guide the user when necessary.

### Command Syntax

- **Basic Usage:** `aws-otp-auth --profile myprofile`
- **Optional Flags:**
  - `--otp <code>`: Provide OTP via command-line argument.
  - `--verbose`: Show additional debugging output.
  - `--force`: Bypass authentication check and request OTP.

### Profile Handling

- If the user has multiple AWS profiles, they must specify one with `--profile`.
- Defaults to the `default` AWS profile if none is provided.

## Technical Implementation

### Authentication Check

- Read `~/.aws/credentials` to check for valid session tokens.
- Run `aws sts get-caller-identity` to confirm authentication.

### OTP Input

- Prompt the user interactively: `Enter OTP:`
- Allow OTP input via `--otp <code>` for automation.

### AWS Authentication

- Use the AWS SDK for Go (`aws-sdk-go-v2`) to:
  - Call `sts.GetSessionToken()` with MFA.
  - Store temporary credentials in `~/.aws/credentials` under the existing profile.
- Backup the original credentials before modifying.

### Credential Handling

- Overwrite existing profile credentials in `~/.aws/credentials`.
- Backup original credentials to `~/.aws/credentials.bak`.
- No file locking mechanism (assumes single-user CLI execution).
- Suppress output unless an error occurs.
- Display success message after updating credentials.

## Error Handling

### Authentication Failures

- If `~/.aws/credentials` is missing or unreadable, print an error and exit.
- If OTP authentication fails (e.g., incorrect OTP, network issues), prompt the user to retry manually.
- If credentials update fails, print new credentials so the user can manually update their config.
- If an issue is obvious, report it clearly to the user; otherwise, assume expiration and request re-authentication.

### User Interrupts

- If the user presses `Ctrl+C`, exit immediately without modifying credentials.

## Testing Plan

### Unit Tests

- Validate authentication check logic.
- Ensure OTP input is correctly parsed.
- Test credential file parsing and updates.

### Integration Tests

- Verify authentication against AWS with valid/invalid OTPs.
- Ensure credentials update correctly for various AWS profiles.
- Check behavior when credentials file is missing or corrupted.

### Edge Cases

- Handling of expired/invalid session tokens.
- Handling incorrect MFA ARN configurations.
- Running without required AWS CLI dependencies.

---

This specification provides all necessary details for a developer to begin implementation immediately.
