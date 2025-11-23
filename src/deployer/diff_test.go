package deployer

import (
	"strings"
	"testing"
)

func TestNewDiff(t *testing.T) {
	resource := "container"
	name := "test-container"

	diff := NewDiff(resource, name)

	if diff.resource != resource {
		t.Errorf("Expected resource %s, got %s", resource, diff.resource)
	}

	if diff.name != name {
		t.Errorf("Expected name %s, got %s", name, diff.name)
	}

	if len(diff.changes) != 0 {
		t.Errorf("Expected empty changes slice, got %d changes", len(diff.changes))
	}
}

func TestDiff_AddChange(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		x        string
		y        string
		expected string
	}{
		{
			name:     "no change",
			key:      "image",
			x:        "nginx:1.0",
			y:        "nginx:1.0",
			expected: " \timage =\tnginx:1.0",
		},
		{
			name:     "addition",
			key:      "image",
			x:        "",
			y:        "nginx:1.0",
			expected: "+\timage =\tnginx:1.0",
		},
		{
			name:     "removal",
			key:      "image",
			x:        "nginx:1.0",
			y:        "",
			expected: "-\timage =\tnginx:1.0",
		},
		{
			name:     "modification",
			key:      "image",
			x:        "nginx:1.0",
			y:        "nginx:2.0",
			expected: "~\timage =\tnginx:1.0  -->  nginx:2.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff := NewDiff("container", "test")
			diff.AddChange(tt.key, tt.x, tt.y)

			if len(diff.changes) != 1 {
				t.Errorf("Expected 1 change, got %d", len(diff.changes))
				return
			}

			// Remove ANSI color codes for testing
			change := removeAnsiCodes(diff.changes[0])
			if change != tt.expected {
				t.Errorf("Expected change %q, got %q", tt.expected, change)
			}
		})
	}
}

func TestDiff_String(t *testing.T) {
	diff := NewDiff("container", "test-container")
	diff.AddChange("image", "nginx:1.0", "nginx:2.0")
	diff.AddChange("port", "80", "8080")

	result := diff.String()

	// Check that the output contains the expected structure
	if !strings.Contains(result, "container \"test-container\" {") {
		t.Error("Expected output to contain resource and name header")
	}

	if !strings.Contains(result, "}") {
		t.Error("Expected output to contain closing brace")
	}

	// Check that changes are included
	if !strings.Contains(result, "image") {
		t.Error("Expected output to contain image change")
	}

	if !strings.Contains(result, "port") {
		t.Error("Expected output to contain port change")
	}
}

func TestDiff_MultipleChanges(t *testing.T) {
	diff := NewDiff("service", "test-service")

	// Add multiple changes
	diff.AddChange("image", "nginx:1.0", "nginx:2.0")
	diff.AddChange("cpu", "256", "512")
	diff.AddChange("memory", "512", "1024")
	diff.AddChange("new_env", "", "production")
	diff.AddChange("old_env", "development", "")

	if len(diff.changes) != 5 {
		t.Errorf("Expected 5 changes, got %d", len(diff.changes))
	}

	result := diff.String()

	// Verify all changes are present in output
	expectedChanges := []string{"image", "cpu", "memory", "new_env", "old_env"}
	for _, change := range expectedChanges {
		if !strings.Contains(result, change) {
			t.Errorf("Expected output to contain change for %s", change)
		}
	}
}

// Helper function to remove ANSI color codes for testing
func removeAnsiCodes(s string) string {
	// Simple regex replacement would be better, but for testing we can use string replacement
	s = strings.ReplaceAll(s, "\x1b[37m", "") // White
	s = strings.ReplaceAll(s, "\x1b[32m", "") // Green
	s = strings.ReplaceAll(s, "\x1b[31m", "") // Red
	s = strings.ReplaceAll(s, "\x1b[33m", "") // Yellow
	s = strings.ReplaceAll(s, "\x1b[0m", "")  // Reset
	return s
}

func TestDiff_EmptyChanges(t *testing.T) {
	diff := NewDiff("container", "empty-container")

	result := diff.String()

	// Should still contain the header and closing brace
	if !strings.Contains(result, "container \"empty-container\" {") {
		t.Error("Expected output to contain resource and name header")
	}

	if !strings.Contains(result, "}") {
		t.Error("Expected output to contain closing brace")
	}

	// Should not contain any change indicators
	changeIndicators := []string{"+\t", "-\t", "~\t", " \t"}
	for _, indicator := range changeIndicators {
		if strings.Contains(result, indicator) {
			t.Errorf("Expected no changes, but found indicator: %s", indicator)
		}
	}
}
