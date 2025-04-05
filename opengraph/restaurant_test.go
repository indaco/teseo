package opengraph

import (
	"context"
	"strings"
	"testing"
)

func TestRestaurant_metaTags(t *testing.T) {
	r := NewRestaurant(
		"Example Restaurant",
		"https://www.example.com/restaurant/example-restaurant",
		"This is an example restaurant description.",
		"https://www.example.com/images/restaurant.jpg",
		"123 Food Street",
		"Gourmet City",
		"CA",
		"12345",
		"USA",
		"+1-800-FOOD-123",
		"https://www.example.com/menu",
		"https://www.example.com/reservations",
	)

	tags := r.metaTags()

	assertHas := func(key string) {
		found := false
		for _, tag := range tags {
			if tag.property == key {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected tag %q to be present", key)
		}
	}

	assertHas("og:type")
	assertHas("og:title")
	assertHas("og:url")
	assertHas("og:description")
	assertHas("og:image")
	assertHas("place:contact_data:street_address")
	assertHas("place:contact_data:locality")
	assertHas("place:contact_data:region")
	assertHas("place:contact_data:postal_code")
	assertHas("place:contact_data:country_name")
	assertHas("place:contact_data:phone_number")
	assertHas("restaurant:menu")
	assertHas("restaurant:reservation")
}

func TestRestaurant_ToMetaTags_WriteError(t *testing.T) {
	r := NewRestaurant(
		"Example Restaurant",
		"https://www.example.com/restaurant/example-restaurant",
		"This is an example restaurant description.",
		"https://www.example.com/images/restaurant.jpg",
		"123 Food Street",
		"Gourmet City",
		"CA",
		"12345",
		"USA",
		"+1-800-FOOD-123",
		"https://www.example.com/menu",
		"https://www.example.com/reservations",
	)
	r.ensureDefaults()

	writer := &failingWriter{}

	err := r.ToMetaTags().Render(context.Background(), writer)
	if err == nil {
		t.Fatal("expected error but got nil")
	}

	expected := "failed to write og:type meta tag"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRestaurant_ToGoHTMLMetaTags_Render(t *testing.T) {
	r := NewRestaurant(
		"Example Restaurant",
		"https://www.example.com/restaurant/example-restaurant",
		"This is an example restaurant description.",
		"https://www.example.com/images/restaurant.jpg",
		"123 Food Street",
		"Gourmet City",
		"CA",
		"12345",
		"USA",
		"+1-800-FOOD-123",
		"https://www.example.com/menu",
		"https://www.example.com/reservations",
	)

	html, err := r.ToGoHTMLMetaTags()
	if err != nil {
		t.Fatalf("failed to render meta tags: %v", err)
	}

	output := string(html)
	if !strings.Contains(output, `property="restaurant:menu"`) {
		t.Errorf("expected meta tag for restaurant:menu in output")
	}
	if !strings.Contains(output, `property="place:contact_data:phone_number"`) {
		t.Errorf("expected meta tag for phone number in output")
	}
}
