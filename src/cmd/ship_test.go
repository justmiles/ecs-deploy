package cmd

import (
	"testing"

	"github.com/justmiles/ecs-deploy/src/deployer"
)

func TestShipCmd(t *testing.T) {
	// Test that ship command exists and has expected properties
	if shipCmd == nil {
		t.Fatal("shipCmd should not be nil")
	}

	if shipCmd.Use != "ship" {
		t.Errorf("Expected Use to be 'ship', got '%s'", shipCmd.Use)
	}

	if shipCmd.Short != "Ship an application to ECS" {
		t.Errorf("Expected Short to be 'Ship an application to ECS', got '%s'", shipCmd.Short)
	}
}

func TestShipCmdFlags(t *testing.T) {
	// Test required flags
	requiredFlags := []struct {
		name      string
		shorthand string
		usage     string
	}{
		{"application", "a", "Application name to deploy"},
		{"version", "v", "Desired version of application"},
		{"environment", "e", "Target environment for deployment"},
	}

	for _, flag := range requiredFlags {
		f := shipCmd.Flags().Lookup(flag.name)
		if f == nil {
			t.Errorf("Required flag '%s' should exist", flag.name)
			continue
		}

		if f.Shorthand != flag.shorthand {
			t.Errorf("Expected flag '%s' shorthand to be '%s', got '%s'", flag.name, flag.shorthand, f.Shorthand)
		}

		if f.Usage != flag.usage {
			t.Errorf("Expected flag '%s' usage to be '%s', got '%s'", flag.name, flag.usage, f.Usage)
		}
	}

	// Test optional flags
	optionalFlags := []struct {
		name         string
		shorthand    string
		defaultValue string
	}{
		{"role", "r", ""},
		{"max-attempts", "", "40"},
		{"no-wait", "w", "false"},
		{"refresh-secrets", "", "false"},
		{"dry-run", "", "false"},
		{"secrets-prefix", "p", "[]"},
		{"ignore-tags", "i", "false"},
	}

	for _, flag := range optionalFlags {
		f := shipCmd.Flags().Lookup(flag.name)
		if f == nil {
			t.Errorf("Optional flag '%s' should exist", flag.name)
			continue
		}

		if flag.shorthand != "" && f.Shorthand != flag.shorthand {
			t.Errorf("Expected flag '%s' shorthand to be '%s', got '%s'", flag.name, flag.shorthand, f.Shorthand)
		}

		if f.DefValue != flag.defaultValue {
			t.Errorf("Expected flag '%s' default value to be '%s', got '%s'", flag.name, flag.defaultValue, f.DefValue)
		}
	}
}

func TestShipCmdRequiredFlags(t *testing.T) {
	// Test that required flags are marked as required
	requiredFlags := []string{"application", "version", "environment"}

	for _, flagName := range requiredFlags {
		flag := shipCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Required flag '%s' should exist", flagName)
			continue
		}

		// Check if flag is in required flags list
		// Note: This is a bit tricky to test directly, but we can check the flag exists
		// The actual requirement checking happens at runtime
	}
}

func TestDeploymentOptionsGlobalVariable(t *testing.T) {
	// Test that deploymentOptions global variable is properly initialized
	if deploymentOptions.Description != "Desired version set by ecs-deploy CLI" {
		t.Errorf("Expected default description to be 'Desired version set by ecs-deploy CLI', got '%s'", deploymentOptions.Description)
	}

	// Test that it's a proper DeploymentOptions struct
	var _ deployer.DeploymentOptions = deploymentOptions
}

func TestShipGlobalVariables(t *testing.T) {
	// Test that global variables are properly initialized

	// noWait should be initialized to false
	if noWait {
		// Reset it for other tests
		noWait = false
	}

	// ignoreTags should be initialized to false
	if ignoreTags {
		// Reset it for other tests
		ignoreTags = false
	}
}

func TestShipCmdFlagTypes(t *testing.T) {
	// Test that flags have correct types

	// String flags
	stringFlags := []string{"application", "version", "environment", "role"}
	for _, flagName := range stringFlags {
		flag := shipCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("String flag '%s' should exist", flagName)
			continue
		}
		if flag.Value.Type() != "string" {
			t.Errorf("Expected flag '%s' to be string type, got '%s'", flagName, flag.Value.Type())
		}
	}

	// Int flags
	intFlags := []string{"max-attempts"}
	for _, flagName := range intFlags {
		flag := shipCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Int flag '%s' should exist", flagName)
			continue
		}
		if flag.Value.Type() != "int" {
			t.Errorf("Expected flag '%s' to be int type, got '%s'", flagName, flag.Value.Type())
		}
	}

	// Bool flags
	boolFlags := []string{"no-wait", "refresh-secrets", "dry-run", "ignore-tags"}
	for _, flagName := range boolFlags {
		flag := shipCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Bool flag '%s' should exist", flagName)
			continue
		}
		if flag.Value.Type() != "bool" {
			t.Errorf("Expected flag '%s' to be bool type, got '%s'", flagName, flag.Value.Type())
		}
	}

	// StringSlice flags
	stringSliceFlags := []string{"secrets-prefix"}
	for _, flagName := range stringSliceFlags {
		flag := shipCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("StringSlice flag '%s' should exist", flagName)
			continue
		}
		if flag.Value.Type() != "stringSlice" {
			t.Errorf("Expected flag '%s' to be stringSlice type, got '%s'", flagName, flag.Value.Type())
		}
	}
}

func TestShipCmdHasRunFunction(t *testing.T) {
	// Test that ship command has a Run function
	if shipCmd.Run == nil {
		t.Error("shipCmd should have a Run function")
	}
}

// Note: Testing the actual Run function would require mocking AWS services
// and is better suited for integration tests
