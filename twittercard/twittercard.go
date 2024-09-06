package twittercard

import (
	"context"
	"html/template"
	"io"
	"log"

	"github.com/a-h/templ"
	"github.com/indaco/teseo"
)

// TwitterCardType represents the type of Twitter Card.
type TwitterCardType string

func (tc TwitterCardType) String() string {
	return string(tc)
}

const (
	CardSummary           TwitterCardType = "summary"
	CardSummaryLargeImage TwitterCardType = "summary_large_image"
	CardApp               TwitterCardType = "app"
	CardPlayer            TwitterCardType = "player"
)

// TwitterCard contains information for generating Twitter Card meta tags.
//
// Example usage:
//
// Pure struct usage:
//
//	twitterCard := &twittercard.TwitterCard{
//		Card:        twittercard.CardSummary,
//		Title:       "Example Title",
//		Description: "This is an example Twitter Card description.",
//		Image:       "https://www.example.com/image.jpg",
//		Site:        "@example_site",
//		Creator:     "@example_creator",
//	}
//
// Factory method usage:
//
//	twitterCard := twittercard.NewCard(
//		twittercard.CardSummary,
//		"Example Title",
//		"This is an example Twitter Card description.",
//		"https://www.example.com/image.jpg",
//		"@example_site",
//		"@example_creator",
//	)
//
//	// Generate the HTML meta tags
//	templ Page() {
//		@twitterCard.ToJsonLd()
//	})
//
// Expected output:
//
//	<meta name="twitter:card" content="summary"/>
//	<meta name="twitter:title" content="Example Title"/>
//	<meta name="twitter:description" content="This is an example Twitter Card description."/>
//	<meta name="twitter:image" content="https://www.example.com/image.jpg"/>
//	<meta name="twitter:site" content="@example_site"/>
//	<meta name="twitter:creator" content="@example_creator"/>
type TwitterCard struct {
	Card        TwitterCardType // Card type, e.g., "summary", "summary_large_image", "app", "player"
	Title       string          // Title of the content
	Description string          // Description of the content
	Image       string          // URL to a thumbnail image to be used in the card
	Site        string          // Twitter username of the website or the content creator
	Creator     string          // Twitter username of the content creator
	AppID       string          // App ID (used in app cards)
	PlayerURL   string          // URL of the player (used in player cards)
}

// NewCard initializes a TwitterCard based on the provided type.
func NewCard(cardType TwitterCardType, title string, description string, image string, site string, creator string) *TwitterCard {
	return &TwitterCard{
		Card:        cardType,
		Title:       title,
		Description: description,
		Image:       image,
		Site:        site,
		Creator:     creator,
	}
}

// SummaryCard represents a Twitter Card of type summary.
//
// Example usage:
//
// Pure struct usage:
//
//	summaryCard := &twittercard.TwitterCard{
//		Card:        twittercard.CardSummary,
//		Title:       "Example Summary",
//		Description: "This is an example summary card.",
//		Image:       "https://www.example.com/summary.jpg",
//		Site:        "@example_site",
//		Creator:     "@example_creator",
//	}
//
// Factory method usage:
//
//	summaryCard := twittercard.NewSummaryCard(
//		"Example Summary",
//		"This is an example summary card.",
//		"https://www.example.com/summary.jpg",
//		"@example_site",
//		"@example_creator",
//	)
//
//	// Generate the HTML meta tags
//	templ Page() {
//		@twitterCard.ToJsonLd()
//	})
//
// Expected output:
//
//	<meta name="twitter:card" content="summary"/>
//	<meta name="twitter:title" content="Example Summary"/>
//	<meta name="twitter:description" content="This is an example summary card."/>
//	<meta name="twitter:image" content="https://www.example.com/summary.jpg"/>
//	<meta name="twitter:site" content="@example_site"/>
//	<meta name="twitter:creator" content="@example_creator"/>
func NewSummaryCard(title string, description string, image string, site string, creator string) *TwitterCard {
	return NewCard(CardSummary, title, description, image, site, creator)
}

// SummaryLargeImageCard represents a Twitter Card of type summary_large_image.
//
// Example usage:
//
// Pure struct usage:
//
//	summaryLargeImageCard := &twittercard.TwitterCard{
//		Card:        twittercard.CardSummaryLargeImage,
//		Title:       "Example Summary Large Image",
//		Description: "This is an example large image summary card.",
//		Image:       "https://www.example.com/large_image.jpg",
//		Site:        "@example_site",
//		Creator:     "@example_creator",
//	}
//
// Factory method usage:
//
//	summaryLargeImageCard := twittercard.NewSummaryLargeImageCard(
//		"Example Summary Large Image",
//		"This is an example large image summary card.",
//		"https://www.example.com/large_image.jpg",
//		"@example_site",
//		"@example_creator",
//	)
//
//	// Generate the HTML meta tags
//	templ Page() {
//		@summaryLargeImageCard.ToJsonLd()
//	})
//
// Expected output:
//
//	<meta name="twitter:card" content="summary_large_image"/>
//	<meta name="twitter:title" content="Example Summary Large Image"/>
//	<meta name="twitter:description" content="This is an example large image summary card."/>
//	<meta name="twitter:image" content="https://www.example.com/large_image.jpg"/>
//	<meta name="twitter:site" content="@example_site"/>
//	<meta name="twitter:creator" content="@example_creator"/>
func NewSummaryLargeImageCard(title string, description string, image string, site string, creator string) *TwitterCard {
	return NewCard(CardSummaryLargeImage, title, description, image, site, creator)
}

// AppCard represents a Twitter App Card.
//
// Example usage:
//
// Pure struct usage:
//
//	appCard := &twittercard.TwitterCard{
//		Card:        twittercard.CardApp,
//		Title:       "Example App",
//		Description: "This is an example app card.",
//		Image:       "https://www.example.com/app.jpg",
//		Site:        "@example_site",
//		AppID:       "1234567890",
//	}
//
// Factory method usage:
//
//	appCard := twittercard.NewAppCard(
//		"Example App",
//		"This is an example app card.",
//		"https://www.example.com/app.jpg",
//		"@example_site",
//		"1234567890",
//	)
//
//	// Generate the HTML meta tags
//	templ Page() {
//		@appCard.ToJsonLd()
//	})
//
// Expected output:
//
//	<meta name="twitter:card" content="app"/>
//	<meta name="twitter:title" content="Example App"/>
//	<meta name="twitter:description" content="This is an example app card."/>
//	<meta name="twitter:image" content="https://www.example.com/app.jpg"/>
//	<meta name="twitter:site" content="@example_site"/>
//	<meta name="twitter:app:id:iphone" content="1234567890"/>
func NewAppCard(title string, description string, image string, site string, appID string) *TwitterCard {
	return &TwitterCard{
		Card:        CardApp,
		Title:       title,
		Description: description,
		Image:       image,
		Site:        site,
		AppID:       appID,
	}
}

// PlayerCard represents a Twitter Player Card.
//
// Example usage:
//
// Pure struct usage:
//
//	playerCard := &twittercard.TwitterCard{
//		Card:        twittercard.CardPlayer,
//		Title:       "Example Player",
//		Description: "This is an example player card.",
//		Image:       "https://www.example.com/player.jpg",
//		Site:        "@example_site",
//		PlayerURL:   "https://www.example.com/player",
//	}
//
// Factory method usage:
//
//	playerCard := twittercard.NewPlayerCard(
//		"Example Player",
//		"This is an example player card.",
//		"https://www.example.com/player.jpg",
//		"@example_site",
//		"https://www.example.com/player",
//	)
//
//	// Generate the HTML meta tags
//	templ Page() {
//		@playerCard.ToMetaTags()
//	})
//
// Expected output:
//
//	<meta name="twitter:card" content="player"/>
//	<meta name="twitter:title" content="Example Player"/>
//	<meta name="twitter:description" content="This is an example player card."/>
//	<meta name="twitter:image" content="https://www.example.com/player.jpg"/>
//	<meta name="twitter:site" content="@example_site"/>
//	<meta name="twitter:player" content="https://www.example.com/player"/>
func NewPlayerCard(title string, description string, image string, site string, playerURL string) *TwitterCard {
	return &TwitterCard{
		Card:        CardPlayer,
		Title:       title,
		Description: description,
		Image:       image,
		Site:        site,
		PlayerURL:   playerURL,
	}
}

// ToMetaTags generates the HTML meta tags for the Twitter Card using templ.Component
func (tc *TwitterCard) ToMetaTags() templ.Component {
	tc.ensureDefaults()
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		// Write each meta tag using the writeMetaTag helper
		for _, tag := range tc.metaTags() {
			if tag.content != "" {
				if err := teseo.WriteMetaTag(w, tag.name, tag.content); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// ToGoHTMLMetaTags generates the HTML meta tags for the Twitter Card as `template.HTML` value for Go's html/template
func (tc *TwitterCard) ToGoHTMLMetaTags() (template.HTML, error) {
	tc.ensureDefaults()

	// Create the templ component.
	templComponent := tc.ToMetaTags()

	// Render the templ component to a `template.HTML` value.
	html, err := templ.ToGoHTML(context.Background(), templComponent)
	if err != nil {
		log.Fatalf("failed to convert to html: %v", err)
	}

	return html, nil
}

// metaTags returns the meta tags for the Twitter Card as a slice of name-content pairs
func (tc *TwitterCard) metaTags() []struct {
	name    string
	content string
} {
	metaTags := []struct {
		name    string
		content string
	}{
		{"twitter:card", tc.Card.String()},
		{"twitter:title", tc.Title},
		{"twitter:description", tc.Description},
	}

	if tc.Image != "" {
		metaTags = append(metaTags, struct {
			name    string
			content string
		}{"twitter:image", tc.Image})
	}
	if tc.Site != "" {
		metaTags = append(metaTags, struct {
			name    string
			content string
		}{"twitter:site", tc.Site})
	}
	if tc.Creator != "" && (tc.Card == CardSummary || tc.Card == CardSummaryLargeImage) {
		metaTags = append(metaTags, struct {
			name    string
			content string
		}{"twitter:creator", tc.Creator})
	}
	if tc.AppID != "" && tc.Card == CardApp {
		metaTags = append(metaTags, struct {
			name    string
			content string
		}{"twitter:app:id:iphone", tc.AppID})
	}
	if tc.PlayerURL != "" && tc.Card == CardPlayer {
		metaTags = append(metaTags, struct {
			name    string
			content string
		}{"twitter:player", tc.PlayerURL})
	}

	return metaTags
}

func (tc *TwitterCard) ensureDefaults() {
	// Set default card type if not specified
	if tc.Card == "" {
		tc.Card = CardSummary
	}
}
