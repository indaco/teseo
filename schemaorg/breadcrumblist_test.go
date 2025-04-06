package schemaorg

import (
	"slices"
	"testing"
)

func TestNewBreadcrumbList_SetsDefaults(t *testing.T) {
	items := []ListItem{
		{Name: "Home", Item: "https://example.com", Position: 1},
	}
	bc := NewBreadcrumbList(items)

	if bc.Type != "BreadcrumbList" {
		t.Errorf("expected type BreadcrumbList, got %s", bc.Type)
	}
	if bc.Context != "https://schema.org" {
		t.Errorf("expected context schema.org, got %s", bc.Context)
	}
	if bc.ItemListElement[0].Type != "ListItem" {
		t.Errorf("expected ListItem type to be ListItem, got %s", bc.ItemListElement[0].Type)
	}
}

func TestBreadcrumbList_EnsureDefaults(t *testing.T) {
	bc := &BreadcrumbList{
		ItemListElement: []ListItem{{}},
	}
	bc.ensureDefaults()

	if bc.Context != "https://schema.org" {
		t.Errorf("expected context to be set")
	}
	if bc.Type != "BreadcrumbList" {
		t.Errorf("expected type to be set")
	}
	if bc.ItemListElement[0].Type != "ListItem" {
		t.Errorf("expected item type to be ListItem")
	}
}

func TestBreadcrumbList_ToGoHTMLJsonLd(t *testing.T) {
	bc := NewBreadcrumbList([]ListItem{
		{Name: "Home", Item: "https://example.com", Position: 1},
	})
	html, err := bc.ToGoHTMLJsonLd()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if html == "" {
		t.Errorf("expected non-empty html")
	}
}

func TestNewBreadcrumbListFromUrl(t *testing.T) {
	bc, err := NewBreadcrumbListFromUrl("https://example.com/about/team")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(bc.ItemListElement) != 3 {
		t.Errorf("expected 3 items, got %d", len(bc.ItemListElement))
	}
	if bc.ItemListElement[1].Name != "About" || bc.ItemListElement[2].Name != "Team" {
		t.Errorf("unexpected breadcrumb names: %+v", bc.ItemListElement)
	}
}

func TestNewBreadcrumbListFromUrl_Invalid(t *testing.T) {
	_, err := NewBreadcrumbListFromUrl("http://[::1]:namedport")
	if err == nil {
		t.Errorf("expected error for invalid URL")
	}
}

func TestBreadcrumbList_Validate(t *testing.T) {
	tests := []struct {
		name     string
		list     *BreadcrumbList
		expected []string
	}{
		{
			name: "valid breadcrumb",
			list: NewBreadcrumbList([]ListItem{
				{Name: "Home", Item: "https://example.com", Position: 1},
			}),
			expected: nil,
		},
		{
			name: "missing item list",
			list: &BreadcrumbList{},
			expected: []string{
				"BreadcrumbList should contain at least one item",
			},
		},
		{
			name: "item with missing name",
			list: &BreadcrumbList{
				ItemListElement: []ListItem{
					{Item: "https://example.com", Position: 1},
				},
			},
			expected: []string{
				"ListItem at position 1 is missing a name",
			},
		},
		{
			name: "item with missing item URL",
			list: &BreadcrumbList{
				ItemListElement: []ListItem{
					{Name: "Home", Position: 1},
				},
			},
			expected: []string{
				"ListItem at position 1 is missing a URL",
			},
		},
		{
			name: "item with missing name and item",
			list: &BreadcrumbList{
				ItemListElement: []ListItem{
					{Position: 1},
				},
			},
			expected: []string{
				"ListItem at position 1 is missing a name",
				"ListItem at position 1 is missing a URL",
			},
		},
		{
			name: "item with missing position",
			list: &BreadcrumbList{
				ItemListElement: []ListItem{
					{Name: "Home", Item: "https://example.com", Position: 0},
				},
			},
			expected: []string{
				"ListItem at position 1 is missing a valid position",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warnings := tt.list.Validate()
			if len(warnings) != len(tt.expected) {
				t.Errorf("expected %d warnings, got %d: %v", len(tt.expected), len(warnings), warnings)
				return
			}
			for _, expectedWarning := range tt.expected {
				found := slices.Contains(warnings, expectedWarning)
				if !found {
					t.Errorf("expected warning %q not found in %v", expectedWarning, warnings)
				}
			}
		})
	}
}

func TestToTitle(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"hello", "Hello"},
		{"ßtraße", "ßtraße"},
		{"123abc", "123abc"},
	}

	for _, tt := range tests {
		result := toTitle(tt.input)
		if result != tt.expected {
			t.Errorf("toTitle(%q) = %q; want %q", tt.input, result, tt.expected)
		}
	}
}
