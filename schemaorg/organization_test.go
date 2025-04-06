package schemaorg

import (
	"html/template"
	"testing"
)

func TestNewOrganization_SetsFieldsAndDefaults(t *testing.T) {
	org := NewOrganization(
		"Example Organization",
		"https://example.com",
		"https://example.com/logo.png",
		[]ContactPoint{{Telephone: "+1-800-1234", ContactType: "customer support"}},
		[]string{"https://twitter.com/example"},
	)

	if org.Context != "https://schema.org" {
		t.Errorf("expected context to be schema.org, got %s", org.Context)
	}
	if org.Type != "Organization" {
		t.Errorf("expected type to be Organization, got %s", org.Type)
	}
	if org.Name != "Example Organization" {
		t.Errorf("name not set correctly")
	}
	if org.URL != "https://example.com" {
		t.Errorf("URL not set correctly")
	}
	if org.Logo == nil || org.Logo.URL != "https://example.com/logo.png" || org.Logo.Type != "ImageObject" {
		t.Errorf("Logo not set correctly")
	}
	if len(org.ContactPoints) != 1 || org.ContactPoints[0].Telephone != "+1-800-1234" {
		t.Errorf("ContactPoints not set correctly")
	}
	if len(org.SameAs) != 1 || org.SameAs[0] != "https://twitter.com/example" {
		t.Errorf("SameAs not set correctly")
	}
}

func TestOrganization_EnsureDefaults_NilSafe(t *testing.T) {
	org := &Organization{}
	org.ensureDefaults()

	if org.Context != "https://schema.org" {
		t.Errorf("expected context to be schema.org, got %s", org.Context)
	}
	if org.Type != "Organization" {
		t.Errorf("expected type to be Organization, got %s", org.Type)
	}
}

func TestOrganization_Validate_AllGood(t *testing.T) {
	org := &Organization{
		Name: "Example Org",
		URL:  "https://example.com",
		Logo: &ImageObject{URL: "https://example.com/logo.jpg"},
	}
	warnings := org.Validate()
	if len(warnings) != 0 {
		t.Errorf("expected no warnings, got %v", warnings)
	}
}

func TestOrganization_Validate_MissingName(t *testing.T) {
	org := &Organization{
		URL:  "https://example.com",
		Logo: &ImageObject{URL: "https://example.com/logo.jpg"},
	}
	warnings := org.Validate()
	want := "missing recommended field: name"
	if len(warnings) != 1 || warnings[0] != want {
		t.Errorf("expected [%s], got %v", want, warnings)
	}
}

func TestOrganization_Validate_MissingURL(t *testing.T) {
	org := &Organization{
		Name: "Org",
		Logo: &ImageObject{URL: "https://example.com/logo.jpg"},
	}
	warnings := org.Validate()
	want := "missing recommended field: url"
	if len(warnings) != 1 || warnings[0] != want {
		t.Errorf("expected [%s], got %v", want, warnings)
	}
}

func TestOrganization_Validate_MissingLogo(t *testing.T) {
	org := &Organization{
		Name: "Org",
		URL:  "https://example.com",
	}
	warnings := org.Validate()
	want := "missing recommended field: logo.url"
	if len(warnings) != 1 || warnings[0] != want {
		t.Errorf("expected [%s], got %v", want, warnings)
	}
}

func TestOrganization_Validate_EmptyLogoURL(t *testing.T) {
	org := &Organization{
		Name: "Org",
		URL:  "https://example.com",
		Logo: &ImageObject{URL: ""},
	}
	warnings := org.Validate()
	want := "missing recommended field: logo.url"
	if len(warnings) != 1 || warnings[0] != want {
		t.Errorf("expected [%s], got %v", want, warnings)
	}
}

func TestOrganization_Validate_MultipleWarnings(t *testing.T) {
	org := &Organization{}
	warnings := org.Validate()

	expected := map[string]bool{
		"missing recommended field: name":     true,
		"missing recommended field: url":      true,
		"missing recommended field: logo.url": true,
	}

	if len(warnings) != len(expected) {
		t.Errorf("expected %d warnings, got %d: %v", len(expected), len(warnings), warnings)
	}

	for _, w := range warnings {
		if !expected[w] {
			t.Errorf("unexpected warning: %s", w)
		}
	}
}

func TestOrganization_ToGoHTMLJsonLd(t *testing.T) {
	org := NewOrganization("Test Org", "https://example.org", "https://example.org/logo.jpg", nil, nil)
	html, err := org.ToGoHTMLJsonLd()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if html == template.HTML("") {
		t.Errorf("expected non-empty HTML")
	}
}
