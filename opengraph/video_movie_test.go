package opengraph

import (
	"context"
	"strings"
	"testing"
)

func TestVideoMovie_metaTags_GeneratesCorrectTags(t *testing.T) {
	video := NewVideoMovie(
		"Example Movie",
		"https://www.example.com/video/movie/example-movie",
		"This is an example movie description.",
		"https://www.example.com/images/movie.jpg",
		"7200",
		[]string{
			"https://www.example.com/actors/jane-doe",
			"https://www.example.com/actors/john-doe",
		},
		"https://www.example.com/directors/jane-director",
		"2024-09-15",
	)

	expected := []metaTag{
		{"og:type", "video.movie"},
		{"og:title", "Example Movie"},
		{"og:url", "https://www.example.com/video/movie/example-movie"},
		{"og:description", "This is an example movie description."},
		{"og:image", "https://www.example.com/images/movie.jpg"},
		{"video:duration", "7200"},
		{"video:director", "https://www.example.com/directors/jane-director"},
		{"video:release_date", "2024-09-15"},
		{"video:actor", "https://www.example.com/actors/jane-doe"},
		{"video:actor", "https://www.example.com/actors/john-doe"},
	}

	actual := video.metaTags()

	if len(expected) != len(actual) {
		t.Fatalf("expected %d tags, got %d", len(expected), len(actual))
	}

	for i, tag := range expected {
		if tag != actual[i] {
			t.Errorf("expected tag %d to be %+v, got %+v", i, tag, actual[i])
		}
	}
}

func TestVideoMovie_ToMetaTags_WriteError(t *testing.T) {
	video := NewVideoMovie(
		"Example Movie",
		"https://www.example.com/video/movie/example-movie",
		"This is an example movie description.",
		"https://www.example.com/images/movie.jpg",
		"7200",
		[]string{
			"https://www.example.com/actors/jane-doe",
			"https://www.example.com/actors/john-doe",
		},
		"https://www.example.com/directors/jane-director",
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

func TestVideoMovie_ToGoHTMLMetaTags_Render(t *testing.T) {
	video := NewVideoMovie(
		"Example Movie",
		"https://www.example.com/video/movie/example-movie",
		"This is an example movie description.",
		"https://www.example.com/images/movie.jpg",
		"7200",
		[]string{
			"https://www.example.com/actors/jane-doe",
			"https://www.example.com/actors/john-doe",
		},
		"https://www.example.com/directors/jane-director",
		"2024-09-15",
	)

	html, err := video.ToGoHTMLMetaTags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content := string(html)
	assertContains := func(sub string) {
		if !strings.Contains(content, sub) {
			t.Errorf("expected HTML to contain %q", sub)
		}
	}

	assertContains(`<meta property="og:type" content="video.movie"`)
	assertContains(`<meta property="video:actor" content="https://www.example.com/actors/john-doe"`)
}
