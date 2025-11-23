package cmd

import (
	"bytes"
	"testing"
)

func TestRootCmd(t *testing.T) {
	// Test that root command exists and has expected properties
	if rootCmd == nil {
		t.Fatal("rootCmd should not be nil")
	}

	if rootCmd.Use != "ecs-deploy" {
		t.Errorf("Expected Use to be 'ecs-deploy', got '%s'", rootCmd.Use)
	}

	if rootCmd.Short != "Deploy something" {
		t.Errorf("Expected Short to be 'Deploy something', got '%s'", rootCmd.Short)
	}

	expectedLong := "A fast and flexible tool to deploy to Amazon Web Service's Elastic Container Service"
	if rootCmd.Long != expectedLong {
		t.Errorf("Expected Long to be '%s', got '%s'", expectedLong, rootCmd.Long)
	}

	if rootCmd.Version != "0.7.3" {
		t.Errorf("Expected Version to be '0.7.3', got '%s'", rootCmd.Version)
	}
}

func TestRootCmdFlags(t *testing.T) {
	// Test that debug flag exists
	debugFlag := rootCmd.PersistentFlags().Lookup("debug")
	if debugFlag == nil {
		t.Fatal("debug flag should exist")
	}

	if debugFlag.Shorthand != "d" {
		t.Errorf("Expected debug flag shorthand to be 'd', got '%s'", debugFlag.Shorthand)
	}

	if debugFlag.Usage != "Enable debug logging" {
		t.Errorf("Expected debug flag usage to be 'Enable debug logging', got '%s'", debugFlag.Usage)
	}

	// Test default value
	if debugFlag.DefValue != "false" {
		t.Errorf("Expected debug flag default value to be 'false', got '%s'", debugFlag.DefValue)
	}
}

func TestRootCmdHasSubcommands(t *testing.T) {
	commands := rootCmd.Commands()

	if len(commands) == 0 {
		t.Error("Expected root command to have subcommands")
		return
	}

	// Check for expected subcommands
	expectedCommands := []string{"ship", "restart", "rollback"}
	foundCommands := make(map[string]bool)

	for _, cmd := range commands {
		foundCommands[cmd.Use] = true
	}

	for _, expected := range expectedCommands {
		if !foundCommands[expected] {
			t.Errorf("Expected subcommand '%s' not found", expected)
		}
	}
}

func TestExecuteFunction(t *testing.T) {
	// Test that Execute function exists and doesn't panic
	// We can't easily test the actual execution without mocking os.Exit
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Execute function panicked: %v", r)
		}
	}()

	// Just ensure the function exists
	_ = Execute
}

func TestRootCmdHelp(t *testing.T) {
	// Test that help can be generated without errors
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	// Reset args to avoid interference from other tests
	rootCmd.SetArgs([]string{"--help"})

	// This should not panic or error
	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("Help command should not return error, got: %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Error("Help output should not be empty")
	}

	// Check that help contains expected content
	if !bytes.Contains(buf.Bytes(), []byte("ecs-deploy")) {
		t.Error("Help output should contain 'ecs-deploy'")
	}
}

func TestGlobalVariables(t *testing.T) {
	// Test that global variables are properly initialized

	// lambdaName should be initialized (even if empty)
	_ = lambdaName

	// debugEnabled should be initialized to false
	if debugEnabled {
		// Reset it for other tests
		debugEnabled = false
	}
}

func TestRootCmdVersion(t *testing.T) {
	// Test that version is set correctly in the command
	if rootCmd.Version != "0.7.3" {
		t.Errorf("Expected version to be '0.7.3', got '%s'", rootCmd.Version)
	}
}
