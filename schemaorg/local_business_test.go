package schemaorg

import (
	"testing"
)

func TestNewLocalBusiness_SetsFieldsAndDefaults(t *testing.T) {
	lb := NewLocalBusiness(
		"Business Name",
		"Great service",
		"https://example.com",
		"+1-800-555-1234",
		&ImageObject{URL: "https://example.com/logo.png"},
		&PostalAddress{StreetAddress: "123 St", AddressLocality: "Town"},
		[]string{"Mo-Fr 09:00-17:00"},
		&GeoCoordinates{Latitude: 45.0, Longitude: 12.0},
		&AggregateRating{RatingValue: 4.5, ReviewCount: 100},
		[]*Review{{ReviewBody: "Awesome!"}},
	)

	if lb.Type != "LocalBusiness" {
		t.Errorf("expected type to be LocalBusiness, got %s", lb.Type)
	}
	if lb.Context != "https://schema.org" {
		t.Errorf("expected context to be schema.org, got %s", lb.Context)
	}
	if lb.Logo == nil || lb.Logo.Type != "ImageObject" {
		t.Errorf("expected logo with type ImageObject")
	}
	if lb.Address == nil || lb.Address.Type != "PostalAddress" {
		t.Errorf("expected address with type PostalAddress")
	}
	if lb.Geo == nil || lb.Geo.Type != "GeoCoordinates" {
		t.Errorf("expected geo with type GeoCoordinates")
	}
	if lb.AggregateRating == nil || lb.AggregateRating.Type != "AggregateRating" {
		t.Errorf("expected aggregateRating with type AggregateRating")
	}
	if len(lb.Review) != 1 || lb.Review[0].Type != "Review" {
		t.Errorf("expected review with type Review")
	}
}

func TestLocalBusiness_EnsureDefaults_NestedNilSafe(t *testing.T) {
	lb := &LocalBusiness{}
	lb.ensureDefaults()

	if lb.Type != "LocalBusiness" {
		t.Errorf("expected default type LocalBusiness")
	}
	if lb.Context != "https://schema.org" {
		t.Errorf("expected default context schema.org")
	}
}

func TestLocalBusiness_Validate(t *testing.T) {
	tests := []struct {
		name     string
		lb       *LocalBusiness
		expected []string
	}{
		{
			name: "valid",
			lb: &LocalBusiness{
				Name:        "Shop",
				Address:     &PostalAddress{},
				Telephone:   "+1-800",
				Description: "desc",
			},
			expected: nil,
		},
		{
			name:     "all missing",
			lb:       &LocalBusiness{},
			expected: []string{"missing recommended field: name", "missing recommended field: address", "missing recommended field: telephone", "missing recommended field: description"},
		},
		{
			name:     "missing address",
			lb:       &LocalBusiness{Name: "x", Telephone: "y", Description: "z"},
			expected: []string{"missing recommended field: address"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warnings := tt.lb.Validate()
			if len(warnings) != len(tt.expected) {
				t.Errorf("expected %d warnings, got %d: %v", len(tt.expected), len(warnings), warnings)
				return
			}
			for _, expected := range tt.expected {
				found := false
				for _, got := range warnings {
					if expected == got {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected warning %q not found in %v", expected, warnings)
				}
			}
		})
	}
}

func TestLocalBusiness_ToGoHTMLJsonLd(t *testing.T) {
	lb := NewLocalBusiness("My Shop", "", "", "", nil, nil, nil, nil, nil, nil)
	html, err := lb.ToGoHTMLJsonLd()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if html == "" {
		t.Errorf("expected non-empty html")
	}
}
