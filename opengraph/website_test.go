package opengraph

import (
	"context"
	"strings"
	"testing"
)

func TestWebSite_metaTags(t *testing.T) {
	ws := &WebSite{
		OpenGraphObject: OpenGraphObject{
			Title:       "Example Website",
			URL:         "https://www.example.com",
			Description: "An example website",
			Image:       "https://www.example.com/image.jpg",
		},
	}
	expected := []metaTag{
		{"og:type", "website"},
		{"og:title", "Example Website"},
		{"og:url", "https://www.example.com"},
		{"og:description", "An example website"},
		{"og:image", "https://www.example.com/image.jpg"},
	}

	actual := ws.metaTags()
	if len(actual) != len(expected) {
		t.Fatalf("expected %d meta tags, got %d", len(expected), len(actual))
	}

	for i, tag := range expected {
		if actual[i] != tag {
			t.Errorf("expected tag %v at index %d, got %v", tag, i, actual[i])
		}
	}
}

func TestNewWebSite_SetsType(t *testing.T) {
	ws := NewWebSite(
		"My Website",
		"https://example.com",
		"Just a test",
		"https://example.com/logo.jpg",
	)

	if ws.OpenGraphObject.Type != "website" {
		t.Errorf("expected type to be 'website', got '%s'", ws.OpenGraphObject.Type)
	}
}

func TestWebSite_ToMetaTags_Render(t *testing.T) {
	ws := NewWebSite(
		"My Website",
		"https://example.com",
		"A test site",
		"https://example.com/logo.jpg",
	)

	var sb strings.Builder
	err := ws.ToMetaTags().Render(context.Background(), &sb)
	if err != nil {
		t.Fatalf("unexpected error rendering meta tags: %v", err)
	}

	out := sb.String()

	if !strings.Contains(out, `property="og:type"`) || !strings.Contains(out, `content="website"`) {
		t.Error("missing or incorrect og:type tag")
	}
	if !strings.Contains(out, `property="og:title"`) || !strings.Contains(out, `content="My Website"`) {
		t.Error("missing or incorrect og:title tag")
	}
	if !strings.Contains(out, `property="og:url"`) || !strings.Contains(out, `content="https://example.com"`) {
		t.Error("missing or incorrect og:url tag")
	}
	if !strings.Contains(out, `property="og:description"`) || !strings.Contains(out, `content="A test site"`) {
		t.Error("missing or incorrect og:description tag")
	}
	if !strings.Contains(out, `property="og:image"`) || !strings.Contains(out, `content="https://example.com/logo.jpg"`) {
		t.Error("missing or incorrect og:image tag")
	}
}

func TestWebSite_ToGoHTMLMetaTags_Render(t *testing.T) {
	ws := NewWebSite(
		"Example Website",
		"https://example.com",
		"Description here",
		"https://example.com/image.jpg",
	)

	html, err := ws.ToGoHTMLMetaTags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	htmlStr := string(html)

	if !strings.Contains(htmlStr, `property="og:title"`) {
		t.Errorf("expected og:title meta tag in HTML, got: %s", htmlStr)
	}
	if !strings.Contains(htmlStr, `property="og:type"`) {
		t.Errorf("expected og:type meta tag in HTML, got: %s", htmlStr)
	}
}
