package opengraph

import (
	"strings"
	"testing"
)

func TestNewMusicRadioStation_SetsFieldsAndDefaults(t *testing.T) {
	mrs := NewMusicRadioStation(
		"Example Radio",
		"https://example.com/radio",
		"A streaming station",
		"https://example.com/image.jpg",
	)

	if mrs.Type != "music.radio_station" {
		t.Errorf("expected Type to be 'music.radio_station', got '%s'", mrs.Type)
	}
	if mrs.Title != "Example Radio" {
		t.Errorf("expected Title to be set")
	}
}

func TestMusicRadioStation_ensureDefaults(t *testing.T) {
	mrs := &MusicRadioStation{}
	mrs.ensureDefaults()
	if mrs.Type != "music.radio_station" {
		t.Errorf("expected default type to be 'music.radio_station'")
	}
}

func TestMusicRadioStation_metaTags(t *testing.T) {
	mrs := NewMusicRadioStation(
		"Radio One",
		"https://radio.com",
		"Your favorite hits",
		"https://radio.com/logo.jpg",
	)

	tags := mrs.metaTags()

	assertTag := func(prop, expected string) {
		for _, tag := range tags {
			if tag.property == prop && tag.content == expected {
				return
			}
		}
		t.Errorf("missing or incorrect tag: %s=%s", prop, expected)
	}

	assertTag("og:type", "music.radio_station")
	assertTag("og:title", "Radio One")
	assertTag("og:url", "https://radio.com")
	assertTag("og:description", "Your favorite hits")
	assertTag("og:image", "https://radio.com/logo.jpg")
}

func TestMusicRadioStation_metaTags_EmptyValues(t *testing.T) {
	mrs := &MusicRadioStation{
		OpenGraphObject: OpenGraphObject{Title: "Radio"},
	}

	tags := mrs.metaTags()

	foundType := false
	foundTitle := false
	for _, tag := range tags {
		if tag.property == "og:type" {
			foundType = true
		}
		if tag.property == "og:title" && tag.content == "Radio" {
			foundTitle = true
		}
	}
	if !foundType {
		t.Errorf("expected og:type to be present")
	}
	if !foundTitle {
		t.Errorf("expected og:title to be 'Radio'")
	}
}

func TestMusicRadioStation_ToGoHTMLMetaTags_Render(t *testing.T) {
	mrs := NewMusicRadioStation(
		"Morning Beats",
		"https://radio.com/morning",
		"Start your day right",
		"https://radio.com/morning.jpg",
	)

	html, err := mrs.ToGoHTMLMetaTags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := string(html)
	if !strings.Contains(out, `property="og:title"`) {
		t.Errorf("expected 'og:title' in rendered HTML")
	}
	if !strings.Contains(out, "Morning Beats") {
		t.Errorf("expected radio station title in rendered HTML")
	}
	if !strings.Contains(out, `property="og:type"`) {
		t.Errorf("expected 'og:type' in rendered HTML")
	}
}
