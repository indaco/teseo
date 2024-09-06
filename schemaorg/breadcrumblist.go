package schemaorg

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"
	"unicode"

	"github.com/a-h/templ"
	"github.com/indaco/teseo"
)

// BreadcrumbList represents a Schema.org BreadcrumbList object.
//
// Example usage:
//
// Pure struct usage:
//
//	breadcrumb := &schemaorg.BreadcrumbList{
//		ItemListElement: []schemaorg.ListItem{
//			{Name: "Home", Item: "https://www.example.com", Position: 1},
//			{Name: "About Us", Item: "https://www.example.com/about", Position: 2},
//		},
//	}
//
// Factory method usage:
//
//	breadcrumb := schemaorg.NewBreadcrumbList(
//		[]schemaorg.ListItem{
//			{Name: "Home", Item: "https://www.example.com", Position: 1},
//			{Name: "About Us", Item: "https://www.example.com/about", Position: 2},
//		},
//	)
//
// Example usage with `NewBreadcrumbListFromUrl`:
//
//	func HandleAbout(w http.ResponseWriter, r *http.Request) {
//		pageURL := teseo.GetFullURL(r)
//		breadcrumbList, err := schemaorg.NewBreadcrumbListFromUrl(pageURL)
//		if err != nil {
//			fmt.Println("Error generating breadcrumb list:", err)
//			return
//		}
//
//	// pass the component to your page and render with `@breadcrumbList.ToJsonLd()`
//		err = pages.AboutPage(breadcrumbList).Render(r.Context(), w)
//			if err != nil {
//			return
//		}
//	}
//
// // Rendering JSON-LD using templ:
//
//	templ Page() {
//		@breadcrumb.ToJsonLd()
//	}
//
// // Rendering JSON-LD as `template.HTML` value:
//
//	jsonLdHtml := breadcrumb.ToGoHTMLJsonLd()
//
// Expected output:
//
//	{
//		"@context": "https://schema.org",
//		"@type": "BreadcrumbList",
//		"itemListElement": [
//			{"@type": "ListItem", "position": 1, "name": "Home", "item": "https://www.example.com"},
//			{"@type": "ListItem", "position": 2, "name": "About Us", "item": "https://www.example.com/about"}
//		]
//	}
type BreadcrumbList struct {
	Context         string     `json:"@context"`
	Type            string     `json:"@type"`
	ItemListElement []ListItem `json:"itemListElement"`
}

// NewBreadcrumbList initializes an BreadcrumbList with default context and type.
func NewBreadcrumbList(listItem []ListItem) *BreadcrumbList {
	return &BreadcrumbList{
		Context:         "https://schema.org",
		Type:            "BreadcrumbList",
		ItemListElement: listItem,
	}
}

// NewBreadcrumbListFromUrl initializes an BreadcrumbList from the URL string.
func NewBreadcrumbListFromUrl(url string) (*BreadcrumbList, error) {
	bcl, err := createBreadcrumbListFromURL(url)
	if err != nil {
		return nil, fmt.Errorf("[NewBreadcrumbListFromUrl] invalid URL: %w", err)
	}
	return bcl, nil
}

// ToJsonLd converts the BreadcrumbList struct to a JSON-LD `templ.Component`.
func (bcl *BreadcrumbList) ToJsonLd() templ.Component {
	bcl.ensureDefaults()
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		return templ.JSONScript(teseo.GenerateUniqueKey(), bcl).WithType("application/ld+json").Render(ctx, w)
	})
}

// ToGoHTMLJsonLd renders the BreadcrumbList struct as a string for Go's `html/template`.
func (bcl *BreadcrumbList) ToGoHTMLJsonLd() (string, error) {
	bcl.ensureDefaults()
	// Create the templ component.
	templComponent := bcl.ToJsonLd()
	// Render the templ component to a `template.HTML` value.
	html, err := templ.ToGoHTML(context.Background(), templComponent)
	if err != nil {
		log.Fatalf("failed to convert to html: %v", err)
	}
	return string(html), nil
}

func (bcl *BreadcrumbList) ensureDefaults() {
	if bcl.Context == "" {
		bcl.Context = "https://schema.org"
	}
	if bcl.Type == "" {
		bcl.Type = "BreadcrumbList"
	}

	// Loop over each ListItem in ile and set its Type to "ListItem"
	for i := range bcl.ItemListElement {
		bcl.ItemListElement[i].Type = "ListItem"
	}
}

// createBreadcrumbListFromURL generates a BreadcrumbList JSON-LD object from a URL string.
func createBreadcrumbListFromURL(rawURL string) (*BreadcrumbList, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("[createBreadcrumbListFromURL] invalid URL: %w", err)
	}

	// Extract segments from the URL path.
	segments := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")

	// Initialize the base URL correctly.
	baseURL := parsedURL.Scheme + "://" + parsedURL.Host

	var listItems []ListItem

	// Always include the base URL as the first breadcrumb item.
	listItems = append(listItems, ListItem{
		Type:     "ListItem",
		Position: 1,
		Name:     "Home",
		Item:     baseURL,
	})

	// Check if there are additional segments beyond the base URL.
	if len(segments) > 0 && segments[0] != "" {
		// Build the ListItem slice for JSON-LD
		for i, segment := range segments {
			// Correctly concatenate the base URL with the segments.
			href := baseURL + "/" + strings.Join(segments[:i+1], "/")
			listItems = append(listItems, ListItem{
				Type:     "ListItem",
				Position: i + 2, // Start from 2 because the base URL is already position 1
				Name:     toTitle(segment),
				Item:     href,
			})
		}
	}

	// Create and return the BreadcrumbList object
	breadcrumbList := &BreadcrumbList{
		Context:         "https://schema.org",
		Type:            "BreadcrumbList",
		ItemListElement: listItems,
	}

	return breadcrumbList, nil
}

// ToTitle converts the first letter of a string to its title case equivalent.
// Useful for handling languages or characters where the title case differs from the uppercase.
// Example: in German, 'ß' will be converted to 'ẞ'.
func toTitle(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToTitle(runes[0])
	return string(runes)
}
