package opengraph

import (
	"strings"
	"testing"
)

func TestNewArticle_Defaults(t *testing.T) {
	a := NewArticle("Title", "https://url", "Desc", "https://image.jpg",
		"2024-01-01T00:00:00Z", "2024-01-02T00:00:00Z", "2024-12-31T23:59:59Z",
		[]string{"https://author1"}, "Tech", []string{"tag1", "tag2"})

	if a.Type != "article" {
		t.Errorf("expected Type to be 'article', got %s", a.Type)
	}
	if a.Title != "Title" {
		t.Errorf("expected Title to be set")
	}
}

func TestArticleEnsureDefaults(t *testing.T) {
	art := &Article{}
	art.ensureDefaults()
	if art.Type != "article" {
		t.Errorf("ensureDefaults failed: expected 'article', got '%s'", art.Type)
	}
}

func TestArticleMetaTagsOutput(t *testing.T) {
	article := NewArticle(
		"Title",
		"https://example.com",
		"Desc",
		"https://example.com/image.jpg",
		"2024-09-15T09:00:00Z",
		"2024-09-15T10:00:00Z",
		"2024-12-31T23:59:59Z",
		[]string{"https://example.com/authors/jane-doe"},
		"Technology",
		[]string{"tag1", "tag2"},
	)

	tags := article.metaTags()

	assertTag := func(prop, val string) {
		found := false
		for _, tag := range tags {
			if tag.property == prop && tag.content == val {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("meta tag not found: %s=%s", prop, val)
		}
	}

	assertTag("og:type", "article")
	assertTag("og:title", "Title")
	assertTag("og:url", "https://example.com")
	assertTag("og:description", "Desc")
	assertTag("og:image", "https://example.com/image.jpg")
	assertTag("article:published_time", "2024-09-15T09:00:00Z")
	assertTag("article:modified_time", "2024-09-15T10:00:00Z")
	assertTag("article:expiration_time", "2024-12-31T23:59:59Z")
	assertTag("article:author", "https://example.com/authors/jane-doe")
	assertTag("article:section", "Technology")
	assertTag("article:tag", "tag1")
	assertTag("article:tag", "tag2")
}

func TestArticleMetaTags_SkipEmptyValues(t *testing.T) {
	article := &Article{
		OpenGraphObject: OpenGraphObject{
			Title: "Only Title",
		},
		Author: []string{"", "https://example.com/valid-author"},
		Tag:    []string{"", "valid-tag"},
	}

	tags := article.metaTags()

	// Should only contain non-empty authors and tags
	for _, tag := range tags {
		if tag.property == "article:author" && tag.content == "" {
			t.Errorf("empty author tag was included")
		}
		if tag.property == "article:tag" && tag.content == "" {
			t.Errorf("empty tag was included")
		}
	}
}

func TestArticleMetaTags_DuplicateKeysAllowed(t *testing.T) {
	article := &Article{
		OpenGraphObject: OpenGraphObject{
			Title: "Test",
		},
		Tag: []string{"a", "b"},
	}

	count := 0
	for _, tag := range article.metaTags() {
		if tag.property == "article:tag" {
			count++
		}
	}
	if count != 2 {
		t.Errorf("expected 2 article:tag meta tags, got %d", count)
	}
}

func TestArticle_ToGoHTMLMetaTags_Render(t *testing.T) {
	article := NewArticle(
		"Test Article",
		"https://example.com/article",
		"A description of the article.",
		"https://example.com/image.jpg",
		"2024-04-01T00:00:00Z",
		"2024-04-02T00:00:00Z",
		"2024-12-31T23:59:59Z",
		[]string{"https://example.com/authors/jane"},
		"Tech",
		[]string{"Go", "Testing"},
	)

	html, err := article.ToGoHTMLMetaTags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := string(html)

	if !strings.Contains(output, `property="og:title"`) || !strings.Contains(output, `content="Test Article"`) {
		t.Errorf("expected og:title tag missing or incorrect, got: %s", output)
	}

	if !strings.Contains(output, `property="article:published_time"`) || !strings.Contains(output, `content="2024-04-01T00:00:00Z"`) {
		t.Errorf("expected article:published_time tag missing or incorrect, got: %s", output)
	}

	if !strings.Contains(output, `property="article:tag"`) || !strings.Contains(output, `content="Go"`) {
		t.Errorf("expected article:tag tag missing or incorrect, got: %s", output)
	}
}
