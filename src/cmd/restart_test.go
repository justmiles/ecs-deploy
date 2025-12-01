package cmd

import (
	"testing"
)

func TestRestartCmd(t *testing.T) {
	// Test that restart command exists and has expected properties
	if restartCmd == nil {
		t.Fatal("restartCmd should not be nil")
	}

	if restartCmd.Use != "restart" {
		t.Errorf("Expected Use to be 'restart', got '%s'", restartCmd.Use)
	}

	if restartCmd.Short != "gracefully restart/redeploy an application" {
		t.Errorf("Expected Short to be 'gracefully restart/redeploy an application', got '%s'", restartCmd.Short)
	}
}

func TestRestartCmdFlags(t *testing.T) {
	// Test required flags
	requiredFlags := []struct {
		name      string
		shorthand string
		usage     string
	}{
		{"application", "a", "Application name to redeploy"},
		{"environment", "e", "Target environment"},
	}

	for _, flag := range requiredFlags {
		f := restartCmd.Flags().Lookup(flag.name)
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
	}

	for _, flag := range optionalFlags {
		f := restartCmd.Flags().Lookup(flag.name)
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

func TestRestartCmdHasRunFunction(t *testing.T) {
	// Test that restart command has a Run function
	if restartCmd.Run == nil {
		t.Error("restartCmd should have a Run function")
	}
}

func TestRestartCmdFlagTypes(t *testing.T) {
	// Test that flags have correct types

	// String flags
	stringFlags := []string{"application", "environment", "role"}
	for _, flagName := range stringFlags {
		flag := restartCmd.Flags().Lookup(flagName)
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
		flag := restartCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Int flag '%s' should exist", flagName)
			continue
		}
		if flag.Value.Type() != "int" {
			t.Errorf("Expected flag '%s' to be int type, got '%s'", flagName, flag.Value.Type())
		}
	}

	// Bool flags
	boolFlags := []string{"no-wait"}
	for _, flagName := range boolFlags {
		flag := restartCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Bool flag '%s' should exist", flagName)
			continue
		}
		if flag.Value.Type() != "bool" {
			t.Errorf("Expected flag '%s' to be bool type, got '%s'", flagName, flag.Value.Type())
		}
	}
}
