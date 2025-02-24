# AWS OTP Authentication CLI - TODO Checklist

This checklist outlines all the steps required to build the AWS OTP Authentication CLI, from project setup to final integration and testing.

---

## Project Setup

- [x] **Initialize Project:**
  - [x] Create a new Go module (e.g., `go mod init aws-otp-auth`)
  - [x] Create a basic project structure with a `main.go` file

---

## CLI Argument Parsing

- [x] **Implement CLI Flags in `main.go`:**
  - [x] Parse `--profile` (string, default "default")
  - [x] Parse `--otp` (string, optional)
  - [x] Parse `--verbose` (boolean)
  - [x] Parse `--force` (boolean)
- [x] **Help Message:**
  - [x] Display usage information when `--help` is provided
- [x] **Testing:**
  - [x] Write unit tests to verify that all flags are parsed correctly

---

## Credentials File Handling

- [x] **Read & Parse Credentials:**
  - [x] Implement `ReadAWSCredentials` to read `~/.aws/credentials`
  - [x] Parse the INI formatted file and extract the credentials for the specified profile
- [x] **Error Handling:**
  - [x] Handle missing or malformed credentials file
- [x] **Testing:**
  - [x] Write unit tests using a sample credentials file to validate parsing

---

## AWS Authentication Check

- [x] **Implement Authentication Check:**
  - [x] Create a function `CheckAuthentication` that calls `aws sts get-caller-identity` using the AWS SDK for Go v2
  - [x] Validate that the credentials are active and correct
- [x] **Testing:**
  - [x] Simulate both valid and invalid responses using mocks or dependency injection

---

## OTP Input Handling

- [x] **Implement OTP Capture:**
  - [x] Create a function `GetOTP` that:
    - [x] Accepts OTP from the `--otp` flag if provided
    - [x] Otherwise, prompts the user with "Enter OTP:" and reads input from stdin
- [x] **Testing:**
  - [x] Write tests that simulate user input and confirm OTP is captured correctly

---

## AWS Session Token Retrieval

- [x] **Implement Session Token Request:**
  - [x] Create a function `GetSessionToken` that calls AWS STS's `GetSessionToken` API with:
    - [x] MFA ARN
    - [x] OTP code
  - [x] Ensure the function returns the new session credentials (access key, secret key, session token, expiration)
- [x] **Error Handling:**
  - [x] Handle scenarios for incorrect OTP or network issues gracefully
- [x] **Testing:**
  - [x] Write unit tests or use mocks to simulate both successful and failed responses

---

## Credentials File Update

- [x] **Backup and Update:**
  - [x] Implement `UpdateCredentials` to:
    - [x] Backup `~/.aws/credentials` to `~/.aws/credentials.bak`
    - [x] Update the specified profile with the new session credentials
  - [x] Ensure that file updates are atomic and consistent
- [x] **Testing:**
  - [x] Write tests to verify the backup creation and the correct update of the credentials file

---

## Integration & Final Wiring

- [x] **Integrate Components in `main.go`:**
  - [x] Parse CLI flags
  - [x] Read credentials using `ReadAWSCredentials`
  - [x] Check authentication with `CheckAuthentication`
  - [x] If credentials are invalid or `--force` is specified, then:
    - [x] Obtain OTP using `GetOTP`
    - [x] Retrieve new credentials using `GetSessionToken`
    - [x] Update the credentials file using `UpdateCredentials`
  - [x] Print a success message after updating
- [x] **Error Handling:**
  - [x] Implement clear error messages for each failure scenario
  - [x] Handle user interrupts (e.g., Ctrl+C) gracefully
- [x] **Testing:**
  - [x] Write an integration test or simulate a full run of the authentication flow

---

## Testing & Documentation

- [x] **Unit Tests:**
  - [x] Write tests for each individual function (CLI parsing, file handling, OTP input, AWS calls)
- [x] **Integration Tests:**
  - [x] Write tests for the complete authentication flow
- [x] **Edge Case Testing:**
  - [x] Validate behavior with missing credentials file
  - [x] Test scenarios with invalid MFA ARN or expired tokens
- [x] **Documentation:**
  - [x] Write clear usage instructions and configuration guidelines
  - [x] Ensure verbose mode outputs helpful debugging information

---

## Final Review & Cleanup

- [ ] **Code Review:**
  - [ ] Verify adherence to best practices and coding standards
- [ ] **Remove Debug Code:**
  - [ ] Remove or disable debugging outputs for production
- [ ] **Final Testing:**
  - [ ] Run a complete integration test to ensure all components work together without orphaned code
- [ ] **Project Documentation:**
  - [ ] Update the project README with usage examples and setup instructions
