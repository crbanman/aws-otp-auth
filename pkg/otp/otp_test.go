package otp

import (
	"strings"
	"testing"
)

func TestGetOTP_WithProvidedOTP(t *testing.T) {
	// When an OTP is already provided via the command-line flag,
	// GetOTP should simply return that value.
	otp, err := GetOTP("123456", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if otp != "123456" {
		t.Fatalf("Expected OTP '123456', got '%s'", otp)
	}
}

func TestGetOTP_Prompt(t *testing.T) {
	// Simulate user input via a strings.Reader.
	simulatedInput := "654321\n"
	otp, err := GetOTP("", strings.NewReader(simulatedInput))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if otp != "654321" {
		t.Fatalf("Expected OTP '654321', got '%s'", otp)
	}
}
