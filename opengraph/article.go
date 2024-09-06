package opengraph

import (
	"context"
	"html/template"
	"io"
	"log"

	"github.com/a-h/templ"
	"github.com/indaco/teseo"
)

// Article represents the Open Graph article metadata.
//
// Example usage:
//
// Pure struct usage:
//
//	// Create an article using pure struct
//	article := &opengraph.Article{
//		OpenGraphObject: opengraph.OpenGraphObject{
//			Title:       "Example Article Title",
//			URL:         "https://www.example.com/articles/example-article",
//			Description: "This is an example article description.",
//			Image:       "https://www.example.com/images/article.jpg",
//		},
//		PublishedTime:  "2024-09-15T09:00:00Z",
//		ModifiedTime:   "2024-09-15T10:00:00Z",
//		ExpirationTime: "2024-12-31T23:59:59Z",
//		Author:         []string{"https://www.example.com/authors/jane-doe"},
//		Section:        "Technology",
//		Tag:            []string{"tech", "innovation", "example"},
//	}
//
// Factory method usage:
//
//	// Create an article
//	article := opengraph.NewArticle(
//		"Example Article Title",
//		"https://www.example.com/articles/example-article",
//		"This is an example article description.",
//		"https://www.example.com/images/article.jpg",
//		"2024-09-15T09:00:00Z",
//		"2024-09-15T10:00:00Z",
//		"2024-12-31T23:59:59Z",
//		[]string{"https://www.example.com/authors/jane-doe"},
//		"Technology",
//		[]string{"tech", "innovation", "example"},
//	)
//
// // Rendering the HTML meta tags using templ:
//
//	templ Page() {
//		@article.ToMetaTgs()
//	}
//
// // Rendering the HTML meta tags as `template.HTML` value:
//
//	metaTagsHtml := article.ToGoHTMLMetaTgs()
//
// Expected output:
//
//	<meta property="og:type" content="article"/>
//	<meta property="og:title" content="Example Article Title"/>
//	<meta property="og:url" content="https://www.example.com/articles/example-article"/>
//	<meta property="og:description" content="This is an example article description."/>
//	<meta property="og:image" content="https://www.example.com/images/article.jpg"/>
//	<meta property="article:published_time" content="2024-09-15T09:00:00Z"/>
//	<meta property="article:modified_time" content="2024-09-15T10:00:00Z"/>
//	<meta property="article:expiration_time" content="2024-12-31T23:59:59Z"/>
//	<meta property="article:section" content="Technology"/>
//	<meta property="article:author" content="https://www.example.com/authors/jane-doe"/>
//	<meta property="article:tag" content="tech"/>
//	<meta property="article:tag" content="innovation"/>
//	<meta property="article:tag" content="example"/>
type Article struct {
	OpenGraphObject
	PublishedTime  string   // article:published_time, the time the article was first published
	ModifiedTime   string   // article:modified_time, the time the article was last modified
	ExpirationTime string   // article:expiration_time, the time the article will expire
	Author         []string // article:author, URLs to the authors of the article
	Section        string   // article:section, a high-level section name
	Tag            []string // article:tag, tags of the article
}

// NewArticle initializes an Article with the default type "article".
func NewArticle(title, url, description, image, publishedTime, modifiedTime, expirationTime string, author []string, section string, tags []string) *Article {
	article := &Article{
		OpenGraphObject: OpenGraphObject{
			Title:       title,
			URL:         url,
			Description: description,
			Image:       image,
		},
		PublishedTime:  publishedTime,
		ModifiedTime:   modifiedTime,
		ExpirationTime: expirationTime,
		Author:         author,
		Section:        section,
		Tag:            tags,
	}
	article.ensureDefaults()
	return article
}

// ToMetaTags generates the HTML meta tags for the Open Graph Article using templ.Component.
func (art *Article) ToMetaTags() templ.Component {
	art.ensureDefaults()
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, tag := range art.metaTags() {
			if tag.content != "" {
				if err := teseo.WriteMetaTag(w, tag.property, tag.content); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// ToGoHTMLMetaTags generates the HTML meta tags for the Open Graph Audio as `template.HTML` value for Go's `html/template`.
func (art *Article) ToGoHTMLMetaTags() (template.HTML, error) {
	art.ensureDefaults()

	// Create the templ component.
	templComponent := art.ToMetaTags()

	// Render the templ component to a `template.HTML` value.
	html, err := templ.ToGoHTML(context.Background(), templComponent)
	if err != nil {
		log.Fatalf("failed to convert to html: %v", err)
	}

	return html, nil
}

// ensureDefaults sets default values for the Article object.
func (art *Article) ensureDefaults() {
	art.OpenGraphObject.ensureDefaults("article")
}

// metaTags returns all meta tags for the Article, including OpenGraphObject fields and article-specific ones.
func (art *Article) metaTags() []struct {
	property string
	content  string
} {
	tags := []struct {
		property string
		content  string
	}{
		{"og:type", "article"},
		{"og:title", art.Title},
		{"og:url", art.URL},
		{"og:description", art.Description},
		{"og:image", art.Image},
		{"article:published_time", art.PublishedTime},
		{"article:modified_time", art.ModifiedTime},
		{"article:expiration_time", art.ExpirationTime},
		{"article:section", art.Section},
	}

	// Add article:author tags
	for _, author := range art.Author {
		if author != "" {
			tags = append(tags, struct {
				property string
				content  string
			}{"article:author", author})
		}
	}

	// Add article:tag tags
	for _, tag := range art.Tag {
		if tag != "" {
			tags = append(tags, struct {
				property string
				content  string
			}{"article:tag", tag})
		}
	}

	return tags
}
