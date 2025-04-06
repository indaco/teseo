package schemaorg

import (
	"html/template"
	"slices"
	"testing"
)

func TestNewWebPage_SetsFieldsAndDefaults(t *testing.T) {
	wp := NewWebPage(
		"https://example.com",
		"Example Page",
		"Headline Here",
		"Description here",
		"About something",
		"keywords, go, templ",
		"en",
		"Main site",
		"2023-01-01",
		"https://example.com/image.jpg",
		"2023-01-01",
		"2023-01-02",
	)

	if wp.Context != "https://schema.org" {
		t.Errorf("expected default context to be schema.org, got %s", wp.Context)
	}
	if wp.Type != "WebPage" {
		t.Errorf("expected default type to be WebPage, got %s", wp.Type)
	}
	if wp.URL != "https://example.com" {
		t.Errorf("URL not set correctly")
	}
	if wp.Name != "Example Page" {
		t.Errorf("Name not set correctly")
	}
	if wp.Headline != "Headline Here" {
		t.Errorf("Headline not set correctly")
	}
}

func TestWebPage_EnsureDefaults(t *testing.T) {
	wp := &WebPage{}
	wp.ensureDefaults()

	if wp.Context != "https://schema.org" {
		t.Errorf("expected schema.org context, got %s", wp.Context)
	}
	if wp.Type != "WebPage" {
		t.Errorf("expected type to be WebPage, got %s", wp.Type)
	}
}

func TestWebPage_ToGoHTMLJsonLd(t *testing.T) {
	wp := NewWebPage("https://test.com", "Test", "Headline", "Desc", "", "", "", "", "", "", "", "")
	html, err := wp.ToGoHTMLJsonLd()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if html == template.HTML("") {
		t.Errorf("expected non-empty HTML output")
	}
}

func TestWebPage_Validate_AllGood(t *testing.T) {
	wp := &WebPage{
		URL:         "https://example.com",
		Name:        "Example",
		Headline:    "Headline",
		Description: "Description",
	}
	warnings := wp.Validate()
	if len(warnings) != 0 {
		t.Errorf("expected no warnings, got %v", warnings)
	}
}

func TestWebPage_Validate_MissingFields(t *testing.T) {
	wp := &WebPage{}
	got := wp.Validate()
	expected := []string{
		"missing recommended field: url",
		"missing recommended field: name",
		"missing recommended field: headline",
		"missing recommended field: description",
	}

	if len(got) != len(expected) {
		t.Errorf("expected %d warnings, got %d: %v", len(expected), len(got), got)
	}

	for _, want := range expected {
		found := slices.Contains(got, want)
		if !found {
			t.Errorf("expected warning %q not found in %v", want, got)
		}
	}
}
