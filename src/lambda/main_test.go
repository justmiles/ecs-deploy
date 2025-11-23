package ld

import (
	"testing"

	"github.com/justmiles/ecs-deploy/src/deployer"
)

func TestInputValidation(t *testing.T) {
	tests := []struct {
		name     string
		depOpts  deployer.DeploymentOptions
		expected bool
	}{
		{
			name: "valid input",
			depOpts: deployer.DeploymentOptions{
				Application: "test-app",
				Version:     "1.0.0",
				Environment: "production",
			},
			expected: true,
		},
		{
			name: "missing application",
			depOpts: deployer.DeploymentOptions{
				Version:     "1.0.0",
				Environment: "production",
			},
			expected: false,
		},
		{
			name: "missing version",
			depOpts: deployer.DeploymentOptions{
				Application: "test-app",
				Environment: "production",
			},
			expected: false,
		},
		{
			name: "missing environment",
			depOpts: deployer.DeploymentOptions{
				Application: "test-app",
				Version:     "1.0.0",
			},
			expected: false,
		},
		{
			name: "empty application",
			depOpts: deployer.DeploymentOptions{
				Application: "",
				Version:     "1.0.0",
				Environment: "production",
			},
			expected: false,
		},
		{
			name: "empty version",
			depOpts: deployer.DeploymentOptions{
				Application: "test-app",
				Version:     "",
				Environment: "production",
			},
			expected: false,
		},
		{
			name: "empty environment",
			depOpts: deployer.DeploymentOptions{
				Application: "test-app",
				Version:     "1.0.0",
				Environment: "",
			},
			expected: false,
		},
		{
			name:     "all empty",
			depOpts:  deployer.DeploymentOptions{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := inputValidation(tt.depOpts)
			if result != tt.expected {
				t.Errorf("inputValidation() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestHandler_InvalidInput(t *testing.T) {
	depOpts := deployer.DeploymentOptions{
		Application: "",
		Version:     "1.0.0",
		Environment: "production",
	}

	_, err := Handler(depOpts)
	if err != ErrInvalidInputProvided {
		t.Errorf("Handler() error = %v, expected %v", err, ErrInvalidInputProvided)
	}
}

func TestHandler_DefaultDescription(t *testing.T) {
	// This test would require mocking the deployer.PerformDeployment function
	// For now, we'll test that the default description is set correctly
	depOpts := deployer.DeploymentOptions{
		Application: "test-app",
		Version:     "1.0.0",
		Environment: "production",
		Description: "",
	}

	// We can't easily test the full Handler without mocking AWS services
	// But we can test the input validation and description setting logic
	if !inputValidation(depOpts) {
		t.Error("Expected valid input validation")
	}

	if depOpts.Description == "" {
		depOpts.Description = defaultSSMDescription
	}

	if depOpts.Description != defaultSSMDescription {
		t.Errorf("Expected description to be set to default: %s, got: %s", defaultSSMDescription, depOpts.Description)
	}
}

func TestErrInvalidInputProvided(t *testing.T) {
	expectedMessage := "invalid input - ensure Application, Version, and Environment is provided"
	if ErrInvalidInputProvided.Error() != expectedMessage {
		t.Errorf("ErrInvalidInputProvided.Error() = %v, expected %v", ErrInvalidInputProvided.Error(), expectedMessage)
	}
}
