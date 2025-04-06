package schemaorg

import (
	"fmt"
	"html/template"

	"github.com/a-h/templ"
	"github.com/indaco/teseo"
)

// Article represents a Schema.org Article object.
// For more details about the meaning of the properties see: https://schema.org/Article
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

// Validate checks if the Article has the recommended fields for SEO.
// It returns a slice of warning messages for missing recommended fields.
func (art *Article) Validate() []string {
	var warnings []string

	if art.Headline == "" {
		warnings = append(warnings, "missing recommended field: headline")
	}
	if len(art.Image) == 0 {
		warnings = append(warnings, "missing recommended field: image")
	}
	if art.DatePublished == "" {
		warnings = append(warnings, "missing recommended field: datePublished")
	}
	if art.Author == nil && art.Publisher == nil {
		warnings = append(warnings, "missing recommended field: author or publisher")
	}

	return warnings
}

// ToJsonLd converts the Article struct to a JSON-LD `templ.Component`.
func (art *Article) ToJsonLd() templ.Component {
	art.ensureDefaults()
	id := fmt.Sprintf("%s-%s", "article", teseo.GenerateUniqueKey())
	return templ.JSONScript(id, art).WithType("application/ld+json")
}

// ToGoHTMLJsonLd renders the Article struct as `template.HTML` value for Go's `html/template`.
func (art *Article) ToGoHTMLJsonLd() (template.HTML, error) {
	return teseo.RenderToHTML(art.ToJsonLd())
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
