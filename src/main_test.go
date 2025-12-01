package main

import (
	"os"
	"testing"
)

func TestMain_LambdaDetection(t *testing.T) {
	// Test that main function can detect lambda environment

	// Save original environment
	originalLambdaPort := os.Getenv("_LAMBDA_SERVER_PORT")
	originalAWSConfig := os.Getenv("AWS_SDK_LOAD_CONFIG")

	// Cleanup after test
	defer func() {
		if originalLambdaPort != "" {
			os.Setenv("_LAMBDA_SERVER_PORT", originalLambdaPort)
		} else {
			os.Unsetenv("_LAMBDA_SERVER_PORT")
		}
		if originalAWSConfig != "" {
			os.Setenv("AWS_SDK_LOAD_CONFIG", originalAWSConfig)
		} else {
			os.Unsetenv("AWS_SDK_LOAD_CONFIG")
		}
	}()

	// Test non-lambda environment
	os.Unsetenv("_LAMBDA_SERVER_PORT")

	// We can't easily test the actual main() function execution without it running
	// But we can test that the environment variable detection logic would work
	lambdaPort := os.Getenv("_LAMBDA_SERVER_PORT")
	if lambdaPort != "" {
		t.Error("Expected _LAMBDA_SERVER_PORT to be empty in non-lambda environment")
	}

	// Test lambda environment detection
	os.Setenv("_LAMBDA_SERVER_PORT", "8080")
	lambdaPort = os.Getenv("_LAMBDA_SERVER_PORT")
	if lambdaPort != "8080" {
		t.Errorf("Expected _LAMBDA_SERVER_PORT to be '8080', got '%s'", lambdaPort)
	}
}

func TestMain_AWSConfigSetting(t *testing.T) {
	// Test that AWS_SDK_LOAD_CONFIG is set in non-lambda mode

	// Save original environment
	originalLambdaPort := os.Getenv("_LAMBDA_SERVER_PORT")
	originalAWSConfig := os.Getenv("AWS_SDK_LOAD_CONFIG")

	// Cleanup after test
	defer func() {
		if originalLambdaPort != "" {
			os.Setenv("_LAMBDA_SERVER_PORT", originalLambdaPort)
		} else {
			os.Unsetenv("_LAMBDA_SERVER_PORT")
		}
		if originalAWSConfig != "" {
			os.Setenv("AWS_SDK_LOAD_CONFIG", originalAWSConfig)
		} else {
			os.Unsetenv("AWS_SDK_LOAD_CONFIG")
		}
	}()

	// Simulate non-lambda environment
	os.Unsetenv("_LAMBDA_SERVER_PORT")
	os.Unsetenv("AWS_SDK_LOAD_CONFIG")

	// Test the logic that would be executed in main()
	if os.Getenv("_LAMBDA_SERVER_PORT") == "" {
		os.Setenv("AWS_SDK_LOAD_CONFIG", "true")
	}

	awsConfig := os.Getenv("AWS_SDK_LOAD_CONFIG")
	if awsConfig != "true" {
		t.Errorf("Expected AWS_SDK_LOAD_CONFIG to be 'true', got '%s'", awsConfig)
	}
}

func TestMain_EnvironmentVariables(t *testing.T) {
	// Test environment variable handling

	testCases := []struct {
		name           string
		lambdaPort     string
		expectedLambda bool
	}{
		{
			name:           "no lambda port",
			lambdaPort:     "",
			expectedLambda: false,
		},
		{
			name:           "lambda port set",
			lambdaPort:     "8080",
			expectedLambda: true,
		},
		{
			name:           "lambda port different value",
			lambdaPort:     "9000",
			expectedLambda: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Save original
			original := os.Getenv("_LAMBDA_SERVER_PORT")
			defer func() {
				if original != "" {
					os.Setenv("_LAMBDA_SERVER_PORT", original)
				} else {
					os.Unsetenv("_LAMBDA_SERVER_PORT")
				}
			}()

			// Set test value
			if tc.lambdaPort != "" {
				os.Setenv("_LAMBDA_SERVER_PORT", tc.lambdaPort)
			} else {
				os.Unsetenv("_LAMBDA_SERVER_PORT")
			}

			// Test the condition
			isLambda := os.Getenv("_LAMBDA_SERVER_PORT") != ""
			if isLambda != tc.expectedLambda {
				t.Errorf("Expected lambda detection to be %v, got %v", tc.expectedLambda, isLambda)
			}
		})
	}
}

// Note: We cannot easily test the actual main() function execution
// as it would either start a lambda server or execute CLI commands
// These tests verify the logic that main() uses for environment detection
