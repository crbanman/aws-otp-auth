package otp

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// GetOTP returns the provided OTP if non-empty.
// Otherwise, it prompts the user with "Enter OTP:" and reads input from the given reader.
func GetOTP(providedOTP string, inReader io.Reader) (string, error) {
	if providedOTP != "" {
		return providedOTP, nil
	}
	if inReader == nil {
		inReader = os.Stdin
	}
	fmt.Print("Enter OTP: ")
	reader := bufio.NewReader(inReader)
	otp, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(otp), nil
}
