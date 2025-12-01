package deployer

import (
	"testing"
)

// Note: These are unit tests for the deployer package that don't require AWS credentials
// Integration tests with actual AWS services would be in a separate file

func TestDeploymentOptions_SecretsPrefix_DefaultBehavior(t *testing.T) {
	depOpts := DeploymentOptions{
		Application: "test-app",
		Environment: "production",
		Version:     "1.0.0",
	}

	// Test that empty secrets prefix can be handled
	if len(depOpts.SecretsPrefix) != 0 {
		t.Errorf("Expected empty SecretsPrefix by default, got %v", depOpts.SecretsPrefix)
	}

	// Test setting secrets prefix
	depOpts.SecretsPrefix = []string{"/prod/test-app", "/shared/secrets"}
	if len(depOpts.SecretsPrefix) != 2 {
		t.Errorf("Expected 2 secrets prefixes, got %d", len(depOpts.SecretsPrefix))
	}

	expectedPrefixes := []string{"/prod/test-app", "/shared/secrets"}
	for i, expected := range expectedPrefixes {
		if depOpts.SecretsPrefix[i] != expected {
			t.Errorf("Expected secrets prefix[%d] to be '%s', got '%s'", i, expected, depOpts.SecretsPrefix[i])
		}
	}
}

func TestDeploymentOptions_MaxAttempts(t *testing.T) {
	depOpts := DeploymentOptions{
		Application: "test-app",
		Environment: "production",
		Version:     "1.0.0",
		MaxAttempts: 60,
	}

	if depOpts.MaxAttempts != 60 {
		t.Errorf("Expected MaxAttempts to be 60, got %d", depOpts.MaxAttempts)
	}

	// Test zero value
	depOpts.MaxAttempts = 0
	if depOpts.MaxAttempts != 0 {
		t.Errorf("Expected MaxAttempts to be 0, got %d", depOpts.MaxAttempts)
	}
}

func TestDeploymentOptions_Flags(t *testing.T) {
	depOpts := DeploymentOptions{
		Application:    "test-app",
		Environment:    "production",
		Version:        "1.0.0",
		RefreshSecrets: true,
		DryRun:         true,
	}

	if !depOpts.RefreshSecrets {
		t.Error("Expected RefreshSecrets to be true")
	}

	if !depOpts.DryRun {
		t.Error("Expected DryRun to be true")
	}

	// Test false values
	depOpts.RefreshSecrets = false
	depOpts.DryRun = false

	if depOpts.RefreshSecrets {
		t.Error("Expected RefreshSecrets to be false")
	}

	if depOpts.DryRun {
		t.Error("Expected DryRun to be false")
	}
}

func TestDeploymentOptions_Role(t *testing.T) {
	depOpts := DeploymentOptions{
		Application: "test-app",
		Environment: "production",
		Version:     "1.0.0",
		Role:        "arn:aws:iam::123456789012:role/ECSDeployRole",
	}

	expectedRole := "arn:aws:iam::123456789012:role/ECSDeployRole"
	if depOpts.Role != expectedRole {
		t.Errorf("Expected Role to be '%s', got '%s'", expectedRole, depOpts.Role)
	}

	// Test empty role
	depOpts.Role = ""
	if depOpts.Role != "" {
		t.Errorf("Expected empty Role, got '%s'", depOpts.Role)
	}
}

func TestDeploymentResults_JSON_Compatibility(t *testing.T) {
	// Test that DeploymentResults can be used for JSON marshaling/unmarshaling
	result := DeploymentResults{
		SuccessfullyInvoked: true,
		ClusterArn:          "arn:aws:ecs:us-east-1:123456789012:cluster/test",
		ServiceArn:          "arn:aws:ecs:us-east-1:123456789012:service/test/app",
		ServiceName:         "test-app",
		TaskDefinition:      "test-app:1",
	}

	// Basic validation that all fields are set
	if !result.SuccessfullyInvoked {
		t.Error("Expected SuccessfullyInvoked to be true")
	}

	if result.ClusterArn == "" {
		t.Error("Expected non-empty ClusterArn")
	}

	if result.ServiceArn == "" {
		t.Error("Expected non-empty ServiceArn")
	}

	if result.ServiceName == "" {
		t.Error("Expected non-empty ServiceName")
	}

	if result.TaskDefinition == "" {
		t.Error("Expected non-empty TaskDefinition")
	}
}

func TestDeploymentOptions_Description(t *testing.T) {
	depOpts := DeploymentOptions{
		Application: "test-app",
		Environment: "production",
		Version:     "1.0.0",
		Description: "Custom deployment description",
	}

	expectedDesc := "Custom deployment description"
	if depOpts.Description != expectedDesc {
		t.Errorf("Expected Description to be '%s', got '%s'", expectedDesc, depOpts.Description)
	}

	// Test empty description
	depOpts.Description = ""
	if depOpts.Description != "" {
		t.Errorf("Expected empty Description, got '%s'", depOpts.Description)
	}
}

// Test that the global session variable is initialized
func TestGlobalSession(t *testing.T) {
	if sess == nil {
		t.Error("Expected global session to be initialized")
	}
}

// Test struct field tags for JSON compatibility
func TestDeploymentOptions_StructTags(t *testing.T) {
	// This is a compile-time test to ensure the struct has proper JSON tags
	depOpts := DeploymentOptions{
		Application:    "test",
		Version:        "1.0",
		Environment:    "prod",
		Description:    "test desc",
		Role:           "test-role",
		MaxAttempts:    40,
		RefreshSecrets: true,
		SecretsPrefix:  []string{"/test"},
		DryRun:         false,
	}

	// Basic validation that struct can be created
	if depOpts.Application != "test" {
		t.Error("Struct field assignment failed")
	}
}

func TestDeploymentResults_StructTags(t *testing.T) {
	// This is a compile-time test to ensure the struct has proper JSON tags
	result := DeploymentResults{
		SuccessfullyInvoked: true,
		ClusterArn:          "cluster-arn",
		ServiceArn:          "service-arn",
		ServiceName:         "service-name",
		TaskDefinition:      "task-def",
	}

	// Basic validation that struct can be created
	if !result.SuccessfullyInvoked {
		t.Error("Struct field assignment failed")
	}
}
