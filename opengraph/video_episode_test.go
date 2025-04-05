package opengraph

import (
	"context"
	"strings"
	"testing"
)

func TestVideoEpisode_metaTags(t *testing.T) {
	video := NewVideoEpisode(
		"Example Video Episode",
		"https://example.com/video/ep1",
		"An episode description.",
		"https://example.com/image.jpg",
		"1800",
		"https://example.com/series",
		[]string{"https://example.com/actor1", "https://example.com/actor2"},
		"https://example.com/director",
		"2024-09-15",
		1,
	)

	tags := video.metaTags()

	assertHas := func(property, content string) {
		t.Helper()
		for _, tag := range tags {
			if tag.property == property && tag.content == content {
				return
			}
		}
		t.Errorf("expected property %q with content %q not found", property, content)
	}

	assertHas("og:type", "video.episode")
	assertHas("og:title", "Example Video Episode")
	assertHas("og:url", "https://example.com/video/ep1")
	assertHas("og:description", "An episode description.")
	assertHas("og:image", "https://example.com/image.jpg")
	assertHas("video:duration", "1800")
	assertHas("video:series", "https://example.com/series")
	assertHas("video:actor", "https://example.com/actor1")
	assertHas("video:actor", "https://example.com/actor2")
	assertHas("video:director", "https://example.com/director")
	assertHas("video:release_date", "2024-09-15")
	assertHas("video:episode", "1")
}

func TestVideoEpisode_metaTags_SkipEmptyValues(t *testing.T) {
	video := NewVideoEpisode("Title", "", "", "", "", "", nil, "", "", 0)
	tags := video.metaTags()

	for _, tag := range tags {
		if tag.property == "video:director" && tag.content != "" {
			t.Errorf("expected empty video:director to be skipped")
		}
		if tag.property == "video:series" && tag.content != "" {
			t.Errorf("expected empty video:series to be skipped")
		}
		if tag.property == "video:episode" && tag.content != "" {
			t.Errorf("expected empty video:episode to be skipped")
		}
	}
}

func TestVideoEpisode_ToMetaTags_WriteError(t *testing.T) {
	video := NewVideoEpisode(
		"Episode Title",
		"https://example.com/ep",
		"A description.",
		"https://example.com/img.jpg",
		"1200",
		"https://example.com/series",
		[]string{"https://example.com/actor"},
		"https://example.com/director",
		"2024-01-01",
		3,
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

func TestVideoEpisode_ToGoHTMLMetaTags_Render(t *testing.T) {
	video := NewVideoEpisode(
		"Episode Title",
		"https://example.com/ep",
		"A description.",
		"https://example.com/img.jpg",
		"1200",
		"https://example.com/series",
		[]string{"https://example.com/actor"},
		"https://example.com/director",
		"2024-01-01",
		3,
	)

	html, err := video.ToGoHTMLMetaTags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := string(html)
	if !strings.Contains(out, `<meta property="og:title" content="Episode Title"`) {
		t.Errorf("expected title meta tag to be rendered")
	}
	if !strings.Contains(out, `<meta property="video:episode" content="3"`) {
		t.Errorf("expected episode number meta tag to be rendered")
	}
}
