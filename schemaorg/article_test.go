package schemaorg

import (
	"testing"
)

func TestNewArticle_SetsFieldsAndDefaults(t *testing.T) {
	author := &Person{Name: "Jane"}
	publisher := &Organization{Name: "Example Publisher"}
	article := NewArticle("Headline", []string{"img1.jpg"}, author, publisher, "2024-01-01", "2024-01-02", "desc")

	if article.Type != "Article" {
		t.Errorf("expected type Article, got %s", article.Type)
	}
	if article.Context != "https://schema.org" {
		t.Errorf("expected schema.org context, got %s", article.Context)
	}
	if article.Headline != "Headline" {
		t.Errorf("headline not set properly")
	}
	if article.Image[0] != "img1.jpg" {
		t.Errorf("image not set properly")
	}
	if article.Author == nil || article.Author.Name != "Jane" {
		t.Errorf("author not set properly")
	}
	if article.Publisher == nil || article.Publisher.Name != "Example Publisher" {
		t.Errorf("publisher not set properly")
	}
}

func TestArticle_EnsureDefaults(t *testing.T) {
	art := &Article{}
	art.ensureDefaults()

	if art.Type != "Article" {
		t.Errorf("expected default type 'Article', got %s", art.Type)
	}
	if art.Context != "https://schema.org" {
		t.Errorf("expected default context 'https://schema.org', got %s", art.Context)
	}
}
func TestArticle_Validate_AllFieldsPresent(t *testing.T) {
	article := &Article{
		Headline:      "Example Headline",
		Image:         []string{"https://example.com/image.jpg"},
		DatePublished: "2024-09-15",
		Author:        &Person{Name: "Jane"},
		Publisher:     &Organization{Name: "Example Org"},
	}

	warnings := article.Validate()
	if len(warnings) != 0 {
		t.Errorf("expected no warnings, got %v", warnings)
	}
}

func TestArticle_Validate_MissingHeadline(t *testing.T) {
	article := &Article{
		Image:         []string{"https://example.com/image.jpg"},
		DatePublished: "2024-09-15",
		Author:        &Person{Name: "Jane"},
	}

	warnings := article.Validate()
	if len(warnings) != 1 || warnings[0] != "missing recommended field: headline" {
		t.Errorf("expected headline warning, got %v", warnings)
	}
}

func TestArticle_Validate_MissingImage(t *testing.T) {
	article := &Article{
		Headline:      "Title",
		DatePublished: "2024-09-15",
		Author:        &Person{Name: "Jane"},
	}

	warnings := article.Validate()
	if len(warnings) != 1 || warnings[0] != "missing recommended field: image" {
		t.Errorf("expected image warning, got %v", warnings)
	}
}

func TestArticle_Validate_MissingDatePublished(t *testing.T) {
	article := &Article{
		Headline: "Title",
		Image:    []string{"https://example.com/image.jpg"},
		Author:   &Person{Name: "Jane"},
	}

	warnings := article.Validate()
	if len(warnings) != 1 || warnings[0] != "missing recommended field: datePublished" {
		t.Errorf("expected datePublished warning, got %v", warnings)
	}
}

func TestArticle_Validate_MissingAuthorAndPublisher(t *testing.T) {
	article := &Article{
		Headline:      "Title",
		Image:         []string{"https://example.com/image.jpg"},
		DatePublished: "2024-09-15",
	}

	warnings := article.Validate()
	if len(warnings) != 1 || warnings[0] != "missing recommended field: author or publisher" {
		t.Errorf("expected author or publisher warning, got %v", warnings)
	}
}

func TestArticle_Validate_MultipleWarnings(t *testing.T) {
	article := &Article{}

	warnings := article.Validate()
	expected := map[string]bool{
		"missing recommended field: headline":            true,
		"missing recommended field: image":               true,
		"missing recommended field: datePublished":       true,
		"missing recommended field: author or publisher": true,
	}

	if len(warnings) != len(expected) {
		t.Errorf("expected %d warnings, got %d: %v", len(expected), len(warnings), warnings)
	}

	for _, w := range warnings {
		if !expected[w] {
			t.Errorf("unexpected warning: %s", w)
		}
	}
}

func TestArticle_ToGoHTMLJsonLd(t *testing.T) {
	article := NewArticle("Test", nil, nil, nil, "", "", "")
	html, err := article.ToGoHTMLJsonLd()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if html == "" {
		t.Errorf("expected non-empty html output")
	}
}
