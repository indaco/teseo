package opengraph

import (
	"context"
	"strings"
	"testing"
)

func TestNewProfile_SetsFieldsAndDefaults(t *testing.T) {
	p := NewProfile(
		"Jane Doe", "Jane", "Doe", "janedoe", "female",
		"https://example.com/jane", "Profile of Jane", "https://example.com/jane.jpg",
	)

	if p.Type != "profile" {
		t.Errorf("expected default type 'profile', got %s", p.Type)
	}
	if p.FirstName != "Jane" {
		t.Errorf("expected FirstName to be set")
	}
	if p.Gender != "female" {
		t.Errorf("expected Gender to be set")
	}
}

func TestProfile_ensureDefaults(t *testing.T) {
	p := &Profile{}
	p.ensureDefaults()
	if p.Type != "profile" {
		t.Errorf("expected type to default to 'profile'")
	}
}

func TestProfile_metaTags(t *testing.T) {
	p := NewProfile(
		"John Smith", "John", "Smith", "jsmith", "male",
		"https://example.com/jsmith", "Bio here", "https://example.com/img.jpg",
	)

	tags := p.metaTags()

	assertTag := func(prop, expected string) {
		for _, tag := range tags {
			if tag.property == prop && tag.content == expected {
				return
			}
		}
		t.Errorf("missing or incorrect tag: %s = %s", prop, expected)
	}

	assertTag("og:type", "profile")
	assertTag("profile:username", "jsmith")
	assertTag("profile:gender", "male")
}

func TestProfile_metaTags_SkipEmptyValues(t *testing.T) {
	p := &Profile{
		OpenGraphObject: OpenGraphObject{
			Title:       "John Doe",
			URL:         "https://example.com/profile/john",
			Description: "Example profile",
			Image:       "https://example.com/images/john.jpg",
		},
		FirstName: "John",
		LastName:  "Doe",
		Username:  "jd",
		Gender:    "",
	}

	tags := p.metaTags()
	for _, tag := range tags {
		if strings.HasPrefix(tag.property, "profile:") && tag.content == "" {
			t.Errorf("should skip tag with empty value: %s", tag.property)
		}
	}
}

func TestProfile_ToMetaTags_WriteError(t *testing.T) {
	p := NewProfile(
		"Dr. Profile", "Dr.", "Profile", "drprofile", "non-binary",
		"https://example.com/profile", "Profile description", "https://example.com/img.jpg",
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

func TestProfile_ToGoHTMLMetaTags_Render(t *testing.T) {
	p := NewProfile(
		"Dr. Profile", "Dr.", "Profile", "drprofile", "non-binary",
		"https://example.com/profile", "Profile description", "https://example.com/img.jpg",
	)

	html, err := p.ToGoHTMLMetaTags()
	if err != nil {
		t.Fatalf("unexpected error rendering HTML: %v", err)
	}

	out := string(html)
	if !strings.Contains(out, "profile:username") {
		t.Errorf("expected rendered HTML to contain profile:username")
	}
	if !strings.Contains(out, "non-binary") {
		t.Errorf("expected gender in rendered HTML")
	}
}
