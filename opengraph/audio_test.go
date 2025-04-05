package opengraph

import (
	"context"
	"strings"
	"testing"
)

func TestNewAudio_SetsFieldsAndDefaults(t *testing.T) {
	audio := NewAudio(
		"Track Title",
		"https://example.com/audio",
		"Desc",
		"https://example.com/image.jpg",
		"180",
		"https://example.com/artist",
	)

	if audio.Type != "music.audio" {
		t.Errorf("expected Type to be music.audio, got %s", audio.Type)
	}
	if audio.Title != "Track Title" {
		t.Errorf("Title not set properly")
	}
	if audio.Duration != "180" {
		t.Errorf("Duration not set properly")
	}
}

func TestAudio_ensureDefaults(t *testing.T) {
	a := &Audio{}
	a.ensureDefaults()
	if a.Type != "music.audio" {
		t.Errorf("ensureDefaults did not set type to music.audio")
	}
}

func TestAudio_metaTags(t *testing.T) {
	audio := NewAudio(
		"Track Title",
		"https://example.com/audio",
		"Audio Description",
		"https://example.com/image.jpg",
		"300",
		"https://example.com/artist",
	)

	tags := audio.metaTags()

	assertMeta := func(prop, expected string) {
		for _, tag := range tags {
			if tag.property == prop && tag.content == expected {
				return
			}
		}
		t.Errorf("missing tag %s=%s", prop, expected)
	}

	assertMeta("og:type", "music.audio")
	assertMeta("og:title", "Track Title")
	assertMeta("og:url", "https://example.com/audio")
	assertMeta("og:description", "Audio Description")
	assertMeta("og:image", "https://example.com/image.jpg")
	assertMeta("music:duration", "300")
	assertMeta("music:musician", "https://example.com/artist")
}

func TestAudio_metaTags_SkipEmptyValues(t *testing.T) {
	audio := &Audio{
		OpenGraphObject: OpenGraphObject{
			Title: "A Title",
		},
		// Duration and ArtistURL are empty
	}

	tags := audio.metaTags()

	foundDuration := false
	foundMusician := false

	for _, tag := range tags {
		if tag.property == "music:duration" {
			foundDuration = true
			if tag.content != "" {
				t.Errorf("expected music:duration to be empty")
			}
		}
		if tag.property == "music:musician" {
			foundMusician = true
			if tag.content != "" {
				t.Errorf("expected music:musician to be empty")
			}
		}
	}

	if !foundDuration {
		t.Errorf("expected music:duration tag to be present")
	}
	if !foundMusician {
		t.Errorf("expected music:musician tag to be present")
	}
}

func TestAudio_ToMetaTags_WriteError(t *testing.T) {
	audio := NewAudio(
		"Track Title",
		"https://example.com/audio",
		"Desc",
		"https://example.com/image.jpg",
		"180",
		"https://example.com/artist",
	)
	audio.ensureDefaults()

	writer := &failingWriter{}

	err := audio.ToMetaTags().Render(context.Background(), writer)
	if err == nil {
		t.Fatal("expected error but got nil")
	}

	expected := "failed to write og:type meta tag"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestAudio_ToGoHTMLMetaTags(t *testing.T) {
	audio := NewAudio(
		"Track",
		"https://example.com/audio",
		"Desc",
		"https://example.com/img.jpg",
		"120",
		"https://example.com/artist",
	)

	html, err := audio.ToGoHTMLMetaTags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := string(html)
	if !strings.Contains(output, `twitter:title`) && !strings.Contains(output, `og:title`) {
		t.Errorf("expected some title meta tag in output, got: %s", output)
	}
}
