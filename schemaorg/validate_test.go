package schemaorg

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

// mockValidator is a helper for simulating a SchemaValidator
type mockValidator struct {
	warnings []string
}

func (m *mockValidator) Validate() []string {
	return m.warnings
}

func TestLogValidationWarnings(t *testing.T) {
	tests := []struct {
		name      string
		validator SchemaValidator
		expected  []string
	}{
		{
			name:      "no warnings",
			validator: &mockValidator{warnings: []string{}},
			expected:  nil,
		},
		{
			name:      "single warning",
			validator: &mockValidator{warnings: []string{"missing headline"}},
			expected:  []string{"schema warning: missing headline"},
		},
		{
			name: "multiple warnings",
			validator: &mockValidator{warnings: []string{
				"missing headline", "missing image", "missing datePublished",
			}},
			expected: []string{
				"schema warning: missing headline",
				"schema warning: missing image",
				"schema warning: missing datePublished",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(nil) // restore default output

			LogValidationWarnings(tt.validator)

			output := buf.String()
			for _, expected := range tt.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("expected log to contain %q, but got: %q", expected, output)
				}
			}

			// Optional: ensure nothing is logged when expected is nil
			if len(tt.expected) == 0 && output != "" {
				t.Errorf("expected no output, got: %q", output)
			}
		})
	}
}
