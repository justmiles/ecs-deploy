//go:build integration
// +build integration

package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// Integration tests require the binary to be built
// Run with: go test -tags=integration

func TestMain(m *testing.M) {
	// Build the binary before running integration tests
	cmd := exec.Command("go", "build", "-o", "build/ecs-deploy-test", "./src")
	if err := cmd.Run(); err != nil {
		panic("Failed to build binary for integration tests: " + err.Error())
	}

	// Run tests
	code := m.Run()

	// Cleanup
	os.Remove("build/ecs-deploy-test")
	os.Exit(code)
}

func TestCLI_Help(t *testing.T) {
	cmd := exec.Command("./build/ecs-deploy-test", "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Help command failed: %v", err)
	}

	outputStr := string(output)

	// Check for expected content in help output
	expectedContent := []string{
		"ecs-deploy",
		"A fast and flexible tool to deploy to Amazon Web Service's Elastic Container Service",
		"Available Commands:",
		"ship",
		"restart",
		"rollback",
		"--debug",
		"--help",
		"--version",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(outputStr, expected) {
			t.Errorf("Help output missing expected content: %s", expected)
		}
	}
}

func TestCLI_Version(t *testing.T) {
	cmd := exec.Command("./build/ecs-deploy-test", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Version command failed: %v", err)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "0.7.3") {
		t.Errorf("Version output should contain version number, got: %s", outputStr)
	}
}

func TestCLI_ShipHelp(t *testing.T) {
	cmd := exec.Command("./build/ecs-deploy-test", "ship", "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Ship help command failed: %v", err)
	}

	outputStr := string(output)

	// Check for expected flags in ship help
	expectedFlags := []string{
		"--application",
		"--version",
		"--environment",
		"--role",
		"--max-attempts",
		"--no-wait",
		"--refresh-secrets",
		"--dry-run",
		"--secrets-prefix",
		"--ignore-tags",
	}

	for _, flag := range expectedFlags {
		if !strings.Contains(outputStr, flag) {
			t.Errorf("Ship help output missing expected flag: %s", flag)
		}
	}
}

func TestCLI_RestartHelp(t *testing.T) {
	cmd := exec.Command("./build/ecs-deploy-test", "restart", "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Restart help command failed: %v", err)
	}

	outputStr := string(output)

	// Check for expected flags in restart help
	expectedFlags := []string{
		"--application",
		"--environment",
		"--role",
		"--max-attempts",
		"--no-wait",
	}

	for _, flag := range expectedFlags {
		if !strings.Contains(outputStr, flag) {
			t.Errorf("Restart help output missing expected flag: %s", flag)
		}
	}
}

func TestCLI_RollbackHelp(t *testing.T) {
	cmd := exec.Command("./build/ecs-deploy-test", "rollback", "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Rollback help command failed: %v", err)
	}

	outputStr := string(output)

	// Check for expected flags in rollback help
	expectedFlags := []string{
		"--application",
		"--environment",
		"--task-definition",
	}

	for _, flag := range expectedFlags {
		if !strings.Contains(outputStr, flag) {
			t.Errorf("Rollback help output missing expected flag: %s", flag)
		}
	}
}

func TestCLI_ShipMissingRequiredFlags(t *testing.T) {
	cmd := exec.Command("./build/ecs-deploy-test", "ship")
	output, err := cmd.CombinedOutput()

	// Should exit with error due to missing required flags
	if err == nil {
		t.Error("Ship command should fail when required flags are missing")
	}

	outputStr := string(output)

	// Should mention required flags
	requiredFlags := []string{"application", "version", "environment"}
	foundMissingFlag := false

	for _, flag := range requiredFlags {
		if strings.Contains(outputStr, flag) {
			foundMissingFlag = true
			break
		}
	}

	if !foundMissingFlag {
		t.Errorf("Error output should mention missing required flags, got: %s", outputStr)
	}
}

func TestCLI_RestartMissingRequiredFlags(t *testing.T) {
	cmd := exec.Command("./build/ecs-deploy-test", "restart")
	output, err := cmd.CombinedOutput()

	// Should exit with error due to missing required flags
	if err == nil {
		t.Error("Restart command should fail when required flags are missing")
	}

	outputStr := string(output)

	// Should mention required flags
	requiredFlags := []string{"application", "environment"}
	foundMissingFlag := false

	for _, flag := range requiredFlags {
		if strings.Contains(outputStr, flag) {
			foundMissingFlag = true
			break
		}
	}

	if !foundMissingFlag {
		t.Errorf("Error output should mention missing required flags, got: %s", outputStr)
	}
}

func TestCLI_RollbackOutput(t *testing.T) {
	// Test rollback with required flags to see the "not implemented" message
	cmd := exec.Command("./build/ecs-deploy-test", "rollback", "--application", "test", "--environment", "test")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Rollback command failed: %v", err)
	}

	outputStr := string(output)
	expectedOutput := "Rollback functionality not yet implemented"

	if !strings.Contains(outputStr, expectedOutput) {
		t.Errorf("Expected rollback output to contain '%s', got: %s", expectedOutput, outputStr)
	}
}

func TestCLI_DebugFlag(t *testing.T) {
	// Test that debug flag is accepted (though we can't easily test its effect)
	cmd := exec.Command("./build/ecs-deploy-test", "--debug", "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Debug flag should be accepted: %v", err)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "ecs-deploy") {
		t.Error("Debug flag should not prevent help output")
	}
}
