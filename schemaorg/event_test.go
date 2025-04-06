package schemaorg

import (
	"testing"
)

func TestNewEvent_SetsDefaults(t *testing.T) {
	event := NewEvent(
		"Sample Event", "desc", "2024-01-01T10:00:00", "2024-01-01T12:00:00",
		&Place{Name: "Venue"}, nil, nil, nil, "", "", nil,
	)

	if event.Type != "Event" {
		t.Errorf("expected type to be Event, got %s", event.Type)
	}
	if event.Context != "https://schema.org" {
		t.Errorf("expected context to be schema.org, got %s", event.Context)
	}
	if event.Location.Type != "Place" {
		t.Errorf("expected nested Place type to be set")
	}
}

func TestEvent_EnsureDefaults(t *testing.T) {
	e := &Event{}
	e.ensureDefaults()
	if e.Type != "Event" {
		t.Errorf("expected default type Event, got %s", e.Type)
	}
	if e.Context != "https://schema.org" {
		t.Errorf("expected default context https://schema.org, got %s", e.Context)
	}
}

func TestEvent_EnsureDefaults_WithNestedValues(t *testing.T) {
	org := &Organization{}
	per := &Person{}
	offer := &Offer{}

	event := &Event{
		Organizer: org,
		Performer: per,
		Offers:    offer,
	}

	event.ensureDefaults()

	if org.Context != "https://schema.org" {
		t.Errorf("expected Organizer context to be schema.org, got %s", org.Context)
	}
	if org.Type != "Organization" {
		t.Errorf("expected Organizer type to be Organization, got %s", org.Type)
	}

	if per.Context != "https://schema.org" {
		t.Errorf("expected Performer context to be schema.org, got %s", per.Context)
	}
	if per.Type != "Person" {
		t.Errorf("expected Performer type to be Person, got %s", per.Type)
	}

	if offer.Type != "Offer" {
		t.Errorf("expected Offer type to be Offer, got %s", offer.Type)
	}
}

func TestEvent_Validate_AllGood(t *testing.T) {
	e := &Event{
		Name:      "My Event",
		StartDate: "2024-01-01T10:00:00",
		Location:  &Place{Name: "Venue"},
	}
	warnings := e.Validate()
	if len(warnings) != 0 {
		t.Errorf("expected no warnings, got %v", warnings)
	}
}

func TestEvent_Validate_MissingName(t *testing.T) {
	e := &Event{
		StartDate: "2024-01-01T10:00:00",
		Location:  &Place{Name: "Venue"},
	}
	w := e.Validate()
	if len(w) != 1 || w[0] != "missing recommended field: name" {
		t.Errorf("expected name warning, got %v", w)
	}
}

func TestEvent_Validate_MissingStartDate(t *testing.T) {
	e := &Event{
		Name:     "My Event",
		Location: &Place{Name: "Venue"},
	}
	w := e.Validate()
	if len(w) != 1 || w[0] != "missing recommended field: startDate" {
		t.Errorf("expected startDate warning, got %v", w)
	}
}

func TestEvent_Validate_MissingLocation(t *testing.T) {
	e := &Event{
		Name:      "My Event",
		StartDate: "2024-01-01T10:00:00",
	}
	w := e.Validate()
	if len(w) != 1 || w[0] != "missing recommended field: location" {
		t.Errorf("expected location warning, got %v", w)
	}
}

func TestEvent_Validate_AllMissing(t *testing.T) {
	e := &Event{}
	expected := map[string]bool{
		"missing recommended field: name":      true,
		"missing recommended field: startDate": true,
		"missing recommended field: location":  true,
	}
	warnings := e.Validate()
	if len(warnings) != len(expected) {
		t.Errorf("expected %d warnings, got %d", len(expected), len(warnings))
	}
	for _, w := range warnings {
		if !expected[w] {
			t.Errorf("unexpected warning: %s", w)
		}
	}
}

func TestEvent_ToGoHTMLJsonLd(t *testing.T) {
	e := NewEvent("Test Event", "", "2024-01-01T10:00", "", &Place{Name: "Venue"}, nil, nil, nil, "", "", nil)
	html, err := e.ToGoHTMLJsonLd()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if html == "" {
		t.Errorf("expected non-empty HTML output")
	}
}

func TestPlace_EnsureDefaults_WithAddressAndGeo(t *testing.T) {
	addr := &PostalAddress{}
	geo := &GeoCoordinates{}
	p := &Place{
		Address: addr,
		Geo:     geo,
	}
	p.ensureDefaults()

	if addr.Type != "PostalAddress" {
		t.Errorf("expected Address.Type to be PostalAddress")
	}
	if geo.Type != "GeoCoordinates" {
		t.Errorf("expected Geo.Type to be GeoCoordinates")
	}
}

func TestPlace_EnsureDefaults_NilAddressAndGeo(t *testing.T) {
	p := &Place{
		Address: nil,
		Geo:     nil,
	}
	p.ensureDefaults()

	if p.Context != "https://schema.org" {
		t.Errorf("expected Context to be https://schema.org, got %s", p.Context)
	}
	if p.Type != "Place" {
		t.Errorf("expected Type to be Place, got %s", p.Type)
	}
	// This test ensures no panic and default values are set even if Address and Geo are nil.
}
