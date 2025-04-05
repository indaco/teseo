package opengraph

import (
	"context"
	"strings"
	"testing"
)

func TestNewProduct_SetsFieldsAndDefaults(t *testing.T) {
	p := NewProduct(
		"Eco Bottle",
		"https://example.com/bottle",
		"Reusable water bottle",
		"https://example.com/img.jpg",
		"19.95",
		"USD",
	)

	if p.Type != "product" {
		t.Errorf("expected type to be 'product', got '%s'", p.Type)
	}
	if p.Price != "19.95" {
		t.Errorf("expected price to be set, got '%s'", p.Price)
	}
	if p.PriceCurrency != "USD" {
		t.Errorf("expected currency to be set, got '%s'", p.PriceCurrency)
	}
}

func TestProduct_ensureDefaults(t *testing.T) {
	p := &Product{}
	p.ensureDefaults()
	if p.Type != "product" {
		t.Errorf("expected default type to be 'product'")
	}
}

func TestProduct_metaTags(t *testing.T) {
	p := NewProduct(
		"Coffee Mug",
		"https://example.com/mug",
		"Stylish ceramic mug",
		"https://example.com/mug.jpg",
		"12.00",
		"EUR",
	)

	tags := p.metaTags()

	assertTag := func(prop, expected string) {
		for _, tag := range tags {
			if tag.property == prop && tag.content == expected {
				return
			}
		}
		t.Errorf("missing or incorrect tag: %s=%s", prop, expected)
	}

	assertTag("og:type", "product")
	assertTag("og:title", "Coffee Mug")
	assertTag("og:url", "https://example.com/mug")
	assertTag("og:description", "Stylish ceramic mug")
	assertTag("og:image", "https://example.com/mug.jpg")
	assertTag("product:price:amount", "12.00")
	assertTag("product:price:currency", "EUR")
}

func TestProduct_metaTags_SkipEmptyValues(t *testing.T) {
	p := &Product{
		OpenGraphObject: OpenGraphObject{Title: "Free Item"},
		Price:           "",
		PriceCurrency:   "",
	}

	tags := p.metaTags()
	for _, tag := range tags {
		if tag.property == "product:price:amount" && tag.content == "" {
			t.Errorf("empty price should not be rendered")
		}
		if tag.property == "product:price:currency" && tag.content == "" {
			t.Errorf("empty currency should not be rendered")
		}
	}
}

func TestProduct_ToMetaTags_WriteError(t *testing.T) {
	p := NewProduct(
		"Noise-Cancelling Headphones",
		"https://example.com/headphones",
		"Wireless over-ear headphones with noise cancellation",
		"https://example.com/img/headphones.jpg",
		"199.99",
		"USD",
	)
	p.ensureDefaults()

	writer := &failingWriter{}

	err := p.ToMetaTags().Render(context.Background(), writer)
	if err == nil {
		t.Fatal("expected error but got nil")
	}

	expected := "failed to write og:type meta tag"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestProduct_ToGoHTMLMetaTags_Render(t *testing.T) {
	p := NewProduct(
		"Noise-Cancelling Headphones",
		"https://example.com/headphones",
		"Wireless over-ear headphones with noise cancellation",
		"https://example.com/img/headphones.jpg",
		"199.99",
		"USD",
	)

	html, err := p.ToGoHTMLMetaTags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := string(html)
	if !strings.Contains(out, `property="og:title"`) {
		t.Errorf("expected 'og:title' tag")
	}
	if !strings.Contains(out, `Noise-Cancelling Headphones`) {
		t.Errorf("expected title content")
	}
	if !strings.Contains(out, `property="product:price:amount"`) {
		t.Errorf("expected 'product:price:amount' tag")
	}
	if !strings.Contains(out, `199.99`) {
		t.Errorf("expected price content")
	}
}
