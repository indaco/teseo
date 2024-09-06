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
		return templ.JSONScript(teseo.GenerateUniqueKey(), p).WithType("application/ld+json").Render(ctx, w)
	})
}

// ToGoHTMLJsonLd renders the Product struct as a string for Go's `html/template`.
func (p *Product) ToGoHTMLJsonLd() string {
	p.ensureDefaults()

	var sb strings.Builder
	sb.WriteString(`<script type="application/ld+json">`)
	sb.WriteString("\n{\n")
	sb.WriteString(fmt.Sprintf(`  "@context": "%s",`, html.EscapeString(p.Context)))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf(`  "@type": "%s",`, html.EscapeString(p.Type)))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf(`  "name": "%s",`, html.EscapeString(p.Name)))
	sb.WriteString("\n")

	writeProductDescription(&sb, p)
	writeProductImages(&sb, p)
	writeProductSKU(&sb, p)
	writeProductBrand(&sb, p)
	writeProductOffers(&sb, p)
	writeProductCategory(&sb, p)
	writeProductAggregateRating(&sb, p)
	writeProductReviews(&sb, p)

	sb.WriteString("}\n</script>")
	return sb.String()
}

func writeProductDescription(sb *strings.Builder, p *Product) {
	if p.Description != "" {
		sb.WriteString(fmt.Sprintf(`  "description": "%s",`, html.EscapeString(p.Description)))
		sb.WriteString("\n")
	}
}

func writeProductImages(sb *strings.Builder, p *Product) {
	if len(p.Image) > 0 {
		sb.WriteString(`  "image": [`)
		for i, img := range p.Image {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf(`"%s"`, html.EscapeString(img)))
		}
		sb.WriteString("],\n")
	}
}

func writeProductSKU(sb *strings.Builder, p *Product) {
	if p.SKU != "" {
		sb.WriteString(fmt.Sprintf(`  "sku": "%s",`, html.EscapeString(p.SKU)))
		sb.WriteString("\n")
	}
}

func writeProductBrand(sb *strings.Builder, p *Product) {
	if p.Brand != nil {
		sb.WriteString(`  "brand": {`)
		sb.WriteString(fmt.Sprintf(`"@type": "Brand", "name": "%s"`, html.EscapeString(p.Brand.Name)))
		sb.WriteString("},\n")
	}
}

func writeProductOffers(sb *strings.Builder, p *Product) {
	if p.Offers != nil {
		sb.WriteString(`  "offers": {`)
		sb.WriteString(fmt.Sprintf(`"@type": "Offer", "price": "%s", "priceCurrency": "%s"`, html.EscapeString(p.Offers.Price), html.EscapeString(p.Offers.PriceCurrency)))
		if p.Offers.URL != "" {
			sb.WriteString(fmt.Sprintf(`, "url": "%s"`, html.EscapeString(p.Offers.URL)))
		}
		if p.Offers.Availability != "" {
			sb.WriteString(fmt.Sprintf(`, "availability": "%s"`, html.EscapeString(p.Offers.Availability)))
		}
		if p.Offers.ItemCondition != "" {
			sb.WriteString(fmt.Sprintf(`, "itemCondition": "%s"`, html.EscapeString(p.Offers.ItemCondition)))
		}
		sb.WriteString("},\n")
	}
}

func writeProductCategory(sb *strings.Builder, p *Product) {
	if p.Category != "" {
		sb.WriteString(fmt.Sprintf(`  "category": "%s",`, html.EscapeString(p.Category)))
		sb.WriteString("\n")
	}
}

func writeProductAggregateRating(sb *strings.Builder, p *Product) {
	if p.AggregateRating != nil {
		sb.WriteString(`  "aggregateRating": {`)
		sb.WriteString(fmt.Sprintf(`"@type": "AggregateRating", "ratingValue": %f, "reviewCount": %d`, p.AggregateRating.RatingValue, p.AggregateRating.ReviewCount))
		sb.WriteString("},\n")
	}
}

func writeProductReviews(sb *strings.Builder, p *Product) {
	if len(p.Review) > 0 {
		sb.WriteString(`  "review": [`)
		for i, review := range p.Review {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString("{\n")
			sb.WriteString(fmt.Sprintf(`    "@type": "Review", "reviewBody": "%s",`, html.EscapeString(review.ReviewBody)))
			sb.WriteString(fmt.Sprintf(`"datePublished": "%s",`, html.EscapeString(review.DatePublished)))
			if review.Author != nil {
				sb.WriteString(`    "author": {`)
				sb.WriteString(fmt.Sprintf(`"@type": "Person", "name": "%s"`, html.EscapeString(review.Author.Name)))
				sb.WriteString("},\n")
			}
			if review.ReviewRating != nil {
				sb.WriteString(`    "reviewRating": {`)
				sb.WriteString(fmt.Sprintf(`"@type": "Rating", "ratingValue": %f, "bestRating": %f`, review.ReviewRating.RatingValue, review.ReviewRating.BestRating))
				sb.WriteString("}\n")
			}
			sb.WriteString("}")
		}
		sb.WriteString("],\n")
	}
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
