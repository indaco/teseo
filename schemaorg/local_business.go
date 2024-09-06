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

// LocalBusiness represents a Schema.org LocalBusiness object.
//
// Example usage:
//
// Pure struct usage:
//
//	localBusiness := &schemaorg.LocalBusiness{
//		Name:        "Example Business",
//		Address:     &schemaorg.PostalAddress{StreetAddress: "123 Main St", AddressLocality: "Anytown", AddressRegion: "CA", PostalCode: "12345"},
//		Telephone:   "+1-800-555-1234",
//		Description: "This is an example local business.",
//	}
//
// Factory method usage:
//
//	localBusiness := schemaorg.NewLocalBusiness(
//		"Example Business",
//		&schemaorg.PostalAddress{StreetAddress: "123 Main St", AddressLocality: "Anytown", AddressRegion: "CA", PostalCode: "12345"},
//		"+1-800-555-1234",
//		"This is an example local business",
//	)
//
// // Rendering JSON-LD using templ:
//
//	templ Page() {
//		@localBusiness.ToJsonLd()
//	}
//
// // Rendering JSON-LD as `template.HTML` value:
//
//	jsonLdHtml := localBusiness.ToGoHTMLJsonLd()
//
// Expected output:
//
//	{
//		"@context": "https://schema.org",
//		"@type": "LocalBusiness",
//		"name": "Example Business",
//		"address": {
//			"@type": "PostalAddress",
//			"streetAddress": "123 Main St",
//			"addressLocality": "Anytown",
//			"addressRegion": "CA",
//			"postalCode": "12345"
//		},
//		"telephone": "+1-800-555-1234",
//		"description": "This is an example local business"
//	}
type LocalBusiness struct {
	Context         string           `json:"@context"`
	Type            string           `json:"@type"`
	Name            string           `json:"name,omitempty"`
	Description     string           `json:"description,omitempty"`
	URL             string           `json:"url,omitempty"`
	Logo            *ImageObject     `json:"logo,omitempty"`
	Telephone       string           `json:"telephone,omitempty"`
	Address         *PostalAddress   `json:"address,omitempty"`
	OpeningHours    []string         `json:"openingHours,omitempty"`
	Geo             *GeoCoordinates  `json:"geo,omitempty"`
	AggregateRating *AggregateRating `json:"aggregateRating,omitempty"`
	Review          []*Review        `json:"review,omitempty"`
}

// GeoCoordinates represents a Schema.org GeoCoordinates object
type GeoCoordinates struct {
	Type      string  `json:"@type"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

// NewLocalBusiness initializes a LocalBusiness with default context and type.
func NewLocalBusiness(name string, description string, url string, telephone string, logo *ImageObject, address *PostalAddress, openingHours []string, geo *GeoCoordinates, aggregateRating *AggregateRating, reviews []*Review) *LocalBusiness {
	localBusiness := &LocalBusiness{
		Name:            name,
		Description:     description,
		URL:             url,
		Logo:            logo,
		Telephone:       telephone,
		Address:         address,
		OpeningHours:    openingHours,
		Geo:             geo,
		AggregateRating: aggregateRating,
		Review:          reviews,
	}
	localBusiness.ensureDefaults()
	return localBusiness
}

// ToJsonLd converts the LocalBusiness struct to a JSON-LD `templ.Component`.
func (lb *LocalBusiness) ToJsonLd() templ.Component {
	lb.ensureDefaults()
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		return templ.JSONScript(teseo.GenerateUniqueKey(), lb).WithType("application/ld+json").Render(ctx, w)
	})
}

// ToGoHTMLJsonLd renders the LocalBusiness struct as a string for Go's `html/template`.
func (lb *LocalBusiness) ToGoHTMLJsonLd() string {
	lb.ensureDefaults()

	var sb strings.Builder
	sb.WriteString(`<script type="application/ld+json">`)
	sb.WriteString("\n{\n")
	sb.WriteString(fmt.Sprintf(`  "@context": "%s",`, html.EscapeString(lb.Context)))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf(`  "@type": "%s",`, html.EscapeString(lb.Type)))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf(`  "name": "%s",`, html.EscapeString(lb.Name)))
	sb.WriteString("\n")

	writeOptionalFields(&sb, lb)
	writeLogo(&sb, lb)
	writeAddress(&sb, lb)
	writeOpeningHours(&sb, lb)
	writeGeo(&sb, lb)
	writeAggregateRating(&sb, lb)
	writeReviews(&sb, lb)

	sb.WriteString("}\n</script>")
	return sb.String()
}

func writeOptionalFields(sb *strings.Builder, lb *LocalBusiness) {
	if lb.Description != "" {
		sb.WriteString(fmt.Sprintf(`  "description": "%s",`, html.EscapeString(lb.Description)))
		sb.WriteString("\n")
	}
	if lb.URL != "" {
		sb.WriteString(fmt.Sprintf(`  "url": "%s",`, html.EscapeString(lb.URL)))
		sb.WriteString("\n")
	}
	if lb.Telephone != "" {
		sb.WriteString(fmt.Sprintf(`  "telephone": "%s",`, html.EscapeString(lb.Telephone)))
		sb.WriteString("\n")
	}
}

func writeLogo(sb *strings.Builder, lb *LocalBusiness) {
	if lb.Logo != nil {
		sb.WriteString(`  "logo": {`)
		sb.WriteString(fmt.Sprintf(`"@type": "ImageObject", "url": "%s"`, html.EscapeString(lb.Logo.URL)))
		sb.WriteString("},\n")
	}
}

func writeAddress(sb *strings.Builder, lb *LocalBusiness) {
	if lb.Address != nil {
		sb.WriteString(`  "address": {`)
		sb.WriteString(fmt.Sprintf(`"@type": "PostalAddress", "addressLocality": "%s", "addressCountry": "%s"`, html.EscapeString(lb.Address.AddressLocality), html.EscapeString(lb.Address.AddressCountry)))
		sb.WriteString("},\n")
	}
}

func writeOpeningHours(sb *strings.Builder, lb *LocalBusiness) {
	if len(lb.OpeningHours) > 0 {
		sb.WriteString(`  "openingHours": [`)
		for i, hours := range lb.OpeningHours {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf(`"%s"`, html.EscapeString(hours)))
		}
		sb.WriteString("],\n")
	}
}

func writeGeo(sb *strings.Builder, lb *LocalBusiness) {
	if lb.Geo != nil {
		sb.WriteString(`  "geo": {`)
		sb.WriteString(fmt.Sprintf(`"@type": "GeoCoordinates", "latitude": %f, "longitude": %f`, lb.Geo.Latitude, lb.Geo.Longitude))
		sb.WriteString("},\n")
	}
}

func writeAggregateRating(sb *strings.Builder, lb *LocalBusiness) {
	if lb.AggregateRating != nil {
		sb.WriteString(`  "aggregateRating": {`)
		sb.WriteString(fmt.Sprintf(`"@type": "AggregateRating", "ratingValue": %f, "reviewCount": %d`, lb.AggregateRating.RatingValue, lb.AggregateRating.ReviewCount))
		sb.WriteString("},\n")
	}
}

func writeReviews(sb *strings.Builder, lb *LocalBusiness) {
	if len(lb.Review) > 0 {
		sb.WriteString(`  "review": [`)
		for i, review := range lb.Review {
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

// ensureDefaults sets default values for LocalBusiness and its nested objects if they are not already set.
func (lb *LocalBusiness) ensureDefaults() {
	if lb.Context == "" {
		lb.Context = "https://schema.org"
	}
	if lb.Type == "" {
		lb.Type = "LocalBusiness"
	}

	if lb.Logo != nil {
		lb.Logo.ensureDefaults()
	}

	if lb.Address != nil {
		lb.Address.ensureDefaults()
	}

	if lb.Geo != nil {
		lb.Geo.ensureDefaults()
	}

	if lb.AggregateRating != nil {
		lb.AggregateRating.ensureDefaults()
	}

	for _, review := range lb.Review {
		review.ensureDefaults()
	}
}

// ensureDefaults sets default values for GeoCoordinates if they are not already set.
func (geo *GeoCoordinates) ensureDefaults() {
	if geo.Type == "" {
		geo.Type = "GeoCoordinates"
	}
}
