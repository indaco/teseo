package opengraph

import (
	"context"
	"html/template"
	"io"
	"log"

	"github.com/a-h/templ"
	"github.com/indaco/teseo"
)

// Book represents the Open Graph book metadata.
//
// Example usage:
//
// Pure struct usage:
//
//	// Create a book using pure struct
//	book := &opengraph.Book{
//		OpenGraphObject: opengraph.OpenGraphObject{
//			Title:       "Example Book Title",
//			URL:         "https://www.example.com/books/example-book",
//			Description: "This is an example book description.",
//			Image:       "https://www.example.com/images/book.jpg",
//		},
//		ISBN:        "978-3-16-148410-0",
//		ReleaseDate: "2024-09-15",
//		Author:      []string{"https://www.example.com/authors/jane-doe"},
//		Tag:         []string{"fiction", "bestseller", "example"},
//	}
//
// Factory method usage:
//
//	// Create a book
//	book := opengraph.NewBook(
//		"Example Book Title",
//		"https://www.example.com/books/example-book",
//		"This is an example book description.",
//		"https://www.example.com/images/book.jpg",
//		"978-3-16-148410-0",
//		"2024-09-15",
//		[]string{"https://www.example.com/authors/jane-doe"},
//		[]string{"fiction", "bestseller", "example"},
//	)
//
// // Rendering the HTML meta tags using templ:
//
//	templ Page() {
//		@book.ToMetaTgs()
//	}
//
// // Rendering the HTML meta tags as `template.HTML` value:
//
//	metaTagsHtml := book.ToGoHTMLMetaTgs()
//
// Expected output:
//
//	<meta property="og:type" content="book"/>
//	<meta property="og:title" content="Example Book Title"/>
//	<meta property="og:url" content="https://www.example.com/books/example-book"/>
//	<meta property="og:description" content="This is an example book description."/>
//	<meta property="og:image" content="https://www.example.com/images/book.jpg"/>
//	<meta property="book:isbn" content="978-3-16-148410-0"/>
//	<meta property="book:release_date" content="2024-09-15"/>
//	<meta property="book:author" content="https://www.example.com/authors/jane-doe"/>
//	<meta property="book:tag" content="fiction"/>
//	<meta property="book:tag" content="bestseller"/>
//	<meta property="book:tag" content="example"/>
type Book struct {
	OpenGraphObject
	Author      []string // book:author, URLs to the authors of the book
	ISBN        string   // book:isbn, ISBN number of the book
	ReleaseDate string   // book:release_date, the release date of the book
	Tag         []string // book:tag, tags for the book
}

// NewBook initializes a Book with the default type "book".
func NewBook(title, url, description, image, isbn, releaseDate string, author, tags []string) *Book {
	book := &Book{
		OpenGraphObject: OpenGraphObject{
			Title:       title,
			URL:         url,
			Description: description,
			Image:       image,
		},
		Author:      author,
		ISBN:        isbn,
		ReleaseDate: releaseDate,
		Tag:         tags,
	}
	book.ensureDefaults()
	return book
}

// ToMetaTags generates the HTML meta tags for the Open Graph Book as templ.Component.
func (book *Book) ToMetaTags() templ.Component {
	book.ensureDefaults()
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, tag := range book.metaTags() {
			if tag.content != "" {
				if err := teseo.WriteMetaTag(w, tag.property, tag.content); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// ToGoHTMLMetaTags generates the HTML meta tags for the Open Graph Book as `template.HTML` value for Go's `html/template`.
func (book *Book) ToGoHTMLMetaTags() (template.HTML, error) {
	// Create the templ component.
	templComponent := book.ToMetaTags()

	// Render the templ component to a `template.HTML` value.
	html, err := templ.ToGoHTML(context.Background(), templComponent)
	if err != nil {
		log.Fatalf("failed to convert to html: %v", err)
	}

	return html, nil
}

// ensureDefaults sets default values for Book.
func (book *Book) ensureDefaults() {
	book.OpenGraphObject.ensureDefaults("book")
}

// metaTags returns all meta tags for the Book object, including OpenGraphObject fields and book-specific ones.
func (book *Book) metaTags() []struct {
	property string
	content  string
} {
	tags := []struct {
		property string
		content  string
	}{
		{"og:type", "book"},
		{"og:title", book.Title},
		{"og:url", book.URL},
		{"og:description", book.Description},
		{"og:image", book.Image},
		{"book:isbn", book.ISBN},
		{"book:release_date", book.ReleaseDate},
	}

	// Add book:author tags
	for _, author := range book.Author {
		if author != "" {
			tags = append(tags, struct {
				property string
				content  string
			}{"book:author", author})
		}
	}

	// Add book:tag tags
	for _, tag := range book.Tag {
		if tag != "" {
			tags = append(tags, struct {
				property string
				content  string
			}{"book:tag", tag})
		}
	}

	return tags
}
