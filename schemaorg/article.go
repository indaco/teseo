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

// Article represents a Schema.org Article object.
//
// Example usage:
//
// Pure struct usage:
//
//	article := &schemaorg.Article{
//		Headline:      "Example Article Headline",
//		Image:         []string{"https://www.example.com/images/article.jpg"},
//		Author:        &schemaorg.Person{Name: "Jane Doe"},
//		Publisher:     &schemaorg.Organization{Name: "Example Publisher"},
//		DatePublished: "2024-09-15",
//		DateModified:  "2024-09-16",
//		Description:   "This is an example article.",
//	}
//
// Factory method usage:
//
//	article := schemaorg.NewArticle(
//		"Example Article Headline",
//		[]string{"https://www.example.com/images/article.jpg"},
//		&schemaorg.Person{Name: "Jane Doe"},
//		&schemaorg.Organization{Name: "Example Publisher"},
//		"2024-09-15",
//		"2024-09-16",
//		"This is an example article",
//	)
//
// // Rendering JSON-LD using templ:
//
//	templ Page() {
//		@article.ToJsonLd()
//	}
//
// // Rendering JSON-LD as `template.HTML` value:
//
//	jsonLdHtml := article.ToGoHTMLJsonLd()
//
// Expected output:
//
//	{
//		"@context": "https://schema.org",
//		"@type": "Article",
//		"headline": "Example Article Headline",
//		"image": ["https://www.example.com/images/article.jpg"],
//		"author": {"@type": "Person", "name": "Jane Doe"},
//		"publisher": {"@type": "Organization", "name": "Example Publisher"},
//		"datePublished": "2024-09-15",
//		"dateModified": "2024-09-16",
//		"description": "This is an example article"
//	}
type Article struct {
	Context       string        `json:"@context"`
	Type          string        `json:"@type"`
	Headline      string        `json:"headline,omitempty"`
	Image         []string      `json:"image,omitempty"`
	Author        *Person       `json:"author,omitempty"`
	Publisher     *Organization `json:"publisher,omitempty"`
	DatePublished string        `json:"datePublished,omitempty"`
	DateModified  string        `json:"dateModified,omitempty"`
	Description   string        `json:"description,omitempty"`
}

// NewArticle initializes an Article with default context and type.
func NewArticle(headline string, images []string, author *Person, publisher *Organization, datePublished, dateModified, description string) *Article {
	article := &Article{
		Headline:      headline,
		Image:         images,
		Author:        author,
		Publisher:     publisher,
		DatePublished: datePublished,
		DateModified:  dateModified,
		Description:   description,
	}
	article.ensureDefaults()
	return article
}

// ToJsonLd converts the Article struct to a JSON-LD `templ.Component`.
func (art *Article) ToJsonLd() templ.Component {
	art.ensureDefaults()
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		return templ.JSONScript(teseo.GenerateUniqueKey(), art).WithType("application/ld+json").Render(ctx, w)
	})
}

// ToGoHTMLJsonLd renders the Article struct as a string for Go's `html/template`.
func (art *Article) ToGoHTMLJsonLd() string {
	art.ensureDefaults()

	var sb strings.Builder
	sb.WriteString(`<script type="application/ld+json">`)
	sb.WriteString("\n{\n")
	sb.WriteString(fmt.Sprintf(`  "@context": "%s",`, html.EscapeString(art.Context)))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf(`  "@type": "%s",`, html.EscapeString(art.Type)))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf(`  "headline": "%s",`, html.EscapeString(art.Headline)))
	sb.WriteString("\n")

	// Images
	if len(art.Image) > 0 {
		sb.WriteString(`  "image": [`)
		for i, img := range art.Image {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf(`"%s"`, html.EscapeString(img)))
		}
		sb.WriteString("],\n")
	}

	// Author
	if art.Author != nil {
		sb.WriteString(`  "author": {`)
		sb.WriteString(fmt.Sprintf(`"@type": "Person", "name": "%s"`, html.EscapeString(art.Author.Name)))
		if art.Author.URL != "" {
			sb.WriteString(fmt.Sprintf(`, "url": "%s"`, html.EscapeString(art.Author.URL)))
		}
		sb.WriteString("},\n")
	}

	// Publisher
	if art.Publisher != nil {
		sb.WriteString(`  "publisher": {`)
		sb.WriteString(fmt.Sprintf(`"@type": "Organization", "name": "%s"`, html.EscapeString(art.Publisher.Name)))
		if art.Publisher.URL != "" {
			sb.WriteString(fmt.Sprintf(`, "url": "%s"`, html.EscapeString(art.Publisher.URL)))
		}
		if art.Publisher.Logo != nil {
			sb.WriteString(`, "logo": {`)
			sb.WriteString(fmt.Sprintf(`"@type": "ImageObject", "url": "%s"`, html.EscapeString(art.Publisher.Logo.URL)))
			sb.WriteString("}")
		}
		sb.WriteString("},\n")
	}

	// Dates and Description
	if art.DatePublished != "" {
		sb.WriteString(fmt.Sprintf(`  "datePublished": "%s",`, html.EscapeString(art.DatePublished)))
		sb.WriteString("\n")
	}
	if art.DateModified != "" {
		sb.WriteString(fmt.Sprintf(`  "dateModified": "%s",`, html.EscapeString(art.DateModified)))
		sb.WriteString("\n")
	}
	if art.Description != "" {
		sb.WriteString(fmt.Sprintf(`  "description": "%s"`, html.EscapeString(art.Description)))
		sb.WriteString("\n")
	}

	sb.WriteString("}\n</script>")
	return sb.String()
}

func (art *Article) ensureDefaults() {
	if art.Context == "" {
		art.Context = "https://schema.org"
	}
	if art.Type == "" {
		art.Type = "Article"
	}

	if art.Author != nil {
		art.Author.ensureDefaults()
	}

	if art.Publisher != nil {
		art.Publisher.ensureDefaults()
	}
}
