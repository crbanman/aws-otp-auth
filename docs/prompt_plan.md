Below is a detailed blueprint for building the AWS OTP Authentication CLI project. Following the blueprint is an iterative breakdown into small, test-driven prompts for a code-generation LLM. Each prompt builds on previous work and ends with wiring the functionality together.

Detailed Blueprint

Project Overview:
A CLI tool to help developers authenticate with AWS when MFA/OTP is required. The tool will:
• Check if current AWS credentials (read from ~/.aws/credentials) are valid by calling aws sts get-caller-identity.
• If valid, exit; if not, prompt the user (or accept a command-line flag) for an OTP.
• Use AWS STS (GetSessionToken) with the provided OTP (and MFA ARN) to obtain temporary credentials.
• Backup and update the credentials file with the new session token.
• Handle errors clearly and support multiple AWS profiles.

Key Components: 1. Project Setup:
• Create a new Go project.
• Initialize a module (using go mod init).
• Set up the main file (main.go) and create additional packages if needed. 2. CLI Argument Parsing:
• Use a simple CLI parsing method (e.g., Go’s built-in flag package) to capture:
• --profile (AWS profile to use; default to “default”)
• --otp (OTP code provided as an argument)
• --verbose (enable debug output)
• --force (bypass authentication check) 3. Credentials File Handling:
• Read and parse ~/.aws/credentials (an INI format file).
• Validate that a session token is present and check its expiry.
• Create unit tests to validate the parsing logic. 4. Authentication Check:
• Implement a function that calls aws sts get-caller-identity using the AWS SDK for Go v2.
• If the call is successful, the existing credentials are valid.
• Write tests (or mocks) to simulate both valid and invalid scenarios. 5. OTP Input Handling:
• If the OTP is not provided via --otp, prompt the user interactively.
• Ensure the OTP is captured correctly.
• Create tests for OTP input (using dependency injection or similar strategies to simulate user input). 6. AWS Session Token Authentication:
• Implement a function that uses the AWS SDK for Go v2 to call sts.GetSessionToken().
• The call should use the MFA ARN and OTP provided by the user.
• Handle error responses (such as incorrect OTP or network issues) gracefully.
• Write unit tests to simulate AWS STS responses. 7. Credentials Update:
• Backup the existing ~/.aws/credentials file to ~/.aws/credentials.bak.
• Update the specified AWS profile with the new session token.
• Ensure that the update is atomic and test for proper file writing. 8. Error Handling & Logging:
• At each stage, include clear error messages and instructions for the user.
• If a step fails (e.g., file missing, invalid OTP), the program should exit cleanly.
• Use logging (or verbose output) to assist in debugging when --verbose is used. 9. Integration & Final Wiring:
• In the main function, wire together:
• Check if credentials are valid.
• If not, prompt for or use the provided OTP.
• Authenticate with AWS STS.
• Update the credentials file.
• Confirm the update by re-checking authentication.
• Ensure that each component is connected and that there is no orphaned code. 10. Testing Plan:
• Unit Tests: For credential parsing, OTP handling, and AWS STS call wrappers.
• Integration Tests: Simulate the full authentication flow with valid and invalid credentials.
• Edge Cases: Handle missing/corrupted credentials file, incorrect MFA ARN, and user interrupts (Ctrl+C).

Iterative, Test-Driven Prompts for Code Generation

Below is a series of prompts designed to guide the implementation in small, integrated steps. Each prompt is tagged as text using code blocks in Markdown.

Prompt 1: Project Setup and CLI Skeleton

Prompt 1:
Create a new Go project for the AWS OTP Authentication CLI. Initialize a Go module and create a file called main.go. In main.go, implement a basic CLI skeleton using the built-in flag package that accepts the following command-line arguments:

- --profile (string, default "default")
- --otp (string, optional)
- --verbose (boolean)
- --force (boolean)

The program should print a help message when run with the --help flag. At this stage, simply print the parsed flag values. Include a minimal unit test that ensures the flags are correctly parsed.

Prompt 2: Implement Credentials File Handling

Prompt 2:
Implement a function (e.g., ReadAWSCredentials) that reads and parses the AWS credentials file from ~/.aws/credentials. Assume the file is in INI format and extract the credentials for a given profile. Create a corresponding unit test to verify that your function correctly reads and parses a sample credentials file. Ensure your function handles cases where the file is missing or malformed.

Prompt 3: Build AWS Authentication Check

Prompt 3:
Develop a function (e.g., CheckAuthentication) that uses the AWS SDK for Go v2 to call `aws sts get-caller-identity`. This function should verify whether the current credentials (from the specified profile) are valid. Simulate both a valid and an invalid response, and include unit tests (or use dependency injection/mocks) to confirm that your function handles both cases correctly.

Prompt 4: Add OTP Input Handling

Prompt 4:
Create a function (e.g., GetOTP) to handle OTP input. This function should:

- Accept an OTP as a command-line argument if provided.
- Otherwise, prompt the user with "Enter OTP:" and read the input from stdin.
  Include tests (e.g., by simulating user input) to ensure that the OTP is captured accurately.

Prompt 5: Implement AWS Session Token Retrieval

Prompt 5:
Implement a function (e.g., GetSessionToken) that calls AWS STS's GetSessionToken API using the aws-sdk-go-v2. The function should take the MFA ARN, OTP code, and any other necessary parameters, then return the new session credentials (access key, secret key, session token, and expiration). Add unit tests (or mocks) to simulate successful and failed responses from AWS.

Prompt 6: Update Credentials File with New Session Token

Prompt 6:
Write a function (e.g., UpdateCredentials) that performs the following:

- Backs up the existing ~/.aws/credentials file to ~/.aws/credentials.bak.
- Overwrites the specified profile in the credentials file with the new session credentials obtained from the previous step.
  Include tests to verify that the backup is created and that the credentials file is updated correctly.

Prompt 7: Integrate and Wire Everything Together

Prompt 7:
Wire together all the functions developed in previous prompts in the main function of your CLI application. The workflow should be:

1. Parse CLI flags.
2. Read and validate the current credentials using ReadAWSCredentials and CheckAuthentication.
3. If credentials are invalid or --force is specified, obtain an OTP using GetOTP.
4. Retrieve new session credentials using GetSessionToken.
5. Update the credentials file using UpdateCredentials.
6. Print a success message.
   Add error handling at each step so that any failure provides clear output to the user. Also, include an integration test (or a simulated full run) that demonstrates the complete authentication flow.

Each prompt is designed to be incremental. Begin with the project skeleton and CLI parsing, then add file handling, AWS checks, OTP input, session token retrieval, and finally integrate everything with proper error handling and tests. This approach ensures that every piece is well tested before moving on to the next, with no orphaned code.
