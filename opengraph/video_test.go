package opengraph

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func TestVideo_metaTags(t *testing.T) {
	video := &Video{
		OpenGraphObject: OpenGraphObject{
			Title:       "Example Video",
			URL:         "https://example.com/video",
			Description: "A test video",
			Image:       "https://example.com/image.jpg",
		},
		Duration:    "300",
		ActorURLs:   []string{"https://example.com/actor1", "https://example.com/actor2"},
		DirectorURL: "https://example.com/director",
		ReleaseDate: "2024-09-15",
	}

	tags := video.metaTags()

	expected := map[string][]string{
		"og:type":            {"video.movie"},
		"og:title":           {"Example Video"},
		"og:url":             {"https://example.com/video"},
		"og:description":     {"A test video"},
		"og:image":           {"https://example.com/image.jpg"},
		"video:duration":     {"300"},
		"video:director":     {"https://example.com/director"},
		"video:release_date": {"2024-09-15"},
		"video:actor": {
			"https://example.com/actor1",
			"https://example.com/actor2",
		},
	}

	// Count occurrences
	actual := make(map[string][]string)
	for _, tag := range tags {
		actual[tag.property] = append(actual[tag.property], tag.content)
	}

	for key, expectedValues := range expected {
		actualValues, ok := actual[key]
		if !ok {
			t.Errorf("expected property %q missing", key)
			continue
		}
		if len(actualValues) != len(expectedValues) {
			t.Errorf("expected %d values for %q, got %d", len(expectedValues), key, len(actualValues))
			continue
		}
		for i, val := range expectedValues {
			if actualValues[i] != val {
				t.Errorf("expected %q for %q at index %d, got %q", val, key, i, actualValues[i])
			}
		}
	}
}

func TestVideo_metaTags_SkipEmptyValues(t *testing.T) {
	video := &Video{
		OpenGraphObject: OpenGraphObject{
			Title: "Video with Missing Fields",
		},
		ActorURLs: []string{"", ""},
	}

	tags := video.metaTags()
	for _, tag := range tags {
		if tag.property == "video:actor" {
			t.Errorf("expected empty actor URLs to be skipped, got: %q", tag.content)
		}
	}
}

func TestVideo_ToMetaTags_WriteError(t *testing.T) {
	video := NewVideo(
		"Example Video",
		"https://example.com/video",
		"Description here",
		"https://example.com/image.jpg",
		"300",
		[]string{"https://example.com/actor1"},
		"https://example.com/director",
		"2024-09-15",
	)
	video.ensureDefaults()

	writer := &failingWriter{}

	err := video.ToMetaTags().Render(context.Background(), writer)
	if err == nil {
		t.Fatal("expected error but got nil")
	}

	expected := "failed to write og:type meta tag"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestVideo_ToGoHTMLMetaTags_Render(t *testing.T) {
	video := NewVideo(
		"Example Video",
		"https://example.com/video",
		"Description here",
		"https://example.com/image.jpg",
		"300",
		[]string{"https://example.com/actor1"},
		"https://example.com/director",
		"2024-09-15",
	)

	html, err := video.ToGoHTMLMetaTags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	htmlStr := string(html)

	// Relaxed assertions: check for substrings without worrying about space before '>'
	if !strings.Contains(htmlStr, `property="og:title"`) || !strings.Contains(htmlStr, `content="Example Video"`) {
		t.Errorf("expected og:title meta tag in HTML, got: %s", htmlStr)
	}
	if !strings.Contains(htmlStr, `property="video:actor"`) || !strings.Contains(htmlStr, `content="https://example.com/actor1"`) {
		t.Errorf("expected video:actor meta tag in HTML, got: %s", htmlStr)
	}
}

func TestVideo_ToMetaTags_ImplementsTemplComponent(t *testing.T) {
	video := NewVideo("Test", "", "", "", "", nil, "", "")
	var _ templ.Component = video.ToMetaTags() // compile-time check
}
