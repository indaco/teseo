package schemaorg

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"

	"github.com/a-h/templ"
	"github.com/indaco/teseo"
)

// Common type definitions used across multiple JSON-LD entities

// ContactPoint represents a Schema.org ContactPoint object
type ContactPoint struct {
	Type              string `json:"@type"`
	Telephone         string `json:"telephone,omitempty"`
	ContactType       string `json:"contactType,omitempty"`
	AreaServed        string `json:"areaServed,omitempty"`
	AvailableLanguage string `json:"availableLanguage,omitempty"`
}

// ImageObject represents a Schema.org ImageObject object
type ImageObject struct {
	Type string `json:"@type"`
	URL  string `json:"url,omitempty"`
}

// ensureDefaults sets default values for ImageObject if they are not already set.
func (img *ImageObject) ensureDefaults() {
	if img.Type == "" {
		img.Type = "ImageObject"
	}
}

// Organization represents a Schema.org Organization object
type Organization struct {
	Context       string         `json:"@context"`
	Type          string         `json:"@type"`
	Name          string         `json:"name,omitempty"`
	URL           string         `json:"url,omitempty"`
	Logo          *ImageObject   `json:"logo,omitempty"`
	ContactPoints []ContactPoint `json:"contactPoint,omitempty"`
	SameAs        []string       `json:"sameAs,omitempty"`
}

func (org *Organization) ensureDefaults() {
	if org.Context == "" {
		org.Context = "https://schema.org"
	}

	if org.Type == "" {
		org.Type = "Organization"
	}

	if org.Logo != nil {
		org.Logo.ensureDefaults()
	}
}

// ToJsonLd converts the Organization struct to a JSON-LD `templ.Component`.
func (org *Organization) ToJsonLd() templ.Component {
	org.ensureDefaults()
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		id := fmt.Sprintf("%s-%s", "org", teseo.GenerateUniqueKey())
		return templ.JSONScript(id, org).WithType("application/ld+json").Render(ctx, w)
	})
}

// ToGoHTMLJsonLd renders the Organization struct as `template.HTML` value for Go's `html/template`.
func (org *Organization) ToGoHTMLJsonLd() (template.HTML, error) {
	// Create the templ component.
	templComponent := org.ToJsonLd()

	// Render the templ component to a `template.HTML` value.
	html, err := templ.ToGoHTML(context.Background(), templComponent)
	if err != nil {
		log.Fatalf("failed to convert to html: %v", err)
	}

	return html, nil
}

// Person represents a Schema.org Person object
type Person struct {
	Context     string         `json:"@context"`
	Type        string         `json:"@type"`
	Name        string         `json:"name,omitempty"`
	URL         string         `json:"url,omitempty"`
	Email       string         `json:"email,omitempty"`
	Image       *ImageObject   `json:"image,omitempty"`
	JobTitle    string         `json:"jobTitle,omitempty"`
	WorksFor    *Organization  `json:"worksFor,omitempty"`
	SameAs      []string       `json:"sameAs,omitempty"`
	Gender      string         `json:"gender,omitempty"`
	BirthDate   string         `json:"birthDate,omitempty"`
	Nationality string         `json:"nationality,omitempty"`
	Telephone   string         `json:"telephone,omitempty"`
	Address     *PostalAddress `json:"address,omitempty"`
	Affiliation *Organization  `json:"affiliation,omitempty"`
}

// ListItem represents a Schema.org ListItem object
type ListItem struct {
	Type     string `json:"@type"`
	Position int    `json:"position,omitempty"`
	Name     string `json:"name,omitempty"`
	Item     string `json:"item,omitempty"`
}
