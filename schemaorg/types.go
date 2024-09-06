package schemaorg

import (
	"context"
	"fmt"
	"html"
	"io"
	"strings"

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
		return templ.JSONScript(teseo.GenerateUniqueKey(), org).WithType("application/ld+json").Render(ctx, w)
	})
}

// ToGoHTMLJsonLd renders the Organization struct as a string for Go's `html/template`.
func (o *Organization) ToGoHTMLJsonLd() string {
	o.ensureDefaults()

	var sb strings.Builder
	sb.WriteString(`<script type="application/ld+json">`)
	sb.WriteString("\n{\n")
	sb.WriteString(fmt.Sprintf(`  "@context": "%s",`, html.EscapeString(o.Context)))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf(`  "@type": "%s",`, html.EscapeString(o.Type)))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf(`  "name": "%s",`, html.EscapeString(o.Name)))
	sb.WriteString("\n")
	if o.URL != "" {
		sb.WriteString(fmt.Sprintf(`  "url": "%s",`, html.EscapeString(o.URL)))
		sb.WriteString("\n")
	}
	if o.Logo != nil {
		sb.WriteString(`  "logo": {`)
		sb.WriteString(fmt.Sprintf(`"@type": "ImageObject", "url": "%s"`, html.EscapeString(o.Logo.URL)))
		sb.WriteString("},\n")
	}
	if len(o.ContactPoints) > 0 {
		sb.WriteString(`  "contactPoint": [`)
		for i, cp := range o.ContactPoints {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString("{\n")
			sb.WriteString(fmt.Sprintf(`    "@type": "ContactPoint", "telephone": "%s", "contactType": "%s"`, html.EscapeString(cp.Telephone), html.EscapeString(cp.ContactType)))
			if cp.AreaServed != "" {
				sb.WriteString(fmt.Sprintf(`, "areaServed": "%s"`, html.EscapeString(cp.AreaServed)))
			}
			if cp.AvailableLanguage != "" {
				sb.WriteString(fmt.Sprintf(`, "availableLanguage": "%s"`, html.EscapeString(cp.AvailableLanguage)))
			}
			sb.WriteString("\n}")
		}
		sb.WriteString("],\n")
	}
	if len(o.SameAs) > 0 {
		sb.WriteString(`  "sameAs": [`)
		for i, sameAs := range o.SameAs {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf(`"%s"`, html.EscapeString(sameAs)))
		}
		sb.WriteString("],\n")
	}
	sb.WriteString("}\n</script>")
	return sb.String()
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
