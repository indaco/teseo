package teseo

import "github.com/a-h/templ"

// TemplRenderer is the interface for rendering content as a templ component.
type TemplRenderer interface {
	ToJsonLd() templ.Component
	ToMetaTags() templ.Component
}

// HtmlRenderer is the interface for rendering content as an HTML string for Go's `template/html`.
type HtmlRenderer interface {
	ToGoHTMLJsonLd() string
	ToGoHTMLMetaTags() string
}

// SitemapRenderer is the interface for rendering content to and from sitemap files.
type SitemapRenderer interface {
	ToSitemapFile() string             // Convert the content to a sitemap file format (e.g., XML).
	FromSitemapFile(data string) error // Load the content from a sitemap file.
}
