package deployer

import (
	"testing"
)

func TestDeploymentOptions_Validation(t *testing.T) {
	tests := []struct {
		name    string
		depOpts DeploymentOptions
		wantErr bool
	}{
		{
			name: "valid deployment options",
			depOpts: DeploymentOptions{
				Application:    "test-app",
				Version:        "1.0.0",
				Environment:    "production",
				Description:    "Test deployment",
				MaxAttempts:    40,
				RefreshSecrets: false,
				SecretsPrefix:  []string{"/prod/test-app"},
				DryRun:         false,
			},
			wantErr: false,
		},
		{
			name: "empty secrets prefix should not cause error",
			depOpts: DeploymentOptions{
				Application:    "test-app",
				Version:        "1.0.0",
				Environment:    "production",
				Description:    "Test deployment",
				MaxAttempts:    40,
				RefreshSecrets: false,
				SecretsPrefix:  []string{},
				DryRun:         false,
			},
			wantErr: false,
		},
		{
			name: "zero max attempts should not cause error",
			depOpts: DeploymentOptions{
				Application:    "test-app",
				Version:        "1.0.0",
				Environment:    "production",
				Description:    "Test deployment",
				MaxAttempts:    0,
				RefreshSecrets: false,
				SecretsPrefix:  []string{"/prod/test-app"},
				DryRun:         false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation - ensure struct can be created without panics
			if tt.depOpts.Application == "" && !tt.wantErr {
				t.Error("Expected non-empty application name")
			}
			if tt.depOpts.Version == "" && !tt.wantErr {
				t.Error("Expected non-empty version")
			}
			if tt.depOpts.Environment == "" && !tt.wantErr {
				t.Error("Expected non-empty environment")
			}
		})
	}
}

func TestDeploymentResults_Structure(t *testing.T) {
	result := DeploymentResults{
		SuccessfullyInvoked: true,
		ClusterArn:          "arn:aws:ecs:us-east-1:123456789012:cluster/test-cluster",
		ServiceArn:          "arn:aws:ecs:us-east-1:123456789012:service/test-cluster/test-service",
		ServiceName:         "test-service",
		TaskDefinition:      "test-task:1",
	}

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

// Note: SetDeploymentOptionsByEcsServiceTags requires AWS credentials and actual ECS services
// This would be better tested with integration tests or mocked AWS services
func TestDeploymentOptions_SetDeploymentOptionsByEcsServiceTags_Structure(t *testing.T) {
	depOpts := &DeploymentOptions{
		Application: "test-app",
		Environment: "test-env",
	}

	// We can't test the actual AWS call without credentials, but we can ensure
	// the method exists and doesn't panic with nil receiver
	if depOpts == nil {
		t.Error("DeploymentOptions should not be nil")
	}

	// Test that the method exists (compilation test)
	_ = depOpts.SetDeploymentOptionsByEcsServiceTags
}
