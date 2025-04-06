package schemaorg

import (
	"html/template"
	"testing"
)

func TestNewWebSite_SetsDefaults(t *testing.T) {
	ws := NewWebSite(
		"https://example.com",
		"Example Website",
		"Alt Name",
		"A description",
		nil,
	)

	if ws.Context != "https://schema.org" {
		t.Errorf("expected context to be schema.org, got %s", ws.Context)
	}
	if ws.Type != "WebSite" {
		t.Errorf("expected type to be WebSite, got %s", ws.Type)
	}
	if ws.Name != "Example Website" || ws.URL != "https://example.com" {
		t.Errorf("fields not set correctly")
	}
}

func TestWebSite_ToGoHTMLJsonLd(t *testing.T) {
	ws := NewWebSite("https://example.com", "Example", "", "", nil)
	html, err := ws.ToGoHTMLJsonLd()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if html == template.HTML("") {
		t.Errorf("expected non-empty HTML output")
	}
}

func TestWebSite_Validate(t *testing.T) {
	tests := []struct {
		name     string
		ws       *WebSite
		expected []string
	}{
		{
			name: "valid website",
			ws: &WebSite{
				URL:         "https://example.com",
				Name:        "My Site",
				Description: "Nice site",
			},
			expected: nil,
		},
		{
			name:     "missing all",
			ws:       &WebSite{},
			expected: []string{"missing recommended field: url", "missing recommended field: name", "missing recommended field: description"},
		},
		{
			name: "potentialAction without target",
			ws: &WebSite{
				URL:             "https://example.com",
				Name:            "My Site",
				Description:     "Cool site",
				PotentialAction: &Action{},
			},
			expected: []string{"potentialAction.target.urlTemplate is recommended when potentialAction is set"},
		},
		{
			name: "potentialAction with empty target",
			ws: &WebSite{
				URL:         "https://example.com",
				Name:        "My Site",
				Description: "Cool site",
				PotentialAction: &Action{
					Target: &Target{},
				},
			},
			expected: []string{"potentialAction.target.urlTemplate is recommended when potentialAction is set"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warnings := tt.ws.Validate()
			if len(warnings) != len(tt.expected) {
				t.Errorf("expected %d warnings, got %d: %v", len(tt.expected), len(warnings), warnings)
			}
			for _, expected := range tt.expected {
				found := false
				for _, w := range warnings {
					if w == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("missing expected warning: %s", expected)
				}
			}
		})
	}
}

func TestWebSite_EnsureDefaults_WithPotentialActionAndTarget(t *testing.T) {
	// Arrange
	target := &Target{}
	action := &Action{
		Target: target,
	}
	ws := &WebSite{
		PotentialAction: action,
	}

	// Act
	ws.ensureDefaults()

	// Assert
	if ws.Context != "https://schema.org" {
		t.Errorf("expected context to be schema.org, got %s", ws.Context)
	}
	if ws.Type != "WebSite" {
		t.Errorf("expected type to be WebSite, got %s", ws.Type)
	}
	if ws.PotentialAction.Type != "Action" {
		t.Errorf("expected Action type to be Action, got %s", ws.PotentialAction.Type)
	}
	if ws.PotentialAction.Target.Type != "EntryPoint" {
		t.Errorf("expected Target type to be EntryPoint, got %s", ws.PotentialAction.Target.Type)
	}
}
