package opengraph

import (
	"strings"
	"testing"
)

func TestNewBusiness_SetsFieldsAndDefaults(t *testing.T) {
	b := NewBusiness(
		"Example Business",
		"https://example.com/business",
		"Description",
		"https://example.com/image.jpg",
		"123 Main St",
		"Anytown",
		"CA",
		"90210",
		"USA",
		"info@example.com",
		"+1-555-555-5555",
		"https://example.com",
	)

	if b.Type != "business.business" {
		t.Errorf("expected type 'business.business', got '%s'", b.Type)
	}
	if b.StreetAddress != "123 Main St" || b.Email != "info@example.com" {
		t.Errorf("fields not set correctly")
	}
}

func TestBusiness_ensureDefaults(t *testing.T) {
	b := &Business{}
	b.ensureDefaults()
	if b.Type != "business.business" {
		t.Errorf("expected default type to be 'business.business', got '%s'", b.Type)
	}
}

func TestBusiness_metaTags(t *testing.T) {
	b := NewBusiness(
		"Biz",
		"https://biz.com",
		"Desc",
		"https://img.com/biz.jpg",
		"100 Market St",
		"Big City",
		"NY",
		"10001",
		"USA",
		"hello@biz.com",
		"123-456-7890",
		"https://biz.com",
	)

	tags := b.metaTags()

	assertTag := func(prop, val string) {
		for _, tag := range tags {
			if tag.property == prop && tag.content == val {
				return
			}
		}
		t.Errorf("missing or incorrect tag: %s=%s", prop, val)
	}

	assertTag("og:type", "business.business")
	assertTag("og:title", "Biz")
	assertTag("og:url", "https://biz.com")
	assertTag("og:description", "Desc")
	assertTag("og:image", "https://img.com/biz.jpg")
	assertTag("business:contact_data:street_address", "100 Market St")
	assertTag("business:contact_data:locality", "Big City")
	assertTag("business:contact_data:region", "NY")
	assertTag("business:contact_data:postal_code", "10001")
	assertTag("business:contact_data:country_name", "USA")
	assertTag("business:contact_data:email", "hello@biz.com")
	assertTag("business:contact_data:phone_number", "123-456-7890")
	assertTag("business:contact_data:website", "https://biz.com")
}

func TestBusiness_metaTags_EmptyValues(t *testing.T) {
	b := &Business{
		OpenGraphObject: OpenGraphObject{Title: "Test"},
		Email:           "",
		Website:         "https://test.com",
	}

	tags := b.metaTags()

	emailSeen := false
	websiteSeen := false

	for _, tag := range tags {
		if tag.property == "business:contact_data:email" {
			if tag.content != "" {
				t.Errorf("expected empty email tag to still be empty (rendering skips it, not metaTags)")
			}
			emailSeen = true
		}
		if tag.property == "business:contact_data:website" && tag.content == "https://test.com" {
			websiteSeen = true
		}
	}

	if !emailSeen {
		t.Errorf("expected email tag (even if empty) to be present")
	}
	if !websiteSeen {
		t.Errorf("website tag was missing")
	}
}

func TestBusiness_ToGoHTMLMetaTags_Render(t *testing.T) {
	b := NewBusiness(
		"Example",
		"https://example.com",
		"Desc",
		"https://img.com/example.jpg",
		"123 St",
		"City",
		"Region",
		"00000",
		"Country",
		"email@example.com",
		"000-0000",
		"https://example.com",
	)

	html, err := b.ToGoHTMLMetaTags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(string(html), `property="og:title"`) {
		t.Errorf("expected meta tags to include og:title")
	}
}
