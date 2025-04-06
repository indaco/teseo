package schemaorg

import (
	"slices"
	"testing"
)

func TestNewPerson_SetsFieldsAndDefaults(t *testing.T) {
	org := &Organization{Name: "Example Company"}
	img := &ImageObject{URL: "https://img.com/avatar.jpg"}

	person := NewPerson("Jane Doe", "https://example.com", "jane@example.com", img, "Engineer", org, []string{"https://twitter.com/jane"}, "Female", "1990-01-01", "US", "+123456", nil, nil)

	if person.Type != "Person" {
		t.Errorf("expected type Person, got %s", person.Type)
	}
	if person.Context != "https://schema.org" {
		t.Errorf("expected context schema.org, got %s", person.Context)
	}
	if person.Image.Type != "ImageObject" {
		t.Errorf("expected ImageObject type")
	}
	if person.WorksFor.Type != "Organization" {
		t.Errorf("expected Organization type")
	}
}

func TestPerson_EnsureDefaults_WithNestedValues(t *testing.T) {
	image := &ImageObject{}
	worksFor := &Organization{}
	address := &PostalAddress{}
	affiliation := &Organization{}

	person := &Person{
		Image:       image,
		WorksFor:    worksFor,
		Address:     address,
		Affiliation: affiliation,
	}

	person.ensureDefaults()

	if person.Context != "https://schema.org" {
		t.Errorf("expected Person context to be schema.org, got %s", person.Context)
	}
	if person.Type != "Person" {
		t.Errorf("expected Person type to be Person, got %s", person.Type)
	}

	if image.Type != "ImageObject" {
		t.Errorf("expected Image type to be ImageObject, got %s", image.Type)
	}

	if worksFor.Context != "https://schema.org" {
		t.Errorf("expected WorksFor context to be schema.org, got %s", worksFor.Context)
	}
	if worksFor.Type != "Organization" {
		t.Errorf("expected WorksFor type to be Organization, got %s", worksFor.Type)
	}

	if address.Type != "PostalAddress" {
		t.Errorf("expected Address type to be PostalAddress, got %s", address.Type)
	}

	if affiliation.Context != "https://schema.org" {
		t.Errorf("expected Affiliation context to be schema.org, got %s", affiliation.Context)
	}
	if affiliation.Type != "Organization" {
		t.Errorf("expected Affiliation type to be Organization, got %s", affiliation.Type)
	}
}

func TestPerson_ToGoHTMLJsonLd(t *testing.T) {
	person := NewPerson("Jane", "", "", nil, "", nil, nil, "", "", "", "", nil, nil)

	html, err := person.ToGoHTMLJsonLd()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if html == "" {
		t.Errorf("expected non-empty HTML")
	}
}

func TestPerson_Validate_AllGood(t *testing.T) {
	p := &Person{
		Name:     "John",
		Email:    "john@example.com",
		JobTitle: "Engineer",
	}
	warnings := p.Validate()
	if len(warnings) > 0 {
		t.Errorf("expected no warnings, got: %v", warnings)
	}
}

func TestPerson_Validate_MissingFields(t *testing.T) {
	tests := []struct {
		name     string
		person   *Person
		expected []string
	}{
		{
			name:     "missing all",
			person:   &Person{},
			expected: []string{"missing recommended field: name", "missing recommended field: email", "missing recommended field: jobTitle"},
		},
		{
			name:     "missing email and job title",
			person:   &Person{Name: "Jane"},
			expected: []string{"missing recommended field: email", "missing recommended field: jobTitle"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.person.Validate()
			if len(got) != len(tt.expected) {
				t.Fatalf("expected %d warnings, got %d: %v", len(tt.expected), len(got), got)
			}
			for _, expected := range tt.expected {
				found := slices.Contains(got, expected)
				if !found {
					t.Errorf("expected warning %q not found in %v", expected, got)
				}
			}
		})
	}
}
