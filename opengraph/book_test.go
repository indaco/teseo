package opengraph

import (
	"strings"
	"testing"
)

func TestNewBook_SetsFieldsAndDefaults(t *testing.T) {
	book := NewBook(
		"Book Title",
		"https://example.com/book",
		"Description",
		"https://example.com/img.jpg",
		"123-4567890123",
		"2024-09-15",
		[]string{"https://example.com/author/jane"},
		[]string{"fiction", "drama"},
	)

	if book.Type != "book" {
		t.Errorf("expected Type to be book, got %s", book.Type)
	}
	if book.ISBN != "123-4567890123" {
		t.Errorf("expected ISBN to be set")
	}
}

func TestBook_ensureDefaults(t *testing.T) {
	b := &Book{}
	b.ensureDefaults()
	if b.Type != "book" {
		t.Errorf("ensureDefaults did not set type to book")
	}
}

func TestBook_metaTags(t *testing.T) {
	book := NewBook(
		"Book Title",
		"https://example.com/book",
		"Desc",
		"https://example.com/img.jpg",
		"1234567890",
		"2024-01-01",
		[]string{"https://example.com/authors/one", "https://example.com/authors/two"},
		[]string{"tag1", "tag2"},
	)

	tags := book.metaTags()

	assertMeta := func(prop, expected string) {
		for _, tag := range tags {
			if tag.property == prop && tag.content == expected {
				return
			}
		}
		t.Errorf("missing tag: %s=%s", prop, expected)
	}

	assertMeta("og:type", "book")
	assertMeta("og:title", "Book Title")
	assertMeta("og:url", "https://example.com/book")
	assertMeta("og:description", "Desc")
	assertMeta("og:image", "https://example.com/img.jpg")
	assertMeta("book:isbn", "1234567890")
	assertMeta("book:release_date", "2024-01-01")
	assertMeta("book:author", "https://example.com/authors/one")
	assertMeta("book:author", "https://example.com/authors/two")
	assertMeta("book:tag", "tag1")
	assertMeta("book:tag", "tag2")
}

func TestBook_metaTags_EmptyValues(t *testing.T) {
	book := &Book{
		OpenGraphObject: OpenGraphObject{Title: "Book"},
		Author:          []string{"", "https://valid.com/author"},
		Tag:             []string{"", "valid-tag"},
	}

	tags := book.metaTags()

	authorFound := false
	tagFound := false

	for _, tag := range tags {
		if tag.property == "book:author" {
			if tag.content == "" {
				t.Errorf("empty author should be skipped")
			} else {
				authorFound = true
			}
		}
		if tag.property == "book:tag" {
			if tag.content == "" {
				t.Errorf("empty tag should be skipped")
			} else {
				tagFound = true
			}
		}
	}

	if !authorFound {
		t.Errorf("valid author not found")
	}
	if !tagFound {
		t.Errorf("valid tag not found")
	}
}

func TestBook_metaTags_Duplicates(t *testing.T) {
	book := &Book{
		OpenGraphObject: OpenGraphObject{Title: "Book"},
		Tag:             []string{"tag1", "tag2"},
	}

	count := 0
	for _, tag := range book.metaTags() {
		if tag.property == "book:tag" {
			count++
		}
	}
	if count != 2 {
		t.Errorf("expected 2 book:tag meta tags, got %d", count)
	}
}

func TestBook_ToGoHTMLMetaTags_Render(t *testing.T) {
	book := NewBook(
		"Test Book",
		"https://example.com/book",
		"A compelling description.",
		"https://example.com/cover.jpg",
		"9876543210",
		"2024-04-01",
		[]string{"https://example.com/author/john"},
		[]string{"fiction", "mystery"},
	)

	html, err := book.ToGoHTMLMetaTags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := string(html)

	if !strings.Contains(output, `property="og:title"`) || !strings.Contains(output, `content="Test Book"`) {
		t.Errorf("expected og:title tag missing or incorrect, got: %s", output)
	}

	if !strings.Contains(output, `property="book:author"`) || !strings.Contains(output, `content="https://example.com/author/john"`) {
		t.Errorf("expected book:author tag missing or incorrect, got: %s", output)
	}

	if !strings.Contains(output, `property="book:tag"`) || !strings.Contains(output, `content="fiction"`) {
		t.Errorf("expected book:tag tag missing or incorrect, got: %s", output)
	}
}
