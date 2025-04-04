package opengraph

import (
	"strings"
	"testing"
)

func TestNewProductGroup_SetsFieldsAndDefaults(t *testing.T) {
	pg := NewProductGroup(
		"Bundle",
		"https://example.com/group",
		"Bundle of items",
		"https://example.com/img.jpg",
		[]string{"https://example.com/p1", "https://example.com/p2"},
	)

	if pg.Type != "product.group" {
		t.Errorf("expected Type to be 'product.group', got '%s'", pg.Type)
	}
	if len(pg.Products) != 2 {
		t.Errorf("expected 2 products, got %d", len(pg.Products))
	}
}

func TestProductGroup_ensureDefaults(t *testing.T) {
	pg := &ProductGroup{}
	pg.ensureDefaults()
	if pg.Type != "product.group" {
		t.Errorf("expected default type to be 'product.group'")
	}
}

func TestProductGroup_metaTags(t *testing.T) {
	pg := NewProductGroup(
		"Tech Pack",
		"https://example.com/tech-pack",
		"A collection of tech products",
		"https://example.com/img.png",
		[]string{"https://example.com/macbook", "https://example.com/keyboard"},
	)

	tags := pg.metaTags()

	assertTag := func(prop, expected string) {
		for _, tag := range tags {
			if tag.property == prop && tag.content == expected {
				return
			}
		}
		t.Errorf("missing or incorrect tag: %s=%s", prop, expected)
	}

	assertTag("og:type", "product.group")
	assertTag("og:title", "Tech Pack")
	assertTag("og:url", "https://example.com/tech-pack")
	assertTag("og:description", "A collection of tech products")
	assertTag("og:image", "https://example.com/img.png")
	assertTag("product:group_item", "https://example.com/macbook")
	assertTag("product:group_item", "https://example.com/keyboard")
}

func TestProductGroup_metaTags_EmptyValuesFiltered(t *testing.T) {
	pg := &ProductGroup{
		OpenGraphObject: OpenGraphObject{Title: "Empty Product"},
		Products:        []string{"", "https://valid.com/product"},
	}

	tags := pg.metaTags()
	validFound := false

	for _, tag := range tags {
		if tag.property == "product:group_item" {
			if tag.content == "" {
				t.Errorf("empty product should not be rendered")
			} else {
				validFound = true
			}
		}
	}

	if !validFound {
		t.Errorf("valid product not rendered")
	}
}

func TestProductGroup_ToGoHTMLMetaTags_Render(t *testing.T) {
	pg := NewProductGroup(
		"Winter Gear",
		"https://example.com/winter",
		"Essential cold weather products",
		"https://example.com/img/winter.jpg",
		[]string{"https://example.com/coat"},
	)

	html, err := pg.ToGoHTMLMetaTags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := string(html)
	if !strings.Contains(out, `property="og:title"`) {
		t.Errorf("expected 'og:title' in rendered HTML")
	}
	if !strings.Contains(out, "Winter Gear") {
		t.Errorf("expected product group title in rendered HTML")
	}
	if !strings.Contains(out, `property="product:group_item"`) {
		t.Errorf("expected 'product:group_item' tag in rendered HTML")
	}
}
