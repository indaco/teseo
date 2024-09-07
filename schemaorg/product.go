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

// Product represents a Schema.org Product object.
//
// Example usage:
//
// Pure struct usage:
//
//	product := &schemaorg.Product{
//		Name:        "Example Product",
//		Description: "This is an example product description.",
//		SKU:         "12345",
//		Brand:       &schemaorg.Brand{Name: "Example Brand"},
//		Offers:      &schemaorg.Offer{Price: "29.99", PriceCurrency: "USD"},
//	}
//
// Factory method usage:
//
//	product := schemaorg.NewProduct(
//		"Example Product",
//		"This is an example product description.",
//		"12345",
//		&schemaorg.Brand{Name: "Example Brand"},
//		&schemaorg.Offer{Price: "29.99", PriceCurrency: "USD"},
//	)
//
// // Rendering JSON-LD using templ:
//
//	templ Page() {
//		@product.ToJsonLd()
//	}
//
// // Rendering JSON-LD as `template.HTML` value:
//
//	jsonLdHtml := product.ToGoHTMLJsonLd()
//
// Expected output:
//
//	{
//		"@context": "https://schema.org",
//		"@type": "Product",
//		"name": "Example Product",
//		"description": "This is an example product description",
//		"sku": "12345",
//		"brand": {"@type": "Brand", "name": "Example Brand"},
//		"offers": {
//			"@type": "Offer",
//			"price": "29.99",
//			"priceCurrency": "USD"
//		}
//	}
type Product struct {
	Context         string           `json:"@context"`
	Type            string           `json:"@type"`
	Name            string           `json:"name,omitempty"`
	Description     string           `json:"description,omitempty"`
	Image           []string         `json:"image,omitempty"`
	SKU             string           `json:"sku,omitempty"`
	Brand           *Brand           `json:"brand,omitempty"`
	Offers          *Offer           `json:"offers,omitempty"`
	Category        string           `json:"category,omitempty"`
	AggregateRating *AggregateRating `json:"aggregateRating,omitempty"`
	Review          []*Review        `json:"review,omitempty"`
}

// Brand represents a Schema.org Brand object
type Brand struct {
	Type string `json:"@type"`
	Name string `json:"name,omitempty"`
}

// Offer represents a Schema.org Offer object
type Offer struct {
	Type          string `json:"@type"`
	URL           string `json:"url,omitempty"`
	PriceCurrency string `json:"priceCurrency,omitempty"`
	Price         string `json:"price,omitempty"`
	Availability  string `json:"availability,omitempty"`
	ItemCondition string `json:"itemCondition,omitempty"`
}

// AggregateRating represents a Schema.org AggregateRating object
type AggregateRating struct {
	Type        string  `json:"@type"`
	RatingValue float64 `json:"ratingValue,omitempty"`
	ReviewCount int     `json:"reviewCount,omitempty"`
}

// Review represents a Schema.org Review object
type Review struct {
	Type          string  `json:"@type"`
	Author        *Person `json:"author,omitempty"`
	DatePublished string  `json:"datePublished,omitempty"`
	ReviewBody    string  `json:"reviewBody,omitempty"`
	ReviewRating  *Rating `json:"reviewRating,omitempty"`
}

// Rating represents a Schema.org Rating object
type Rating struct {
	Type        string  `json:"@type"`
	RatingValue float64 `json:"ratingValue,omitempty"`
	BestRating  float64 `json:"bestRating,omitempty"`
}

// NewProduct initializes a Product with default context and type.
func NewProduct(name, description string, image []string, sku string, brand *Brand, offers *Offer, category string, aggregateRating *AggregateRating, reviews []*Review) *Product {
	product := &Product{
		Name:            name,
		Description:     description,
		Image:           image,
		SKU:             sku,
		Brand:           brand,
		Offers:          offers,
		Category:        category,
		AggregateRating: aggregateRating,
		Review:          reviews,
	}
	product.ensureDefaults()
	return product
}

// ToJsonLd converts the Product struct to a JSON-LD `templ.Component`.
func (p *Product) ToJsonLd() templ.Component {
	p.ensureDefaults()
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		id := fmt.Sprintf("%s-%s", "product", teseo.GenerateUniqueKey())
		return templ.JSONScript(id, p).WithType("application/ld+json").Render(ctx, w)
	})
}

// ToGoHTMLJsonLd renders the Product struct as `template.HTML` value for Go's `html/template`.
func (p *Product) ToGoHTMLJsonLd() (template.HTML, error) {
	// Create the templ component.
	templComponent := p.ToJsonLd()

	// Render the templ component to a `template.HTML` value.
	html, err := templ.ToGoHTML(context.Background(), templComponent)
	if err != nil {
		log.Fatalf("failed to convert to html: %v", err)
	}

	return html, nil
}

// ensureDefaults sets default values for Product and its nested objects if they are not already set.
func (p *Product) ensureDefaults() {
	if p.Context == "" {
		p.Context = "https://schema.org"
	}

	if p.Type == "" {
		p.Type = "Product"
	}

	if p.Brand != nil {
		p.Brand.ensureDefaults()
	}

	if p.Offers != nil {
		p.Offers.ensureDefaults()
	}

	if p.AggregateRating != nil {
		p.AggregateRating.ensureDefaults()
	}

	for _, review := range p.Review {
		review.ensureDefaults()

	}
}

// ensureDefaults sets default values for Brand if they are not already set.
func (b *Brand) ensureDefaults() {
	if b.Type == "" {
		b.Type = "Brand"
	}
}

// ensureDefaults sets default values for Offer if they are not already set.
func (o *Offer) ensureDefaults() {
	if o.Type == "" {
		o.Type = "Offer"
	}
}

// ensureDefaults sets default values for AggregateRating if they are not already set.
func (ar *AggregateRating) ensureDefaults() {
	if ar.Type == "" {
		ar.Type = "AggregateRating"
	}
}

// ensureDefaults sets default values for Review if they are not already set.
func (r *Review) ensureDefaults() {
	if r.Type == "" {
		r.Type = "Review"
	}

	if r.Author != nil {
		r.Author.ensureDefaults()
	}

	if r.ReviewRating != nil {
		r.ReviewRating.ensureDefaults()
	}
}

// ensureDefaults sets default values for Rating if they are not already set.
func (ra *Rating) ensureDefaults() {
	if ra.Type == "" {
		ra.Type = "Rating"
	}
}
