package opengraph

import (
	"context"
	"html/template"
	"io"
	"log"

	"github.com/a-h/templ"
	"github.com/indaco/teseo"
)

// ProductGroup represents the Open Graph product group metadata.
//
// Example usage:
//
// Pure struct usage:
//
//	// Create a product group using pure struct
//	productGroup := &opengraph.ProductGroup{
//		OpenGraphObject: opengraph.OpenGraphObject{
//			Title:       "Example Product Group",
//			URL:         "https://www.example.com/product-group/example-product-group",
//			Description: "This is an example product group description.",
//			Image:       "https://www.example.com/images/product-group.jpg",
//		},
//		Products: []string{
//			"https://www.example.com/product/product-1",
//			"https://www.example.com/product/product-2",
//		},
//	}
//
// Factory method usage:
//
//	// Create a product group using the factory method
//	productGroup := opengraph.NewProductGroup(
//		"Example Product Group",
//		"https://www.example.com/product-group/example-product-group",
//		"This is an example product group description.",
//		"https://www.example.com/images/product-group.jpg",
//		[]string{
//			"https://www.example.com/product/product-1",
//			"https://www.example.com/product/product-2",
//		},
//	)
//
// // Rendering the HTML meta tags using templ:
//
//	templ Page() {
//		@productGroup.ToMetaTgs()
//	}
//
// // Rendering the HTML meta tags as `template.HTML` value:
//
//	metaTagsHtml := productGroup.ToGoHTMLMetaTgs()
//
// Expected output:
//
//	<meta property="og:type" content="product.group"/>
//	<meta property="og:title" content="Example Product Group"/>
//	<meta property="og:url" content="https://www.example.com/product-group/example-product-group"/>
//	<meta property="og:description" content="This is an example product group description."/>
//	<meta property="og:image" content="https://www.example.com/images/product-group.jpg"/>
//	<meta property="product:group_item" content="https://www.example.com/product/product-1"/>
//	<meta property="product:group_item" content="https://www.example.com/product/product-2"/>
type ProductGroup struct {
	OpenGraphObject
	Products []string // product:group_item, URLs to individual products in the group
}

// NewProductGroup initializes a ProductGroup with the default type "product.group".
func NewProductGroup(title, url, description, image string, products []string) *ProductGroup {
	productGroup := &ProductGroup{
		OpenGraphObject: OpenGraphObject{
			Title:       title,
			URL:         url,
			Description: description,
			Image:       image,
		},
		Products: products,
	}
	productGroup.ensureDefaults()
	return productGroup
}

// ToMetaTags generates the HTML meta tags for the Open Graph Product Group as templ.Component.
func (pg *ProductGroup) ToMetaTags() templ.Component {
	pg.ensureDefaults()
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, tag := range pg.metaTags() {
			if tag.content != "" {
				if err := teseo.WriteMetaTag(w, tag.property, tag.content); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// ToGoHTMLMetaTags generates the HTML meta tags for the Open Graph Product Group as `template.HTML` value for Go's `html/template`.
func (pg *ProductGroup) ToGoHTMLMetaTags() (template.HTML, error) {
	// Create the templ component.
	templComponent := pg.ToMetaTags()

	// Render the templ component to a `template.HTML` value.
	html, err := templ.ToGoHTML(context.Background(), templComponent)
	if err != nil {
		log.Fatalf("failed to convert to html: %v", err)
	}

	return html, nil
}

// ensureDefaults sets default values for ProductGroup.
func (pg *ProductGroup) ensureDefaults() {
	pg.OpenGraphObject.ensureDefaults("product.group")
}

// metaTags returns all meta tags for the ProductGroup object, including OpenGraphObject fields and product-specific ones.
func (pg *ProductGroup) metaTags() []struct {
	property string
	content  string
} {
	tags := []struct {
		property string
		content  string
	}{
		{"og:type", "product.group"},
		{"og:title", pg.Title},
		{"og:url", pg.URL},
		{"og:description", pg.Description},
		{"og:image", pg.Image},
	}

	// Add product:group_item tags for each product in the group
	for _, product := range pg.Products {
		if product != "" {
			tags = append(tags, struct {
				property string
				content  string
			}{"product:group_item", product})
		}
	}

	return tags
}
