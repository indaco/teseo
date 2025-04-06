package twittercard

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestNewCard(t *testing.T) {
	card := NewCard(CardSummary, "Title", "Desc", "https://image.jpg", "@site", "@creator")

	if card.Card != CardSummary {
		t.Errorf("expected CardSummary, got %s", card.Card)
	}
	if card.Title != "Title" || card.Description != "Desc" {
		t.Errorf("unexpected title or description")
	}
}

func TestNewSummaryCard(t *testing.T) {
	card := NewSummaryCard("Title", "Desc", "https://img.jpg", "@site", "@creator")
	if card.Card != CardSummary {
		t.Errorf("expected summary card, got %s", card.Card)
	}
}

func TestNewSummaryLargeImageCard(t *testing.T) {
	card := NewSummaryLargeImageCard("Title", "Desc", "https://img.jpg", "@site", "@creator")
	if card.Card != CardSummaryLargeImage {
		t.Errorf("expected summary_large_image, got %s", card.Card)
	}
}

func TestNewAppCard(t *testing.T) {
	card := NewAppCard("App Title", "App Desc", "https://img.jpg", "@site", "12345")
	if card.Card != CardApp || card.AppID != "12345" {
		t.Errorf("AppCard not set up properly")
	}
}

func TestNewPlayerCard(t *testing.T) {
	card := NewPlayerCard("Title", "Desc", "https://img.jpg", "@site", "https://player.url")
	if card.Card != CardPlayer || card.PlayerURL != "https://player.url" {
		t.Errorf("PlayerCard not set up properly")
	}
}

func TestMetaTagsContent(t *testing.T) {
	card := NewSummaryCard("Title", "Desc", "https://img.jpg", "@site", "@creator")
	tags := card.metaTags()

	expected := map[string]string{
		"twitter:card":        "summary",
		"twitter:title":       "Title",
		"twitter:description": "Desc",
		"twitter:image":       "https://img.jpg",
		"twitter:site":        "@site",
		"twitter:creator":     "@creator",
	}

	for _, tag := range tags {
		if val, ok := expected[tag.name]; !ok || val != tag.content {
			t.Errorf("unexpected tag: %s=%s", tag.name, tag.content)
		}
	}
}

func TestToMetaTags_Render(t *testing.T) {
	card := NewSummaryCard("Title", "Desc", "https://img.jpg", "@site", "@creator")
	var buf bytes.Buffer
	err := card.ToMetaTags().Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("rendering failed: %v", err)
	}

	output := buf.String()
	required := []string{
		`<meta property="twitter:card" content="summary" >`,
		`<meta property="twitter:title" content="Title" >`,
		`<meta property="twitter:description" content="Desc" >`,
		`<meta property="twitter:image" content="https://img.jpg" >`,
		`<meta property="twitter:site" content="@site" >`,
		`<meta property="twitter:creator" content="@creator" >`,
	}

	for _, line := range required {
		if !strings.Contains(output, line) {
			t.Errorf("expected meta tag not found: %s", line)
		}
	}
}

func TestMetaTags_AppCardIncludesAppID(t *testing.T) {
	card := &TwitterCard{
		Card:        CardApp,
		Title:       "App Title",
		Description: "App Description",
		AppID:       "9876543210",
	}

	tags := card.metaTags()

	found := false
	for _, tag := range tags {
		if tag.name == "twitter:app:id:iphone" && tag.content == "9876543210" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected twitter:app:id:iphone tag to be included with correct value")
	}
}

func TestEnsureDefaults_SetsDefaultCardType(t *testing.T) {
	card := &TwitterCard{} // Card is empty
	card.ensureDefaults()

	if card.Card != CardSummary {
		t.Errorf("expected default card type to be CardSummary, got %s", card.Card)
	}
}

func TestToGoHTMLMetaTags_Output(t *testing.T) {
	card := NewPlayerCard("Title", "Desc", "https://img.jpg", "@site", "https://player.url")
	html, err := card.ToGoHTMLMetaTags()
	if err != nil {
		t.Fatalf("ToGoHTMLMetaTags failed: %v", err)
	}

	if !strings.Contains(string(html), "twitter:player") {
		t.Errorf("expected twitter:player in HTML output")
	}
}

func TestToMetaTags_WriteMetaTagError(t *testing.T) {
	// Backup and replace WriteMetaTag
	original := WriteMetaTag
	defer func() { WriteMetaTag = original }()

	WriteMetaTag = func(w io.Writer, property, content string) error {
		return fmt.Errorf("simulated failure on property: %s", property)
	}

	card := NewSummaryCard("Title", "Desc", "https://img.jpg", "@site", "@creator")

	var buf bytes.Buffer
	err := card.ToMetaTags().Render(context.Background(), &buf)
	if err == nil || !strings.Contains(err.Error(), "simulated failure") {
		t.Errorf("expected simulated failure, got: %v", err)
	}
}
