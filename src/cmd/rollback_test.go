package cmd

import (
	"testing"
)

func TestRollbackCmd(t *testing.T) {
	// Test that rollback command exists and has expected properties
	if rollbackCmd == nil {
		t.Fatal("rollbackCmd should not be nil")
	}

	if rollbackCmd.Use != "rollback" {
		t.Errorf("Expected Use to be 'rollback', got '%s'", rollbackCmd.Use)
	}

	if rollbackCmd.Short != "Redeploy the previous task definition" {
		t.Errorf("Expected Short to be 'Redeploy the previous task definition', got '%s'", rollbackCmd.Short)
	}
}

func TestRollbackCmdFlags(t *testing.T) {
	// Test required flags
	requiredFlags := []struct {
		name      string
		shorthand string
		usage     string
	}{
		{"application", "a", "Application name to rollback"},
		{"environment", "e", "Target environment for rollback"},
	}

	for _, flag := range requiredFlags {
		f := rollbackCmd.PersistentFlags().Lookup(flag.name)
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
	taskDefFlag := rollbackCmd.PersistentFlags().Lookup("task-definition")
	if taskDefFlag == nil {
		t.Error("task-definition flag should exist")
	} else {
		if taskDefFlag.Shorthand != "t" {
			t.Errorf("Expected task-definition flag shorthand to be 't', got '%s'", taskDefFlag.Shorthand)
		}
		if taskDefFlag.DefValue != "0" {
			t.Errorf("Expected task-definition flag default value to be '0', got '%s'", taskDefFlag.DefValue)
		}
	}
}

func TestRollbackCmdHasRunFunction(t *testing.T) {
	// Test that rollback command has a Run function
	if rollbackCmd.Run == nil {
		t.Error("rollbackCmd should have a Run function")
	}
}

func TestRollbackCmdRunFunction(t *testing.T) {
	// Test that the run function exists and doesn't panic
	// The actual output goes to fmt.Println which we can't easily capture in unit tests
	// This is better tested in integration tests

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Rollback run function panicked: %v", r)
		}
	}()

	// Just ensure the function can be called without panicking
	rollbackCmd.Run(rollbackCmd, []string{})
}

func TestRollbackCmdFlagTypes(t *testing.T) {
	// Test that flags have correct types

	// String flags
	stringFlags := []string{"application", "environment"}
	for _, flagName := range stringFlags {
		flag := rollbackCmd.PersistentFlags().Lookup(flagName)
		if flag == nil {
			t.Errorf("String flag '%s' should exist", flagName)
			continue
		}
		if flag.Value.Type() != "string" {
			t.Errorf("Expected flag '%s' to be string type, got '%s'", flagName, flag.Value.Type())
		}
	}

	// Int flags
	intFlags := []string{"task-definition"}
	for _, flagName := range intFlags {
		flag := rollbackCmd.PersistentFlags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Int flag '%s' should exist", flagName)
			continue
		}
		if flag.Value.Type() != "int" {
			t.Errorf("Expected flag '%s' to be int type, got '%s'", flagName, flag.Value.Type())
		}
	}
}

func TestRollbackCmdDefaultValues(t *testing.T) {
	// Test default values for flags
	appFlag := rollbackCmd.PersistentFlags().Lookup("application")
	if appFlag != nil && appFlag.DefValue != "" {
		t.Errorf("Expected application flag default value to be empty, got '%s'", appFlag.DefValue)
	}

	envFlag := rollbackCmd.PersistentFlags().Lookup("environment")
	if envFlag != nil && envFlag.DefValue != "" {
		t.Errorf("Expected environment flag default value to be empty, got '%s'", envFlag.DefValue)
	}

	taskDefFlag := rollbackCmd.PersistentFlags().Lookup("task-definition")
	if taskDefFlag != nil && taskDefFlag.DefValue != "0" {
		t.Errorf("Expected task-definition flag default value to be '0', got '%s'", taskDefFlag.DefValue)
	}
}
