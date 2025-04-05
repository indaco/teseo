package opengraph

import (
	"context"
	"strings"
	"testing"
)

func TestNewPlace_SetsFieldsAndDefaults(t *testing.T) {
	place := NewPlace(
		"NYC Office",
		"https://example.com/place",
		"Main office location",
		"https://example.com/img.jpg",
		40.7128,
		-74.0060,
		"123 Main St",
		"New York",
		"NY",
		"10001",
		"USA",
	)

	if place.Type != "place" {
		t.Errorf("expected Type to be 'place', got '%s'", place.Type)
	}
	if place.Latitude != 40.7128 || place.Longitude != -74.0060 {
		t.Errorf("latitude or longitude not set correctly")
	}
	if place.StreetAddress != "123 Main St" {
		t.Errorf("expected StreetAddress to be set")
	}
}

func TestPlace_ensureDefaults(t *testing.T) {
	p := &Place{}
	p.ensureDefaults()
	if p.Type != "place" {
		t.Errorf("expected default type to be 'place'")
	}
}

func TestPlace_metaTags(t *testing.T) {
	p := NewPlace(
		"Colosseum",
		"https://rome.com/colosseum",
		"Historic Roman landmark",
		"https://rome.com/img.jpg",
		41.8902,
		12.4922,
		"Piazza del Colosseo, 1",
		"Rome",
		"RM",
		"00184",
		"Italy",
	)

	tags := p.metaTags()

	assertTag := func(prop, expected string) {
		for _, tag := range tags {
			if tag.property == prop && tag.content == expected {
				return
			}
		}
		t.Errorf("missing or incorrect tag: %s=%s", prop, expected)
	}

	assertTag("og:type", "place")
	assertTag("og:title", "Colosseum")
	assertTag("og:url", "https://rome.com/colosseum")
	assertTag("og:description", "Historic Roman landmark")
	assertTag("og:image", "https://rome.com/img.jpg")
	assertTag("place:location:latitude", "41.8902")
	assertTag("place:location:longitude", "12.4922")
	assertTag("place:contact_data:street_address", "Piazza del Colosseo, 1")
	assertTag("place:contact_data:locality", "Rome")
	assertTag("place:contact_data:region", "RM")
	assertTag("place:contact_data:postal_code", "00184")
	assertTag("place:contact_data:country_name", "Italy")
}

func TestPlace_metaTags_EmptyValuesAreRenderedIfZeroValue(t *testing.T) {
	p := &Place{
		OpenGraphObject: OpenGraphObject{Title: "Null Island"},
		Latitude:        0.0,
		Longitude:       0.0,
	}

	tags := p.metaTags()
	var foundLat, foundLon bool

	for _, tag := range tags {
		if tag.property == "place:location:latitude" && tag.content == "0.0000" {
			foundLat = true
		}
		if tag.property == "place:location:longitude" && tag.content == "0.0000" {
			foundLon = true
		}
	}

	if !foundLat {
		t.Errorf("expected latitude to be rendered as 0.0000")
	}
	if !foundLon {
		t.Errorf("expected longitude to be rendered as 0.0000")
	}
}

func TestPlace_ToMetaTags_WriteError(t *testing.T) {
	p := NewPlace(
		"Times Square",
		"https://example.com/ts",
		"Famous NYC location",
		"https://example.com/img.png",
		40.7580,
		-73.9855,
		"Broadway & 7th",
		"New York",
		"NY",
		"10036",
		"USA",
	)
	p.ensureDefaults()

	writer := &failingWriter{}

	err := p.ToMetaTags().Render(context.Background(), writer)
	if err == nil {
		t.Fatal("expected error but got nil")
	}

	expected := "failed to write og:type meta tag"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestPlace_ToGoHTMLMetaTags_Render(t *testing.T) {
	p := NewPlace(
		"Times Square",
		"https://example.com/ts",
		"Famous NYC location",
		"https://example.com/img.png",
		40.7580,
		-73.9855,
		"Broadway & 7th",
		"New York",
		"NY",
		"10036",
		"USA",
	)

	html, err := p.ToGoHTMLMetaTags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := string(html)
	if !strings.Contains(out, `property="og:title"`) {
		t.Errorf("expected 'og:title' in rendered HTML")
	}
	if !strings.Contains(out, "Times Square") {
		t.Errorf("expected place title in rendered HTML")
	}
	if !strings.Contains(out, `property="place:location:latitude"`) {
		t.Errorf("expected 'place:location:latitude' tag in rendered HTML")
	}
}
